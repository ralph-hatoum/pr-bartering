package functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	peers := getPeersFromBootstrapHTTP("127.0.0.1", "8080")

	return storage_pool, pending_requests, fulfilled_requests, peers
}

func Store(path string, storage_pool []string, pending_requests []StorageRequest) {
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

func getPeersFromBootstrapTCP(IP string, port string) {

	serverAddress := IP + ":" + port
	conn, err := net.Dial("tcp", serverAddress)
	errorHandler(err)

	defer conn.Close()

	message := "hello\n"

	_, err = io.WriteString(conn, message)
	errorHandler(err)

	fmt.Println("Called bootstrap, awaiting response")

	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	errorHandler(err)
	fmt.Println(response)

}

func getPeersFromBootstrapHTTP(IP string, port string) []string {
	serverUrl := IP + ":" + port

	response, err := http.Get("http://" + serverUrl)
	errorHandler(err)

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code:", response.StatusCode)
		panic(-1)
	}

	body, err := ioutil.ReadAll(response.Body)
	errorHandler(err)

	var peers []string

	err = json.Unmarshal(body, &peers)

	errorHandler(err)

	return peers

}
