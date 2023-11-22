package main

import (
	"fmt"
	"sync"
	"time"

	"bartering/functions"
	peersconnect "bartering/peers-connect"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
)

var NodeStorage float64
var port = "8084"

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
	fmt.Println("stored for peers : ", storedForPeers)
	fmt.Println("")
	fmt.Println("Node started ! Listening on port ", port)

	// functions.Store(path, storage_pool, pending_requests)

	var wg sync.WaitGroup // Import "sync" package to use WaitGroup.

	wg.Add(1)

	go func() {
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers)
	}()

	storage_request := storagerequests.StorageRequest{CID: "QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps", FileSize: 5.5}

	storagerequests.RequestStorageFromPeer("127.0.0.1", storage_request, "8081", bytesAtPeers, scores, &fulfilled_requests)

	fmt.Println("fulfilled requests : ", fulfilled_requests)
	fmt.Println("bytesAtPeers : ", bytesAtPeers)

	fmt.Println("Waiting 10 secs to request proof ...")
	time.Sleep(10 * time.Second)

	storagetesting.ContactPeerForTest(storage_request.CID, "127.0.0.1", scores)

	fmt.Println("scores after test : ", scores)

	// Wait for the goroutine to finish.
	wg.Wait()

}
