package grpc

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"

	serverConfig "github.com/ramil063/gometrics/cmd/server/config"
)

// ServerConfigFlags содержит переменные флагов
// Address основной урл на которым поднят сервис
// StoreInterval на какой промежуток сохраняются данные
// FileStoragePath путь для сохранения данных
// Restore флаг восстановления данных с сохраненного файла
// DatabaseDSN настройки подключения к БД
// HashKey ключ для декодирования зашифрованных данных
// CryptoKey путь до приватного ключа шифрования
// TrustedSubnet доверенная подсеть для пропуска на сервер
type ServerConfigFlags struct {
	Address         string `env:"GRPC_ADDRESS"`
	FileStoragePath string `env:"GRPC_FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"GRPC_DATABASE_DSN"`
	HashKey         string `env:"GRPC_KEY"`
	CryptoKey       string `env:"GRPC_CRYPTO_KEY"`
	TrustedSubnet   string `env:"GRPC_TRUSTED_SUBNET"`
	StoreInterval   int    `env:"GRPC_STORE_INTERVAL"`
	Restore         bool   `env:"GRPC_RESTORE"`
}

// GetFlags парсит глобальные переменные системы, или парсит флаги, или подменяет их значениями по умолчанию
func GetFlags(config *serverConfig.ServerConfig) (*ServerConfigFlags, error) {
	//значения флагов по умолчанию
	flags := &ServerConfigFlags{
		Address:         "localhost:3202",
		FileStoragePath: "internal/storage/files/grpc/metrics.json",
		StoreInterval:   300,
	}

	var (
		address         string
		fileStoragePath string
		databaseDSN     string
		hashKey         string
		cryptoKey       string
		trustedSubnet   string
		storeInterval   int
		restore         bool
	)

	flag.StringVar(&address, "grpc-a", config.GetAddress(flags.Address), "address and port to run server")
	flag.StringVar(&fileStoragePath, "grpc-f", config.GetFileStoragePath(flags.FileStoragePath), "file storage path")
	flag.StringVar(&databaseDSN, "grpc-d", config.GetDatabaseDSN(flags.DatabaseDSN), "database DSN")
	flag.StringVar(&hashKey, "grpc-k", config.GetHashKey(flags.HashKey), "key for hash")
	flag.StringVar(&cryptoKey, "grpc-crypto-key", config.GetCryptoKey(flags.CryptoKey), "key for encryption")
	flag.StringVar(&trustedSubnet, "grpc-t", config.GetTrustedSubnet(flags.TrustedSubnet), "allowed subnet")
	flag.IntVar(&storeInterval, "grpc-i", config.GetStoreInterval(flags.StoreInterval), "interval of saving metrics to file")
	flag.BoolVar(&restore, "grpc-r", config.GetRestore(flags.Restore), "restore from file")
	flag.Parse()

	var envVars ServerConfigFlags
	err := env.Parse(&envVars)
	if err != nil {
		return flags, fmt.Errorf("error parsing environment variables: %w", err)
	}

	applyFlags(flags, address, fileStoragePath, databaseDSN, hashKey, cryptoKey, trustedSubnet, storeInterval, restore)
	applyEnvVars(flags, envVars)

	return flags, nil
}

// applyFlags присваивание флагов переданных в командной строке
func applyFlags(
	flags *ServerConfigFlags,
	address string,
	fileStoragePath string,
	databaseDSN string,
	hashKey string,
	cryptoKey string,
	trustedSubnet string,
	storeInterval int,
	restore bool,
) {

	if address != "" && address != flags.Address {
		flags.Address = address
	}
	if fileStoragePath != "" && fileStoragePath != flags.FileStoragePath {
		flags.FileStoragePath = fileStoragePath
	}
	if databaseDSN != "" && databaseDSN != flags.DatabaseDSN {
		flags.DatabaseDSN = databaseDSN
	}
	if hashKey != "" && hashKey != flags.HashKey {
		flags.HashKey = hashKey
	}
	if cryptoKey != "" && cryptoKey != flags.CryptoKey {
		flags.CryptoKey = cryptoKey
	}
	if trustedSubnet != "" && trustedSubnet != flags.TrustedSubnet {
		flags.TrustedSubnet = trustedSubnet
	}
	if storeInterval != 0 && storeInterval != flags.StoreInterval {
		flags.StoreInterval = storeInterval
	}
	if restore && restore != flags.Restore {
		flags.Restore = restore
	}
}

// applyEnvVars присваивание переменных окружения
func applyEnvVars(flags *ServerConfigFlags, envVars ServerConfigFlags) {
	if envVars.Address != "" {
		flags.Address = envVars.Address
	}
	if envVars.FileStoragePath != "" {
		flags.FileStoragePath = envVars.FileStoragePath
	}
	if envVars.DatabaseDSN != "" {
		flags.DatabaseDSN = envVars.DatabaseDSN
	}
	if envVars.HashKey != "" {
		flags.HashKey = envVars.HashKey
	}
	if envVars.CryptoKey != "" {
		flags.CryptoKey = envVars.CryptoKey
	}
	if envVars.TrustedSubnet != "" {
		flags.TrustedSubnet = envVars.TrustedSubnet
	}
	if envVars.StoreInterval != 0 {
		flags.StoreInterval = envVars.StoreInterval
	}
	if envVars.Restore {
		flags.Restore = envVars.Restore
	}
}
