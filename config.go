package main

import (
	"encoding/json"
	"io"
	"os"
)

type Configs []Mailbox

type Mailbox struct {
	Host     string
	Port     int
	Username string
	Password string
	PostHook string
	Box      string
}

func NewConfig(confPath *string) (*Configs, error) {
	file, err := os.Open(*confPath)
	if err != nil {
		return nil, err
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	configs := &Configs{}
	err = json.Unmarshal(fileBytes, configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}
