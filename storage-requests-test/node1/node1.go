package main

import (
	"fmt"
	"sync"

	configextractor "bartering/config-extractor"
	"bartering/functions"
)

var NodeStorage float64
var port = "8083"

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

	wg.Add(1)

	// go func() {
	// 	defer wg.Done()
	// 	peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio)
	// }()

	// stoRq := datastructures.StorageRequest{CID: "QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps", FileSize: 5.0}

	// storagerequests.RequestStorageFromPeer("127.0.0.1", stoRq, "8084", bytesAtPeers, scores, &fulfilled_requests, config.StoragerequestsScoreDecreaseRefusedStoReq)
	// fmt.Println(bytesAtPeers)
	// fmt.Println(scores)
	// fmt.Println("fulfilled requests :", fulfilled_requests)
	// // Wait for the goroutine to finish.
	wg.Wait()

}
