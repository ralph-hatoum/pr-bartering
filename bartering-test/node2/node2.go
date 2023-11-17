package main

import (
	"fmt"
	"sync"

	"bartering/bartering-api"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
)

var PORT = "8082"

var NodeStorage = 400000000.0

func main() {

	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, scores, ratios := functions.NodeStartup()
	fmt.Println(ratios)
	peers = append(peers, "127.0.0.1")

	// path := "test-data/test.txt"
	fmt.Println(fulfilled_requests)
	fmt.Println(storage_pool, pending_requests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Node started !")
	// functions.Store(path, storage_pool, pending_requests)

	bytesAtPeers = append(bytesAtPeers, bartering.PeerStorageUse{NodeIP: "127.0.0.1", StorageAtNode: 4000.0})
	scores = append(scores, bartering.NodeScore{NodeIP: "127.0.0.1", Score: 1.0})

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(PORT, NodeStorage, bytesAtPeers, scores, ratios)
	}()

	// Wait for the goroutine to finish.
	wg.Wait()

}
