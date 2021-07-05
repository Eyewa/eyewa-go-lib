package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eyewa/eyewa-go-lib/base"
	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	amqptracing "github.com/eyewa/eyewa-go-lib/tracing/amqp"
	"github.com/ory/viper"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	config        Config
	exchangeTypes = map[string]string{
		amqp.ExchangeDirect:  amqp.ExchangeDirect,
		amqp.ExchangeFanout:  amqp.ExchangeFanout,
		amqp.ExchangeHeaders: amqp.ExchangeHeaders,
		amqp.ExchangeTopic:   amqp.ExchangeTopic,
	}
	defaultPrefetchCount = 5
)

func initConfig() (Config, string, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envVars := []string{
		"SERVICE_NAME",
		"HOSTNAME",
		"RABBITMQ_SERVER",
		"RABBITMQ_AMQP_PORT",
		"RABBITMQ_USERNAME",
		"RABBITMQ_PASSWORD",
		"PUBLISHER_QUEUE_NAME",
		"CONSUMER_QUEUE_NAME",
		"QUEUE_PREFETCH_COUNT",
		"RABBITMQ_PUBLISHER_EXCHANGE_TYPE",
		"RABBITMQ_CONSUMER_EXCHANGE_TYPE",
	}

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return config, "", err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password,
		config.Server, config.AmqpPort), nil
}

// NewRMQClient new rmq client
func NewRMQClient() *RMQClient {
	return &RMQClient{
		mutex:      new(sync.RWMutex),
		connection: nil,
		channels:   make(map[string]*amqp.Channel),
	}
}

// Connect establishes connnection to the message broker of choice
func (rmq *RMQClient) Connect() error {
	// if a connection already exists, back off.
	if rmq.connection != nil {
		return nil
	}

	// init configs
	_, connStr, err := initConfig()
	if err != nil {
		return err
	}

	// if no queues are specified, back off.
	if config.ConsumerQueueName == "" && config.PublisherQueueName == "" {
		return libErrs.ErrorNoQueuesSpecified
	}

	// establish connection
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}

	rmq.connection = conn
	rmq.mutex = new(sync.RWMutex)
	rmq.channels = make(map[string]*amqp.Channel)

	// create channel for consuming (if any)
	if config.ConsumerQueueName != "" {
		if err := rmq.createConsumerChannel(); err != nil {
			return err
		}
	}

	// create channel for publish (if any)
	if config.PublisherQueueName != "" {
		if err := rmq.createPublisherChannel(); err != nil {
			return err
		}
	}

	// connection listener
	rmq.ConnectionListener()

	return nil
}

