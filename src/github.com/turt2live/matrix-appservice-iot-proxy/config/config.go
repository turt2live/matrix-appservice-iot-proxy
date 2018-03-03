package config

import (
	"sync"
	"os"
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type ProxyConfig struct {
	HomeserverUrl string   `yaml:"homeserverUrl"`
	AllowedTokens []string `yaml:"allowedTokens,flow"`
	LogDirectory  string   `yaml:"logDirectory"`
	BindAddress   string   `yaml:"bindAddress"`
	BindPort      int      `yaml:"bindPort"`
}

var instance *ProxyConfig
var lock = &sync.Once{}
var Path = "iot-proxy.yaml"

func ReloadConfig() (error) {
	c := NewDefaultConfig()

	// Write a default config if the one given doesn't exist
	_, err := os.Stat(Path)
	exists := err == nil || !os.IsNotExist(err)
	if !exists {
		fmt.Println("Generating new configuration...")
		configBytes, err := yaml.Marshal(c)
		if err != nil {
			return err
		}

		newFile, err := os.Create(Path)
		if err != nil {
			return err
		}

		_, err = newFile.Write(configBytes)
		if err != nil {
			return err
		}

		err = newFile.Close()
		if err != nil {
			return err
		}
	}

	f, err := os.Open(Path)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer, err := ioutil.ReadAll(f)
	err = yaml.Unmarshal(buffer, &c)
	if err != nil {
		return err
	}

	instance = c
	return nil
}

func Get() (*ProxyConfig) {
	if instance == nil {
		lock.Do(func() {
			err := ReloadConfig()
			if err != nil {
				panic(err)
			}
		})
	}
	return instance
}

func NewDefaultConfig() *ProxyConfig {
	return &ProxyConfig{
		HomeserverUrl: "http://localhost:8008",
		AllowedTokens: make([]string, 0),
		LogDirectory:  "logs",
		BindAddress:   "0.0.0.0",
		BindPort:      4232,
	}
}
