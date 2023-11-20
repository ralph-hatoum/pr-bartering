package main

import (
	"fmt"
	"sync"

	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var NodeStorage float64
var port = "8084"

func main() {

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, scores, ratios := functions.NodeStartup()

	// path := "test-data/test.txt"
	fmt.Println("Bytes at peers :", bytesAtPeers)
	fmt.Println("Fulfilled requests : ", fulfilled_requests)
	fmt.Println("Storage pool : ", storage_pool)
	fmt.Println("Pending requests : ", pending_requests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Scores : ", scores)
	fmt.Println("Node ratios : ", ratios)
	fmt.Println("")
	fmt.Println("Node started ! Listening on port ", port)

	// functions.Store(path, storage_pool, pending_requests)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratios)
	}()

	// Wait for the goroutine to finish.
	wg.Wait()

}