// Consume consumes messages from a queue
func (rmq *RMQClient) Consume(queue string, callback base.MessageBrokerCallbackFunc) {
	var (
		ctx     = context.Background()
		endSpan func()
	)

	defer func() {
		// reaching here means the connection meant to be long lived has died.
		_ = callback(ctx, nil, libErrs.ErrorLostConnectionToMessageBroker)
		defer endSpan()
	}()

	rmq.mutex.RLock()
	channel, exists := rmq.channels[queue]
	rmq.mutex.RUnlock()

	// check if channel exists for queue
	// if not create it
	if !exists && channel == nil {
		channel, err := rmq.CreateNewChannel(config.ConsumerQueueName)
		if err != nil {
			_ = callback(ctx, nil, err)
			return
		}

		errQ := rmq.declareQueue(channel, config.ConsumerQueueName, config.ConsumerExchangeType)
		if errQ != nil {
			_ = callback(ctx, nil, errQ)
			return
		}
	}

	rmq.mutex.RLock()
	channel, exists = rmq.channels[queue]
	rmq.mutex.RUnlock()

	if exists && channel != nil {
		log.Info(fmt.Sprintf("Listening to %s for new messages...", queue))

		// attempt to consume events from broker
		msgs, err := channel.Consume(queue, getNameForChannel(queue), false, false, false, false, nil)
		if err != nil {
			_ = callback(ctx, nil, fmt.Errorf(libErrs.ErrorConsumeFailure.Error(), queue, err))
			return
		}

		var event *base.EyewaEvent
		for msg := range msgs {
			// start tracing
			ctx, endSpan = amqptracing.StartDeliverySpan(ctx, msg)

			// attempt to unmarshal event
			err := json.Unmarshal(msg.Body, &event)
			if err != nil {
				errMsg := fmt.Errorf(libErrs.ErrorEventUnmarshalFailure.Error(), queue, err)
				_ = callback(ctx, nil, errMsg)

				// nack message and remove from queue
				err = msg.Nack(false, false)
				if err != nil {
					_ = callback(ctx, nil, err)
				}

				// publish message to DL
				err := rmq.sendToDeadletterQueue(msg, errMsg)
				if err != nil {
					callback(ctx, nil, err)
				}
			} else {
				// ONLY ack message once there is no error from callback.
				if err := callback(ctx, event, nil); err == nil {
					if err = msg.Ack(false); err != nil {
						// nack message and push back to queue
						err = msg.Nack(false, true)
						if err != nil {
							log.Error(libErrs.ErrorNackFailure.Error(),
								zap.String("queue", queue),
								zap.String("event", string(msg.Body)),
								zap.Error(err))
						}

						log.Error(libErrs.ErrorAckFailure.Error(),
							zap.String("queue", queue),
							zap.String("event", string(msg.Body)),
							zap.Error(err))
					}
				}
			}
		}
	}
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string, event *base.EyewaEvent, callback base.MessageBrokerCallbackFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	rmq.mutex.RLock()
	channel, exists := rmq.channels[queue]
	rmq.mutex.RUnlock()

	ctx := context.Background()

	// determine if channel exists for queue
	if !exists && channel == nil {
		channel, err := rmq.CreateNewChannel(config.PublisherQueueName)
		if err != nil {
			_ = callback(ctx, event, err)
			return
		}

		errQ := rmq.declareQueue(channel, config.PublisherQueueName, config.PublisherExchangeType)
		if errQ != nil {
			_ = callback(ctx, event, errQ)
			return
		}
	}

	rmq.mutex.RLock()
	channel, exists = rmq.channels[queue]
	rmq.mutex.RUnlock()

	if exists && channel != nil {
		// attempt to marshal event for publishing
		eventJSON, err := json.Marshal(&event)
		if err != nil {
			_ = callback(ctx, event, err)
			return
		}

		msg := amqp.Publishing{
			ContentType:  "application/json",
			Body:         eventJSON,
			DeliveryMode: amqp.Persistent,
		}

		// start tracing the publishing span
		ctx, endSpan := amqptracing.StartPublishingSpan(ctx, msg)
		defer endSpan()

		// attempt to publish event
		err = channel.Publish("", config.PublisherQueueName, false, false, msg)
		if err != nil {
			_ = callback(ctx, event, libErrs.ErrorFailedToPublishEvent)
		}

		_ = callback(ctx, event, nil)
		return
	}
}

// CloseConnection closes a connection as well as any
// underlying channels associated to it.
func (rmq *RMQClient) CloseConnection() error {
	if rmq.connection != nil {
		return rmq.connection.Close()
	}

	rmq.mutex.Lock()
	defer rmq.mutex.Unlock()
	rmq.channels = make(map[string]*amqp.Channel)

	return nil
}

func (rmq *RMQClient) declareQueue(channel *amqp.Channel, queue, exchangeType string) error {
	var err error

	if channel == nil {
		if channel, err = rmq.CreateNewChannel(queue); err != nil {
			return err
		}
	}

	exchType := exchangeTypes[exchangeType]
	exchName := ""
	if exchType != "" {
		exchName = fmt.Sprintf("%s.%s", queue, exchType)
	}

	// declare queue
	q, err := channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf(libErrs.ErrorQueueDeclareFailure.Error(), q.Name, err)
	}

	// declare exchange
	err = channel.ExchangeDeclare(exchName, exchType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf(libErrs.ErrorExchangeDeclareFailure.Error(), q.Name, err)
	}

	// bind them together
	err = channel.QueueBind(queue, queue, exchName, false, nil)
	if err != nil {
		return fmt.Errorf(libErrs.ErrorExchangeBindFailure.Error(), q.Name, err)
	}

	return nil
}

func (rmq *RMQClient) createConsumerChannel() error {
	if config.ConsumerQueueName == "" {
		return libErrs.ErrorNoConsumerQueueSpecified
	}

	rmq.mutex.Lock()
	defer rmq.mutex.Unlock()

	if _, exists := rmq.channels[config.ConsumerQueueName]; !exists {
		conCh, err := rmq.connection.Channel()
		if err != nil {
			return err
		}

		prefetchCount, _ := strconv.Atoi(config.QueuePrefetchCount)
		if prefetchCount == 0 {
			prefetchCount = defaultPrefetchCount
		}

		if err := conCh.Qos(prefetchCount, 0, true); err != nil {
			return err
		}

		if err := rmq.declareQueue(conCh, config.ConsumerQueueName, config.ConsumerExchangeType); err != nil {
			return err
		}

		rmq.channels[config.ConsumerQueueName] = conCh
	}

	return nil
}

