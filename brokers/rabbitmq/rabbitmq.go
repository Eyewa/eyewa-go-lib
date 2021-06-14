package rabbitmq

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/eyewa/eyewa-go-lib/base"
	libErrs "github.com/eyewa/eyewa-go-lib/errors"
	"github.com/eyewa/eyewa-go-lib/log"
	"github.com/ory/viper"
	"github.com/streadway/amqp"
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
		mutex:      new(sync.Mutex),
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
	rmq.mutex = new(sync.Mutex)
	rmq.channels = make(map[string]*amqp.Channel)

	rmq.mutex.Lock()
	defer rmq.mutex.Unlock()

	// create channel for consuming (if any)
	if config.ConsumerQueueName != "" {
		for err := range rmq.createConsumerChannel() {
			return err
		}
	}

	// create channel for publish (if any)
	if config.PublisherQueueName != "" {
		if err := rmq.createPublisherChannel(); err != nil {
			return err
		}
	}

	return nil
}

// Consume consumes messages from a queue
func (rmq *RMQClient) Consume(wg *sync.WaitGroup, queue string, errChan chan<- error) {
	defer wg.Done()

	if channel, ok := rmq.channels[queue]; ok {
		if channel != nil {
			log.Info(fmt.Sprintf("Listening to %s for new messages...", queue))

			msgs, err := channel.Consume(queue, queue, false, false, false, false, nil)
			if err != nil {
				errChan <- fmt.Errorf("Failed to consume from queue(%s). %s", queue, err)
			}

			for msg := range msgs {
				var event *base.EyewaEvent
				fmt.Println(string(msg.Body))

				err := json.Unmarshal(msg.Body, &event)
				if err != nil {
					// TODO: add message to deadletter queue
					errChan <- fmt.Errorf("Failed to unmarshaling event from queue(%s). %s", queue, err)
					err = msg.Nack(false, false)
					if err != nil {
						errChan <- err
					}
				} else {
					fmt.Println(string(msg.Body))
					// TODO: save event to db
					err = msg.Ack(false)
					if err != nil {
						errChan <- fmt.Errorf("Failed to acknowledge new messages from queue(%s)", queue)
					}
				}
			}
		}
	}
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string) error {
	// rmq.channels[queue].Publish()
	return nil
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
		return fmt.Errorf("Failed to declare queue(%s). %s", q.Name, err)
	}

	// declare exchange
	err = channel.ExchangeDeclare(exchName, exchType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to declare an exchange for queue(%s). %s", q.Name, err)
	}

	// bind them together
	err = channel.QueueBind(queue, queue, exchName, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to bind exchange to queue(%s). %s", q.Name, err)
	}

	return nil
}

func (rmq *RMQClient) createConsumerChannel() chan error {
	errChan := make(chan error)
	defer close(errChan)

	if config.ConsumerQueueName == "" {
		errChan <- libErrs.ErrorNoConsumerQueueSpecified
		return errChan
	}

	if _, exists := rmq.channels[config.ConsumerQueueName]; !exists {
		conCh, err := rmq.connection.Channel()
		if err != nil {
			errChan <- err
			return errChan
		}

		prefetchCount, _ := strconv.Atoi(config.QueuePrefetchCount)
		if prefetchCount == 0 {
			prefetchCount = defaultPrefetchCount
		}

		if err := conCh.Qos(prefetchCount, 0, true); err != nil {
			errChan <- err
			return errChan
		}

		if err := rmq.declareQueue(conCh, config.ConsumerQueueName, config.ConsumerExchangeType); err != nil {
			errChan <- err
			return errChan
		}

		rmq.channels[config.ConsumerQueueName] = conCh
	}

	return errChan
}

func (rmq *RMQClient) createPublisherChannel() error {
	if config.PublisherQueueName == "" {
		return libErrs.ErrorNoPublisherQueueSpecified
	}

	if _, ok := rmq.channels[config.PublisherQueueName]; !ok {
		pubCh, err := rmq.connection.Channel()
		if err != nil {
			return err
		}

		rmq.channels[config.PublisherQueueName] = pubCh
	}

	return nil
}

// CreateNewChannel creates a new channel for specified queue
func (rmq *RMQClient) CreateNewChannel(queue string) (*amqp.Channel, error) {
	if rmq.connection != nil {
		rmq.mutex.Lock()
		defer rmq.mutex.Unlock()

		channel, err := rmq.connection.Channel()
		if err != nil {
			return nil, fmt.Errorf("Cannot create new channel for queue(%s). %s", queue, err)
		}

		rmq.channels[queue] = channel
		return rmq.channels[queue], nil
	}

	return nil, libErrs.ErrorNoRMQConnection
}

// QueueInspect inspects a queue and returns no. of consumers + messages
// currently unacked.
func (rmq *RMQClient) QueueInspect(queue string) (map[string]int, error) {
	var inspect = make(map[string]int, 2)

	if channel, exists := rmq.channels[queue]; exists {
		q, err := channel.QueueInspect(queue)
		if err != nil {
			return nil, fmt.Errorf("Failed to inspect queue(%s). %s", queue, err)
		}

		inspect["Total Consumers"] = q.Consumers
		inspect["Total Messages"] = q.Messages

		return inspect, nil
	}

	return nil, fmt.Errorf("Queue specified to inspect doesn't exist queue(%s)", queue)
}
