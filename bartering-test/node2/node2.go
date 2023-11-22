package main

import (
	"fmt"
	"sync"

	"bartering/bartering-api"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var PORT = "8084"

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

	bytesAtPeers = append(bytesAtPeers, bartering.PeerStorageUse{NodeIP: "127.0.0.1", StorageAtNode: 4000.0})
	scores = append(scores, bartering.NodeScore{NodeIP: "127.0.0.1", Score: 1.0})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(PORT, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers)
	}()

	// Wait for the goroutine to finish.
	wg.Wait()

}
