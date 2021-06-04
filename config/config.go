package config

import (
	"strings"

	"github.com/eyewa/eyewa-go-lib/brokers"
	"github.com/ory/viper"
)

var (
	config  configuration
	envVars = []string{
		"ENV",
		"LOG_LEVEL",
		"MESSAGE_BROKER",
	}
)

type configuration struct {
	Env           string
	LogLevel      string             `mapstructure:"log_level"`
	MessageBroker brokers.BrokerType `mapstructure:"message_broker"`
}

func initConfig() error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	for _, v := range envVars {
		if err := viper.BindEnv(v); err != nil {
			return err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}

// GetConfigEnvVars returns all the env vars used/supported by pkg
func GetConfigEnvVars() []string {
	return envVars
}
