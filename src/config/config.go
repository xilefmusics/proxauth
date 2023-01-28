package config

import (
	"io/ioutil"
	"log"
	"os"
	"proxauth/login"
	"proxauth/rule"
	"strconv"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Users        []login.User `json:"users" yaml:"users"`
	Rules        []rule.Rule  `json:"rules" yaml:"rules"`
	ServerSecret []byte       `json:"serverSecret yaml:"serverSecret"`
	Port         int          `json:"port" yaml:"port"`
}

func getEnv() (string, string, int) {
	configFile := os.Getenv("CONFIG_FILE")
	if len(configFile) == 0 {
		log.Println("WARNING: ENV CONFIG_FILE is not set. Use default \"../config/config.yaml\".")
		configFile = "../config/config.yaml"
	}

	serverSecret := os.Getenv("SERVER_SECRET")
	if len(serverSecret) == 0 {
		log.Println("WARNING: ENV SERVER_SECRET is not set. Use default \"changeMe\".")
		serverSecret = "changeMe"
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("WARNING: ENV PORT is not set. Use default \"8081\".")
		port = 8081
	}

	return configFile, serverSecret, port
}

func Load() (*Config, error) {
	configFile, serverSecret, port := getEnv()

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	config.ServerSecret = []byte(serverSecret)
	config.Port = port

	return &config, nil
}
