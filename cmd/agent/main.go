package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	agentConfig "github.com/ramil063/gometrics/cmd/agent/config"
	"github.com/ramil063/gometrics/cmd/agent/handlers"
	"github.com/ramil063/gometrics/cmd/agent/handlers/grpc"
	"github.com/ramil063/gometrics/internal/constants"
	"github.com/ramil063/gometrics/internal/logger"
	"github.com/ramil063/gometrics/internal/security/crypto"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	params := agentConfig.NewConfigParams(
		constants.ConfigHTTPConsoleShortKey,
		constants.ConfigHTTPConsoleFullKey,
		constants.ConfigHTTPTypeAlias)
	config, err := agentConfig.GetConfig(params)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "config")
	}

	flags, err := handlers.GetFlags(config)
	if err != nil {
		logger.WriteErrorLog(err.Error(), "flags")
	}

	manager := crypto.NewCryptoManager()
	if flags != nil && flags.CryptoKey != "" {
		encryptor, err := crypto.NewRSAEncryptor(flags.CryptoKey)

		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create encryptor")
		}
		manager.SetDefaultEncryptor(encryptor)
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctxGrSh, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	var serversWg sync.WaitGroup
	serversWg.Add(1)

	go grpc.StartClient(ctxGrSh, &serversWg)

	serversWg.Add(1)
	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	go r.SendMultipleMetricsJSON(c, -1, ctxGrSh, flags, manager, &serversWg)

	serversWg.Wait()
	fmt.Println("Server shutdown gracefully")
}
