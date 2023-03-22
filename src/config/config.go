package config

import (
	"io/ioutil"
	"log"
	"os"
	"proxauth/login"
	"proxauth/rule"
	"strconv"
	"time"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Users                 []login.User  `json:"users" yaml:"users"`
	Rules                 []rule.Rule   `json:"rules" yaml:"rules"`
	ServerSecret          []byte        `json:"serverSecret yaml:"serverSecret"`
	Port                  int           `json:"port" yaml:"port"`
	JWTExpirationDuration time.Duration `json:"jwtExpirationDuration" yaml:"jwtExpirationDuration"`
}

func getEnv() (string, string, string, int, time.Duration) {
	config := os.Getenv("CONFIG")
	if len(config) != 0 {
		log.Println("INFO: ENV CONFIG is set. Ignore ENV CONFIG_FILE.")
	}

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

	jwtExpirationDuration, err := time.ParseDuration(os.Getenv("JWT_EXPIRATION_DURATION"))
	if err != nil {
		log.Println("WARNING: ENV JWT_EXPIRATION_DURATION is not set. Use default \"24h\".")
		jwtExpirationDuration, _ = time.ParseDuration("24h")
	}

	return config, configFile, serverSecret, port, jwtExpirationDuration
}

func Load() (*Config, error) {
	configString, configFile, serverSecret, port, jwtExpirationDuration := getEnv()

	var bytes []byte
	if len(configString) > 0 {
		bytes = []byte(configString)
	} else {
		file, err := os.Open(configFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		bytes, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}

	var config Config
	err := yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	config.ServerSecret = []byte(serverSecret)
	config.Port = port
	config.JWTExpirationDuration = jwtExpirationDuration

	return &config, nil
}
