package main

import "./functions"

func main() {

	// storage_pool, pending_requests, fulfilled_storage, peers := functions.NodeStartup()
	storage_pool, _, _, _ := functions.NodeStartup()

	path := "test-data/test.txt"
	functions.Store(path, storage_pool)
}
