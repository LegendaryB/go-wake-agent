package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ApplicationConfiguration struct {
	ListenPort int
}

type Configuration struct {
	Application    ApplicationConfiguration
	AllowedClients []*AllowedClient `toml:"allowed_clients"`
}

const configurationFileName string = "config.toml"

func NewConfiguration() (*Configuration, error) {
	var configuration Configuration

	if _, err := toml.DecodeFile(configurationFileName, &configuration); err != nil {
		fmt.Print("Failed to load configuration!")
		return nil, err
	}

	return &configuration, nil
}
