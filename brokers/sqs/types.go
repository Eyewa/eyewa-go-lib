package sqs

// Config for all SQS env vars
type Config struct {
	Region          string `mapstructure:"aws_region"`
	AccessKeyID     string `mapstructure:"aws_access_key_id"`
	SecretAccessKey string `mapstructure:"aws_secret_access_key"`
}

// SQSClient SQS client for implementing the MessageBroker interface and handling all things SQS.
type SQSClient struct{}
