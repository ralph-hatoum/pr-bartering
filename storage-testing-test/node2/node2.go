package main

import (
	"fmt"
	"os"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var NodeStorage float64
var port = "8081"

func main() {

	args := os.Args

	bootstrapIp := args[1]

	fmt.Println("Extracting configuration")
	config := configextractor.ConfigExtractor("config.yaml")

	configextractor.ConfigPrinter(config)

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, bytesForPeers, scores, ratiosAtPeers, ratiosForPeers, storedForPeers := functions.NodeStartup(bootstrapIp)

	// path := "test-data/test.txt"
	fmt.Println("Bytes at peers :", bytesAtPeers)
	fmt.Println("Bytes stored for peers : ", bytesForPeers)
	fmt.Println("Fulfilled requests : ", fulfilled_requests)
	fmt.Println("Storage pool : ", storage_pool)
	fmt.Println("Pending requests : ", pending_requests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Scores : ", scores)
	fmt.Println("Node ratios : ", ratiosForPeers)
	fmt.Println("ratios at peers : ", ratiosAtPeers)
	fmt.Println("stored for peers : ", storedForPeers)
	fmt.Println("")
	fmt.Println("Node started ! Listening on port ", port)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}
	wg.Add(1)

	// Message queue
	storageRequestsChannel := make(chan datastructures.StorageRequestQueueMessage)
	testRequestsChannel := make(chan datastructures.TestRequestQueueMessage)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue, storageRequestsChannel, testRequestsChannel)
	}()

	// Wait for the goroutine to finish.
	wg.Wait()

}
