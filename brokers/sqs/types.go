package sqs

// EnvConfig for all SQS env vars
type EnvConfig struct {
	Region    string `env:"AWS_REGION,required"`
	AccessKey string `env:"AWS_ACCESS_KEY_ID,required"`
	Secret    string `env:"AWS_SECRET_ACCESS_KEY,required"`
}

// SQSClient SQS client for implementing the MessageBroker interface and handling all things SQS.
type SQSClient struct{}
