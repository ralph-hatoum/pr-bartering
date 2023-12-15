package main

import (
	"fmt"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	failuresimulation "bartering/failure-simulation"
	fswatcher "bartering/fs-watcher"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var NodeStorage float64
var port = "8081"

func main() {

	fmt.Println("Extracting configuration")
	config := configextractor.ConfigExtractor("config.yaml")

	configextractor.ConfigPrinter(config)

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, bytesForPeers, scores, ratiosAtPeers, ratiosForPeers, storedForPeers := functions.NodeStartup()

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

	var failureMutex sync.Mutex

	go failuresimulation.Failure(config, 2.0, 100, &failureMutex)

	wg.Add(1)
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}
	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue)
	}()

	// to_request, err := storagerequests.ElectStorageNodes(scores, 3)
	// utils.ErrorHandler(err)
	// fmt.Println(to_request)

	fswatcher.FsWatcher("./test-data", storage_pool, pending_requests)
	fmt.Println(pending_requests)
	// Wait for the goroutine to finish.
	wg.Wait()

}
