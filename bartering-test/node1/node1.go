package main

import (
	"fmt"
	"sync"

	"bartering/functions"

	peersconnect "bartering/peers-connect"

	"bartering/bartering-api"
)

var PORT = "8083"

var NodeStorage = 400000000.0

func main() {

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

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(PORT, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers)
	}()

	err := bartering.InitiateBarter("127.0.0.1", ratiosAtPeers)
	fmt.Println(ratiosAtPeers)

	if err != nil {
		fmt.Println("Bartering request failed")
	}

	wg.Wait()

}
