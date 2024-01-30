package main

import (
	"fmt"
	"os"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
	"bartering/utils"
)

var NodeStorage float64
var port = "8083"

func main() {

	msgCounter, _ := 0, 0

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

	DecreaseBehavior, IncreaseBehavior := functions.IncreaseDecreaseBehaviors(config)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.

	wg.Add(1)
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}
	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue, &msgCounter)
	}()

	to_request, err := storagerequests.ElectStorageNodes(scores, 1)
	utils.ErrorHandler(err)

	single_node := to_request[0]

	stoRq := datastructures.StorageRequest{CID: "QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps", FileSize: 5.0}

	storagerequests.RequestStorageFromPeer(single_node, stoRq, "8081", bytesAtPeers, scores, &fulfilled_requests, config.StoragerequestsScoreDecreaseRefusedStoReq)
	fmt.Println(bytesAtPeers)
	fmt.Println(fulfilled_requests)
	// Wait for the goroutine to finish.

	wg.Add(1)
	go func() {
		defer wg.Done()
		storagetesting.PeriodicTests(fulfilled_requests, scores, config.StoragetestingTimerTimeoutSec, port, config.StoragetestingTestingPeriod, DecreaseBehavior, IncreaseBehavior, bytesAtPeers, config.StoragerequestsScoreDecreaseRefusedStoReq)
	}()
	wg.Wait()

}
