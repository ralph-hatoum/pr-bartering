package dumper

import (
	datastructures "bartering/data-structures"
	"encoding/json"
	"fmt"
	"net/http"
)

func Dumper(bytesAtPeers []datastructures.PeerStorageUse, bytesForPeers []datastructures.PeerStorageUse, fulfilledRequests []datastructures.FulfilledRequest, storagePool []string, pendingRequests []datastructures.StorageRequest, peers []string, scores []datastructures.NodeScore, ratiosForPeers []datastructures.NodeRatio, ratiosAtPeers []datastructures.NodeRatio, storedForPeers []datastructures.FulfilledRequest) {
	address := "0.0.0.0"
	port := "8083"

	serverAddr := address + ":" + port

	http.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {
		response, err := responseBuilder(bytesAtPeers, bytesForPeers, fulfilledRequests, storagePool, pendingRequests, peers, scores, ratiosForPeers, ratiosAtPeers, storedForPeers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	})

	err := http.ListenAndServe(serverAddr, nil)

	if err != nil {
		fmt.Println("Could not start dumper")
	}
}

func responseBuilder(bytesAtPeers []datastructures.PeerStorageUse, bytesForPeers []datastructures.PeerStorageUse, fulfilledRequests []datastructures.FulfilledRequest, storagePool []string, pendingRequests []datastructures.StorageRequest, peers []string, scores []datastructures.NodeScore, ratiosForPeers []datastructures.NodeRatio, ratiosAtPeers []datastructures.NodeRatio, storedForPeers []datastructures.FulfilledRequest) ([]byte, error) {
	output := make(map[string]any)

	output["bytesAtPeers"] = bytesAtPeers
	output["bytesForPeers"] = bytesForPeers
	output["fulfilledRequests"] = fulfilledRequests
	output["storagePool"] = storagePool
	output["pendingRequests"] = pendingRequests
	output["peers"] = peers
	output["scores"] = scores
	output["ratiosAtPeers"] = ratiosAtPeers
	output["ratiosForPeers"] = ratiosForPeers

	jsonResponse, err := json.Marshal(output)

	if err != nil {
		fmt.Println("Could not build /dump response")
		return nil, fmt.Errorf("could not build /dump resp")
	}

	return jsonResponse, nil
}
