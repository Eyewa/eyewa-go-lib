package kafka

// EnvConfig for all Kafka env vars
type EnvConfig struct {
	Brokers           []string `mapstructure:"kafka_brokers"`
	NumPartitions     int      `mapstructure:"kafka_num_partitions"`
	ReplicationFactor int      `mapstructure:"kafka_replication_factor"`
}

// KafkaClient Kafka client for implementing the MessageBroker interface and handling all things Kafka.
type KafkaClient struct{}
