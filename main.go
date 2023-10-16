package main

import (
	"fmt"

	"./functions"
)

func main() {

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, pending_requests, _, peers := functions.NodeStartup()

	path := "test-data/test.txt"
	fmt.Println("Peers : ", peers)
	functions.Store(path, storage_pool, pending_requests)
}
