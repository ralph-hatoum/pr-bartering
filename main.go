package main

import (
	"fmt"
	"os"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	"bartering/dumper"
	fswatcher "bartering/fs-watcher"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
)

func main() {

	msgCounter, _ := 0, 0

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Not enough arguments ; use : ./bartering <bootstrap-IP>")
		panic(-1)
	}

	bootstrapIp := args[1]

	config := configextractor.ConfigExtractor("config.yaml")

	port := fmt.Sprint(config.Port)
	NodeStorage := config.TotalStorage

	configextractor.ConfigPrinter(config)

	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, bytesForPeers, scores, ratiosAtPeers, ratiosForPeers, storedForPeers := functions.NodeStartup(bootstrapIp)

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

	DecreaseBehavior, IncreaseBehavior := functions.IncreaseDecreaseBehaviors(config)

	var wg sync.WaitGroup

	// Message queues
	storageRequestsChannel := make(chan datastructures.StorageRequestQueueMessage)
	newFilesChannel := make(chan datastructures.StorageRequest)
	testRequestsChannel := make(chan datastructures.TestRequestQueueMessage)

	wg.Add(1)
	go func() {
		// DUMPER - to have access to the state of all data structures (for testing purposes)
		defer wg.Done()
		datastructures := []dumper.Datastructure{
			{Name: "bytesAtPeers", Value: bytesAtPeers},
			{Name: "bytesForPeers", Value: bytesForPeers},
			{Name: "fulfilledRequests", Value: fulfilled_requests},
			{Name: "storagePool", Value: storage_pool},
			{Name: "pendingRequests", Value: pending_requests},
			{Name: "peers", Value: peers},
			{Name: "scores", Value: scores},
			{Name: "ratiosForPeers", Value: ratiosForPeers},
			{Name: "ratiosAtPeers", Value: ratiosAtPeers},
			{Name: "storedForPeers", Value: storedForPeers},
		}
		dumper.Dumper(datastructures)
	}()

	wg.Add(1)
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}
	go func() {
		// PEER LISTENER - to receive messages from other peers
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue, &msgCounter, storageRequestsChannel, testRequestsChannel)
	}()

	wg.Add(1)
	go func() {
		// Handle storage requests pool
		defer wg.Done()
		storagerequests.HandleStorageRequest(bytesForPeers, &storedForPeers, storageRequestsChannel)
	}()

	wg.Add(1)
	go func() {
		// STORAGE TESTING - to test storage at peers
		defer wg.Done()
		storagetesting.PeriodicTests(&fulfilled_requests, scores, config.StoragetestingTimerTimeoutSec, port, config.StoragetestingTestingPeriod, DecreaseBehavior, IncreaseBehavior, bytesAtPeers, config.StoragerequestsScoreDecreaseRefusedStoReq, newFilesChannel)
	}()

	wg.Add(1)
	go func() {
		// Handle new files pool
		defer wg.Done()
		storagerequests.StoreKCopiesOnNetwork(scores, 2, port, bytesAtPeers, &fulfilled_requests, config.StoragerequestsScoreDecreaseRefusedStoReq, newFilesChannel)
	}()

	wg.Add(1)
	go func() {
		// FSWATCHER - to upload data on network
		defer wg.Done()
		fswatcher.FsWatcher("./data", scores, config.DataCopies, port, bytesAtPeers, &fulfilled_requests, config.StoragerequestsScoreDecreaseRefusedStoReq, newFilesChannel)
	}()

	wg.Add(1)
	go func() {
		// Test Handler pool
		defer wg.Done()
		storagetesting.HandleTest(testRequestsChannel)
	}()

	fmt.Println("Node started ! Listening on port ", port)

	// TODO : BARTERER, FAILURESIM, DATASIM
	wg.Wait()

	// to_request, err := storagerequests.ElectStorageNodes(scores, 1)
	// utils.ErrorHandler(err)

	// single_node := to_request[0]

	// stoRq := datastructures.StorageRequest{CID: "QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps", FileSize: 5.0}

	// storagerequests.RequestStorageFromPeer(single_node, stoRq, "8081", bytesAtPeers, scores, &fulfilled_requests, config.StoragerequestsScoreDecreaseRefusedStoReq)
	// fmt.Println(bytesAtPeers)
	// fmt.Println(fulfilled_requests)
	// Wait for the goroutine to finish.

}
