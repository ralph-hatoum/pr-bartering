package functions

import (
	"fmt"
	"strings"

	api_ipfs "../api-ipfs"

	bootstrapconnect "../bootstrap-connect"

	"../utils"
)

type StorageRequest struct {
	CID      string
	fileSize float64
}

func NodeStartup() ([]string, []StorageRequest, []StorageRequest, []string) {

	fmt.Println("Starting node")
	// Create all needed data structures
	fmt.Println("Creating storage pool and requests lists")
	storage_pool, pending_requests, fulfilled_requests := createStorageRequestsLists()
	fmt.Println("Storage pool and requests lists created successfully")
	// Connect to bootstrap

	//Create peers list
	fmt.Println("Creating peers list")
	peers := bootstrapconnect.GetPeersFromBootstrapHTTP("127.0.0.1", "8080")

	return storage_pool, pending_requests, fulfilled_requests, peers
}

func Store(path string, storage_pool []string, pending_requests []StorageRequest) {
	// Function called to store a file on the network

	// Uploading file to IPFS & retrieving its CID
	upload_command_result := api_ipfs.UploadToIPFS(path)
	CID := strings.Split(upload_command_result, " ")[1]

	// Add the CID to the storage pool
	storage_pool = append(storage_pool, CID)

	fmt.Println(storage_pool)

	// Pin file to IPFS
	//pin_command_result := api_ipfs.PinToIPFS(CID)

	file_size := utils.GetFileSize(path)

	fmt.Println(file_size)

	storage_request := StorageRequest{CID, file_size}

	pending_requests = append(pending_requests, storage_request)

	fmt.Println("Pending requests : ", pending_requests)
	// TODO build storage request
	// TODO propagate to network
}

func createStorageRequestsLists() ([]string, []StorageRequest, []StorageRequest) {
	storage_pool := []string{}

	pending_requests := []StorageRequest{}

	fulfilled_requests := []StorageRequest{}

	return storage_pool, pending_requests, fulfilled_requests

}
