package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	Username  string `json:"d2emu_username"`
	Token     string `json:"d2emu_token"`
	Port      int    `json:"port"`
	Host      string `json:"host"`
	RateLimit int    `json:"rate_limit"`
}

func LoadConfig(path string) *Config {
	config := Config{}

	if !PathExists(path) {
		log.Println("config path doesn't exist")
		return &config
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)
	return &config
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
