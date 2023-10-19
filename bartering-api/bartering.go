package bartering

import (
	"../functions"
	"./functions"
)

type NodeScore struct {
	NodeIP string
	Score  float64
}

var InitScore = 10.0

var AcceptanceTolerance = 10.0

var NodeTotalStorageSpace = 200

func InitNodeScores(peers []string) []NodeScore {
	/*
	Function to initiate node scores 
	*/

	scores := []NodeScore{}

	for _, peer := range peers {
		score := NodeScore{NodeIP: peer, Score: InitScore}
		scores = append(scores, score)
	}

	return scores
}

func ElectStorageNodes() []string {
	return []string{}
}

func CheckCIDValidity(storageRequest functions.StorageRequest) {
	/* Check if : CID is valid and exists */
}

func CheckFileSizeValidity(storageRequest functions.StorageRequest) {
	/* Check if fileSize announced in storage request is declared honestly */
}

func CheckEnoughSpace(storageRequest functions.StorageRequest, currentStorageSpace float64) bool {
	/* Check if self has enough space to store the file
	Arguments : storage request of type StorageRequest, current storage space used as float64
	Returns : boolean
	*/
	if storageRequest.fileSize+currentStorageSpace < NodeTotalStorageSpace {
		return true
	}
	return false
}

func dealWithRefusedRequest(storageRequest functions.StorageRequest) {
	/*
		Function to deal with a refused storage request
	*/
}

func craftNewRq(storageRequest functions.StorageRequest) functions.StorageRequest {
	/*
	Function to craft a new better suited request aftet it was refused
	*/

}

func updateTolerance() {
	/*
	Function to update tolerance 
	*/

}

func shouldReqBeAccepted() bool {

}
