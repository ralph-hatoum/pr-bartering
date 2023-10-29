package main

import (
	"fmt"
	"sync"

	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var NodeStorage float64
var port string

func main() {

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, scores, _ := functions.NodeStartup()

	peers = append(peers, "127.0.0.1")

	// path := "test-data/test.txt"
	fmt.Println(fulfilled_requests)
	fmt.Println(storage_pool, pending_requests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Node started !")
	// functions.Store(path, storage_pool, pending_requests)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores)
	}()

	// Wait for the goroutine to finish.
	wg.Wait()

}
