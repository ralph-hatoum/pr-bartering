package functions

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func NodeStartup() ([]string, []string, []string, []string) {

	fmt.Println("Starting node")
	// Create all needed data structures
	fmt.Println("Creating storage pool and requests lists")
	storage_pool, pending_requests, fulfilled_requests := createStorageRequestsLists()
	fmt.Println("Storage pool and requests lists created successfully")
	// Connect to bootstrap

	//Create peers list
	fmt.Println("Creating peers list")
	peers := createPeersList()
	fmt.Println("Peers list created successfully")

	return storage_pool, pending_requests, fulfilled_requests, peers
}

func Store(path string, storage_pool []string) {
	// Function called to store a file on the network

	// Uploading file to IPFS & retrieving its CID
	upload_command_result := uploadToIPFS(path)
	CID := strings.Split(upload_command_result, " ")[1]

	// Add the CID to the storage pool
	storage_pool = append(storage_pool, CID)

	// Pin file to IPFS
	//pin_command_result := pinToIPFS(CID)

	file_size := getFileSize(path)

	fmt.Println(file_size)

	//storage_request := (CID)
	// TODO build storage request
	// TODO propagate to network
}

func createPeersList() []string {

	peersList := []string{}

	return peersList
}

func createStorageRequestsLists() ([]string, []string, []string) {
	storage_pool := []string{}

	pending_requests := []string{}

	fulfilled_requests := []string{}

	return storage_pool, pending_requests, fulfilled_requests

}

func uploadToIPFS(path string) string {
	cmd := "ipfs"
	args := []string{"add", path}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)

}

func errorHandler(err error) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		panic(0)
	}
}

func listPrint(list []string) {
	for _, element := range list {
		fmt.Print(element + " ")
	}
}

func pinToIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "add", cid}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)

}

func unpinIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "rm", cid}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)
}

func getFileSize(path string) float64 {
	// Returns file size in KB

	fileInfo, err := os.Stat(path)
	errorHandler(err)
	fileSize := fileInfo.Size()

	return float64(fileSize) / 1024.0
}
