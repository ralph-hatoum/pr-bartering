package main

import (
	"go.uber.org/zap"

	configextractor "bartering/config-extractor"
	peersconnect "bartering/peers-connect"
	worker "bartering/worker"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Starting ...")

	config := configextractor.ConfigExtractor("config.yaml")

	// TODO: Initiate needed storage structures : sqlite ? redis ? dicedb ? in memory + regular flush to disk ?

	workloadChannel := make(chan worker.Workload)

	for i := 0; i <= 5; i++ {
		go worker.Worker(workloadChannel, logger)
	}

	peersconnect.ListenPeersRequestsTCP(config.Port, workloadChannel, logger)
}
