package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from config.json
// You can change the implementation, for example load from env file, consul, etcd, etc
// NewViper initializes a new Viper instance, sets the configuration file name and type,
// adds configuration paths, and reads the configuration file. If the configuration file
// is not found, it prints a message and uses the default configuration. If there is an
// error reading the configuration file, it panics with a fatal error message.
//
// Returns:
//
//	*viper.Viper: A pointer to the initialized Viper instance.
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath(".")
	config.AddConfigPath("..")
	config.AddConfigPath("/")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return config
}
