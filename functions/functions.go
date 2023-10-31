package functions

import (
	"fmt"

	api_ipfs "bartering/api-ipfs"

	bootstrapconnect "bartering/bootstrap-connect"

	"bartering/utils"

	"bartering/bartering-api"
)

type StorageRequest struct {
	/*
		Data structure to represent storage requests ; consist of the CID of a file and its size
	*/

	CID      string
	fileSize float64
}

func NodeStartup() ([]string, []StorageRequest, []StorageRequest, []string, []bartering.PeerStorageUse, []bartering.NodeScore, []bartering.NodeRatio) {
	/*
		UNFINISHED
		Function called upon a node's startup
		This function will create all needed lists :
			- storagePool : list of CIDs of node's data
			- pendingRequests : list of storage requests that are awaiting to be given to a peer in the network
			- fulfilledRequests : requests accepted by other peers in the network
			- peers : list of peers ids
			- bytesAtPeers : list of bytes stored by other peers
			- scores : list of scores attributed to each peer
			- ratios : storage ratios for each peer

		This function will call the bootstrap to retrieve peers' IP addresses, and store them in the peers list

		Arguments : None
		Returns : storagePool as list of strings, pendingRequests and fulfilledRequests as StorageRequests lists, peers as list of strings
	*/

	fmt.Println("Starting node")

	fmt.Println("Creating storage pool and requests lists")
	storage_pool, pending_requests, fulfilled_requests := createStorageRequestsLists()

	fmt.Println("Creating peers list")
	peers := bootstrapconnect.GetPeersFromBootstrapHTTP("127.0.0.1", "8080")

	fmt.Println("Creating bytes at peers list, scores and ratios")
	bytesAtPeers := []bartering.PeerStorageUse{}
	scores := []bartering.NodeScore{}
	ratios := []bartering.NodeRatio{}

	return storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, scores, ratios
}

func Store(path string, storage_pool []string, pending_requests []StorageRequest) {
	/*
		UNFINISHED
		Function called when a new file needs to be stored on the network
		This function will :
			- add the file to IPFS, pin it and get its CID
			- retrieve the file's size and build a StorageRequest data object with the CID and the file's size
			- add the storage requests to the pendingRequests list
	*/

	CID := api_ipfs.UploadToIPFS(path)

	storage_pool = append(storage_pool, CID)

	fmt.Println(storage_pool)

	file_size := utils.GetFileSize(path)

	fmt.Println(file_size)

	storage_request := StorageRequest{CID, file_size}

	pending_requests = append(pending_requests, storage_request)

	fmt.Println("Pending requests : ", pending_requests)
}

func createStorageRequestsLists() ([]string, []StorageRequest, []StorageRequest) {
	/*
		Function to create all needed data structures
		Argument : None
		Returns : storage_pool as string list, pending and fulfilled requests lists as StorageRequest lists
	*/

	storage_pool := []string{}

	pending_requests := []StorageRequest{}

	fulfilled_requests := []StorageRequest{}

	return storage_pool, pending_requests, fulfilled_requests

}

func propagateToPeers(storageRequest StorageRequest) {
	messageToPropagate := buildStorageRequestMessage(storageRequest)
	fmt.Println(messageToPropagate)

	// Choose peers to propagate to
	// send request, await accept ?
	// If refuse or no answer, make better offer ?
}

func buildStorageRequestMessage(storageRequest StorageRequest) string {

	fileSizeString := fmt.Sprintf("%.*f", 10, storageRequest.fileSize)
	storageRequestMessage := "StoReq" + storageRequest.CID + fileSizeString

	return storageRequestMessage
}
