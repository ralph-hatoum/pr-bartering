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
	storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, scores, ratios := functions.NodeStartup()
	fmt.Println(ratios)
	peers = append(peers, "127.0.0.1")

	// path := "test-data/test.txt"
	fmt.Println(fulfilled_requests)
	fmt.Println(storage_pool, pending_requests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Node started !")
	// functions.Store(path, storage_pool, pending_requests)

	var wg sync.WaitGroup

	// Adding 127.0.0.1 to all lists for the barter test
	bytesAtPeers = append(bytesAtPeers, bartering.PeerStorageUse{NodeIP: "127.0.0.1", StorageAtNode: 4000.0})
	scores = append(scores, bartering.NodeScore{NodeIP: "127.0.0.1", Score: 100.0})
	ratios = append(ratios, bartering.NodeRatio{NodeIP: "127.0.0.1", Ratio: 1.0})

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(PORT, NodeStorage, bytesAtPeers, scores)
	}()

	err := bartering.InitiateBarter("127.0.0.1", ratios)

	if err != nil {
		fmt.Println("Bartering request failed")
	}

	// Wait for the goroutine to finish.
	wg.Wait()

}
