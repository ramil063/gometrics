package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	agentConfig "github.com/ramil063/gometrics/cmd/agent/config"
	"github.com/ramil063/gometrics/cmd/agent/handlers"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	config, err := agentConfig.GetConfig()
	if err != nil {
		logger.WriteErrorLog(err.Error(), "config")
	}
	handlers.InitFlags(config)

	if handlers.CryptoKey != "" {
		crypto.DefaultEncryptor, err = crypto.NewRSAEncryptor(handlers.CryptoKey)

		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create encryptor")
		}
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctxGrSh, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	r.SendMultipleMetricsJSON(c, -1, ctxGrSh)
	fmt.Println("Server shutdown gracefully")
}