func (rmq *RMQClient) createPublisherChannel() error {
	if config.PublisherQueueName == "" {
		return libErrs.ErrorNoPublisherQueueSpecified
	}

	rmq.mutex.Lock()
	defer rmq.mutex.Unlock()

	if _, exists := rmq.channels[config.PublisherQueueName]; !exists {
		pubCh, err := rmq.connection.Channel()
		if err != nil {
			return err
		}

		prefetchCount, _ := strconv.Atoi(config.QueuePrefetchCount)
		if prefetchCount == 0 {
			prefetchCount = defaultPrefetchCount
		}

		if err := pubCh.Qos(prefetchCount, 0, true); err != nil {
			return err
		}

		if err := rmq.declareQueue(pubCh, config.PublisherQueueName, config.PublisherExchangeType); err != nil {
			return err
		}

		rmq.channels[config.PublisherQueueName] = pubCh
	}

	return nil
}

// CreateNewChannel creates a new channel for specified queue
func (rmq *RMQClient) CreateNewChannel(queue string) (*amqp.Channel, error) {
	if rmq.connection != nil {
		channel, err := rmq.connection.Channel()
		if err != nil {
			return nil, fmt.Errorf(libErrs.ErrorChannelCreateFailure.Error(), queue, err)
		}

		rmq.mutex.Lock()
		rmq.channels[queue] = channel
		rmq.mutex.Unlock()

		return rmq.channels[queue], nil
	}

	return nil, libErrs.ErrorNoRMQConnection
}

// QueueInspect inspects a queue and returns no. of consumers + messages
// currently unacked.
func (rmq *RMQClient) QueueInspect(queue string) (map[string]int, error) {
	var inspect = make(map[string]int, 2)

	rmq.mutex.RLock()
	channel, exists := rmq.channels[queue]
	rmq.mutex.RUnlock()

	if exists {
		q, err := channel.QueueInspect(queue)
		if err != nil {
			return nil, fmt.Errorf(libErrs.ErrorQueueInspectFailure.Error(), queue, err)
		}

		inspect["Total Consumers"] = q.Consumers
		inspect["Total Messages"] = q.Messages

		return inspect, nil
	}

	return nil, fmt.Errorf(libErrs.ErrorQueueInspectMissingQueueFailure.Error(), queue)
}

func (rmq *RMQClient) sendToDeadletterQueue(msg amqp.Delivery, eventErr error) error {
	deadletterQ := fmt.Sprintf("%s-%s", "deadletter", config.ConsumerQueueName)

	rmq.mutex.RLock()
	channel, exists := rmq.channels[deadletterQ]
	rmq.mutex.RUnlock()

	// channel doesn't exist yet for DL - create one
	if !exists && channel == nil {
		channel, err := rmq.CreateNewChannel(deadletterQ)
		if err != nil {
			return err
		}

		errQ := rmq.declareQueue(channel, deadletterQ, amqp.ExchangeDirect)
		if errQ != nil {
			return errQ
		}
	}

	// define error
	eyewaEventErr := base.EyewaEventError{
		Event:        string(msg.Body),
		ErrorMessage: eventErr.Error(),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	errJSON, err := json.Marshal(eyewaEventErr)
	if err != nil {
		return err
	}

	rmq.mutex.RLock()
	channel, exists = rmq.channels[deadletterQ]
	rmq.mutex.RUnlock()

	// publish event error to DL exchange
	if exists && channel != nil {
		err = channel.Publish("", deadletterQ, false, false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         errJSON,
				DeliveryMode: amqp.Persistent,
			})
		if err != nil {
			log.Error(libErrs.ErrorFailedToPublishToDeadletter.Error(),
				zap.String("event", string(errJSON)),
				zap.String("deadletter_queue", deadletterQ), zap.Error(err))

			return err
		}
	}

	return nil
}

func getNameForChannel(queue string) string {
	if config.HostName != "" {
		return config.HostName
	}

	if config.ServiceName == "" {
		return fmt.Sprintf("%s-%d", queue, rand.Uint64())
	}

	return fmt.Sprintf("%s-%d", config.ServiceName, rand.Uint64())
}

// ConnectionListener simply listens for a closed connection and
// simply logs it for visibility. The `Consume` func is responsible
// for notifying a consumer via a callback on lost connetion.
func (rmq *RMQClient) ConnectionListener() {
	go func() {
		notify := rmq.connection.NotifyClose(make(chan *amqp.Error))
		for range notify {
			log.Warn("RMQ connection has closed!")
			return
		}
	}()
}

// IsConnectionOpen gets the connection status to RMQ
func (rmq *RMQClient) IsConnectionOpen() bool {
	return !rmq.connection.IsClosed()
}
