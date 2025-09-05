package config

import (
	"github.com/spf13/viper"
)

type Server struct {
	ConfigPath string
	// logger
}

type Config struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func New(configPath string) (*Server, error) {

	return &Server{
		ConfigPath: configPath,
	}, nil
}

func (l *Server) Loader() (*Config, error) {

	config := &Config{}
	err := loadYaml(l.ConfigPath, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadYaml(path string, config interface{}) error {

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return err
	}

	return nil
}
