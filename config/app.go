package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type KafkaConfig struct {
	Version       string   `toml:"version"`
	Assignor      string   `toml:"assignor"`
	ConsumerGroup string   `toml:"consumer_group"`
	BrokerList    []string `toml:"broker_list"`
	Topics        []string `toml:"topics"`
}

type ElasticConfig struct {
	AddressList           []string `toml:"address_list"`
	MaxIdleConnsPerHost   uint32   `toml:"max_idle_conns_per_host"`
	ResponseHeaderTimeout uint32   `toml:"response_header_timeout"`
}

type Config struct {
	KafkaConfig   KafkaConfig   `toml:"kafka"`
	ElasticConfig ElasticConfig `toml:"elastic"`
}

func NewConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrNotFound.New(fmt.Sprintf("[NewConfig]%s", err.Error()))
	}

	var c Config

	_, err = toml.Decode(string(data), &c)
	if err != nil {
		return nil, ErrParse.New(fmt.Sprintf("[NewConfig]%s", err.Error()))
	}

	return &c, nil
}
