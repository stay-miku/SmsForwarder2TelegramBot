package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	BotToken string `yaml:"botToken"`
	Receiver string `yaml:"receiver"`

	Port   string `yaml:"port"`
	Secret string `yaml:"secret"`

	Https struct {
		Enable bool   `yaml:"enable"`
		Cert   string `yaml:"cert"`
		Key    string `yaml:"key"`
	} `yaml:"https"`
}

var config Config
var entrypoint string

func init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	entrypoint = fmt.Sprintf("https://api.telegram.org/bot%s/", config.BotToken)
}
