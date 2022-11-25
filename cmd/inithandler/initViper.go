package inithandler

import (
	"log"

	"github.com/invoice-service/spec"

	"github.com/spf13/viper"
)

// use viper package to read .env file
// return the value of the key
func InitViper() {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	viper.AddConfigPath(".")
	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config spec.Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath("/app")
	viper.AddConfigPath("/app/env") //for use in docker
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
