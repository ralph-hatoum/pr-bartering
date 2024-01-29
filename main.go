package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	fswatcher "bartering/fs-watcher"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var NodeStorage float64
var port = "8081"

func main() {

	msgCounter, _ := 0, 0

	args := os.Args

	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	if len(args) != 2 {
		fmt.Println("Not enough arguments ; use : ./bartering <bootstrap-IP>")
		panic(-1)
	} else if !ipRegex.MatchString(args[1]) {
		fmt.Println("Argument invalid : must be an IP adress")
		panic(-1)
	}

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

	// DOIT TOURNER EN ARRIERE PLAN : LISTENPEERSCONNECT (AWAIT AND ANSWER PEER MESSAGES), FS WATCHER (WATCH FOR NEW FILES),
	//  STORAGE TESTER (TEST FOR STORAGE OF OWN FILES, ASK FOR STORAGE AFTER LEASE EXPIRES,
	// MAINTAIN K COPIES OF DATA), BARTERER (SHOULD WE ASK FOR MORE SPACE?)

	// functions.Store(path, storage_pool, pending_requests)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.

	wg.Add(1)
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}
	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue, &msgCounter)
	}()

	// to_request, err := storagerequests.ElectStorageNodes(scores, 3)
	// utils.ErrorHandler(err)
	// fmt.Println(to_request)

	fswatcher.FsWatcher("./test-data", storage_pool, pending_requests)
	fmt.Println(pending_requests)
	// Wait for the goroutine to finish.
	wg.Wait()

}
