package main

import (
	"fmt"

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
	handlers.ParseFlags()

	if handlers.CryptoKey != "" {
		var err error
		crypto.DefaultEncryptor, err = crypto.NewRSAEncryptor(handlers.CryptoKey)

		if err != nil {
			logger.WriteErrorLog(err.Error(), "Failed to create encryptor")
		}
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	c := handlers.NewJSONClient()
	r := handlers.NewRequest()
	r.SendMultipleMetricsJSON(c, -1)
}
