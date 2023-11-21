package main

import (
	"fmt"
	"sync"

	"bartering/functions"

	peersconnect "bartering/peers-connect"

	"bartering/bartering-api"
)

var PORT = "8081"

var NodeStorage = 400000000.0

func main() {

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
	fmt.Println("")
	fmt.Println("Node started ! Listening on port ", PORT)

	var wg sync.WaitGroup

	// Adding 127.0.0.1 to all lists for the barter test
	bytesAtPeers = append(bytesAtPeers, bartering.PeerStorageUse{NodeIP: "127.0.0.1", StorageAtNode: 4000.0})
	scores = append(scores, bartering.NodeScore{NodeIP: "127.0.0.1", Score: 100.0})
	ratiosAtPeers = append(ratiosAtPeers, bartering.NodeRatio{NodeIP: "127.0.0.1", Ratio: 1.0})

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(PORT, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers)
	}()

	err := bartering.InitiateBarter("127.0.0.1", ratiosAtPeers)

	if err != nil {
		fmt.Println("Bartering request failed")
	}

	// Wait for the goroutine to finish.
	wg.Wait()

}
