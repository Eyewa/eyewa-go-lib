package kafka

// EnvConfig for all Kafka env vars
type EnvConfig struct {
	Brokers           []string `env:"KAFKA_BROKERS,required" envSeparator:","`
	NumPartitions     int      `env:"KAFKA_NUM_PARTITIONS,required"`
	ReplicationFactor int      `env:"KAFKA_REPLICATION_FACTOR,required" envDefault:"1"`
}

// KafkaClient Kafka client for implementing the MessageBroker interface and handling all things Kafka.
type KafkaClient struct{}
