package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port  int   `yaml:"port"`
	Task  Task  `yaml:"task"`
	Redis Redis `yaml:"redis"`
}

type Task struct {
	Concurrent uint16 `yaml:"concurrent"`
	Queues     struct {
		Critical uint16 `yaml:"critical"`
		Default  uint16 `yaml:"default"`
		Low      uint16 `yaml:"low"`
	} `yaml:"queues"`
}

type Redis struct {
	Addr              string `yaml:"addr"`
	DB                uint16 `yaml:"db"`
	Password          string `yaml:"password"`
	RedisTLS          string `yaml:"redis_tls"`
	RedisURL          string `yaml:"redis_url"`
	RedisInsecureTLS  bool   `yaml:"redis_insecure_tls"`
	RedisClusterNodes string `yaml:"redis_cluster_nodes"`
}

// LoadSecret reads the file from path and return Secret
func LoadSecret(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadSecretFromBytes(data)
}

// LoadSecretFromBytes reads the secret file from data bytes
func LoadSecretFromBytes(data []byte) (*Config, error) {
	flag := viper.New()
	flag.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	flag.AutomaticEnv()
	flag.SetEnvPrefix("workerine")
	flag.SetConfigType("yaml")

	if err := flag.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	var cred Config
	err := flag.Unmarshal(&cred)
	if err != nil {
		log.Fatalf("Error loading cred: %v", err)
	}

	return &cred, nil
}
