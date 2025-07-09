package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

type envConfig struct {
	Config string `env:"CONFIG"`
}

// ServerConfig структура для парсинга файла конфигурации
type ServerConfig struct {
	Restore         *bool  `json:"restore"`
	Address         string `json:"address"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
	HashKey         string `json:"hash_key"`
	CryptoKey       string `json:"crypto_key"`
	StoreInterval   string `json:"store_interval"`
}

// loadConfig загружает конфигурацию из файла
func (cfg *ServerConfig) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return nil
}

// prepareConfig подготавливает параметры конфигурации для дальнейшей работы
func (cfg *ServerConfig) prepareConfig() error {
	storeInterval, err := time.ParseDuration(cfg.StoreInterval)
	if err != nil {
		return err
	}
	cfg.StoreInterval = strconv.FormatFloat(storeInterval.Seconds(), 'f', 0, 64)

	return nil
}

// getConfigName получение названия файла конфигурации
func getConfigName() string {
	// configName имя файла конфигурации
	var configName = ""
	flag.StringVar(&configName, "c", "", "key for configuration")
	flag.StringVar(&configName, "config", "", "key for configuration")

	var ev envConfig
	_ = env.Parse(&ev)

	if ev.Config != "" {
		configName = ev.Config
	}

	return configName
}

// GetConfig установка значений конфигурации
func GetConfig() (*ServerConfig, error) {
	configName := getConfigName()

	var config ServerConfig
	var err error

	if configName == "" {
		return &config, nil
	}

	if err = config.loadConfig(configName); err != nil {
		return nil, err
	}

	if err = config.prepareConfig(); err != nil {
		return nil, err
	}

	return &config, err
}
