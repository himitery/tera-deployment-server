package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	conf *Config
	once sync.Once

	errorMessage string
)

func Load() *Config {
	profile := strings.ToLower(os.Getenv("PROFILE"))

	filename := "config/application.yaml"
	if profile != "" {
		filename = fmt.Sprintf("config/application-%s.yaml", profile)
	}

	once.Do(func() {
		data, err := os.ReadFile(filename)
		if err != nil {
			errorMessage = "can't find config file"
		}

		if err == nil && yaml.Unmarshal(data, &conf) != nil {
			errorMessage = "can't parse config file"
		}
	})

	if errorMessage != "" {
		panic(errorMessage)
	}

	conf.Profile = profile

	return conf
}
