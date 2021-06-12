package rabbitmq

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

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
		return errors.New("No queues to consume or publish to specified!")
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

	return nil
}

// Consume consumes messages from a queue
func (rmq *RMQClient) Consume(queue string) error {
	return nil
}

// Publish publishes a message to a queue
func (rmq *RMQClient) Publish(queue string) error {
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
	if channel == nil {
		if err := rmq.CreateNewChannel(queue); err != nil {
			return err
		}
	}

	exchType := exchangeTypes[exchangeType]

	// declare queue
	q, err := rmq.channels[queue].QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to declare queue(%s). %s", q.Name, err)
	}

	// declare exchange
	err = rmq.channels[queue].ExchangeDeclare(queue, exchType, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to declare an exchange for queue(%s). %s", q.Name, err)
	}

	// bind them together
	err = rmq.channels[queue].QueueBind(queue, queue, queue, false, nil)
	if err != nil {
		return fmt.Errorf("Failed to bind exchange to queue(%s). %s", q.Name, err)
	}

	return nil
}

func (rmq *RMQClient) createConsumerChannel() error {
	if config.ConsumerQueueName == "" {
		return errors.New("No queue specified to consume from!")
	}

	if _, ok := rmq.channels[config.ConsumerQueueName]; !ok {
		conCh, err := rmq.connection.Channel()
		if err != nil {
			return err
		}

		rmq.channels[config.ConsumerQueueName] = conCh

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
	}

	return nil
}

func (rmq *RMQClient) createPublisherChannel() error {
	if config.PublisherQueueName == "" {
		return errors.New("No queue specified to publish to!")
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
func (rmq *RMQClient) CreateNewChannel(queue string) error {
	if rmq.connection != nil {
		rmq.mutex.Lock()
		defer rmq.mutex.Unlock()

		channel, err := rmq.connection.Channel()
		if err != nil {
			return fmt.Errorf("Cannot create new channel for queue(%s). %s", queue, err)
		}

		rmq.channels[queue] = channel
	}

	return nil
}
