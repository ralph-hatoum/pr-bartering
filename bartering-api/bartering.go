package bartering

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"

	"../functions"
	"../utils"
)

type NodeScore struct {
	NodeIP string
	Score  float64
}

type PeerStorageUse struct {
	NodeIP        string
	StorageAtNode float64
}

var InitScore = 10.0

var AcceptanceTolerance = 0.5

var AcceptanceToleranceMutex sync.Mutex

var NodeTotalStorageSpace = 200

var RatioIncreaseRate = 0.1

var factorAcceptableRatio = 0.3

var PORT = "8083"

var currentRatio float64

func InitiateBarter(peer string) {

	newRatio := calculateNewRatio(currentRatio)

	barterMessage := "BartReq" + strconv.FormatFloat(newRatio, 'f', -1, 64)

	response := contactNodeForBarter(peer, barterMessage)

	if response == "OK" {
		// update that ratio value
		updateRatio(newRatio)
	} else {
		// in this case we have received a response to our barter message, we have to deal w it
		ratio, err := strconv.ParseFloat(response, 64)
		utils.ErrorHandler(err)
		shouldResponseRatioBeAccepted(ratio)
	}

}

func respondToBarterMsg(barterMsg string, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore) {

	barterMsg_ratioRq := barterMsg[7:]
	barterMsg_ratio, err := strconv.ParseFloat(barterMsg_ratioRq, 64)
	utils.ErrorHandler(err)

	shouldRatioBeAccepted(barterMsg_ratio, peer, storageSpace, bytesAtPeers, scores)

}

func calculateNewRatio(ratio float64) float64 {
	return ratio * (1 + RatioIncreaseRate)
}

func updateRatio(ratio float64) {

}

func shouldRatioBeAccepted(ratio float64, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore) bool {
	return isRatioTolerableGivenStorageSpace(peer, ratio, storageSpace, bytesAtPeers) && (ratio < calculateMaxAcceptableRatio(peer, scores))
}

func shouldResponseRatioBeAccepted(ratio float64) bool {
	return true
}

func isRatioTolerableGivenStorageSpace(peer string, ratio float64, storageSpace float64, bytesAtPeers []PeerStorageUse) bool {
	peerStorageUse, err := findPeerStorageUse(peer, bytesAtPeers)
	utils.ErrorHandler(err)

	return peerStorageUse.StorageAtNode < storageSpace
}

func findPeerStorageUse(peer string, bytesAtPeers []PeerStorageUse) (PeerStorageUse, error) {

	for _, peerStorageUse := range bytesAtPeers {
		if peerStorageUse.NodeIP == peer {
			return peerStorageUse, nil
		}
	}

	return PeerStorageUse{}, errors.New("Peer not found")
}

func calculateMaxAcceptableRatio(peer string, scores []NodeScore) float64 {
	peerScore, err := fincPeerScore(peer, scores)
	utils.ErrorHandler(err)

	return factorAcceptableRatio * peerScore.Score
}

func fincPeerScore(peer string, scores []NodeScore) (NodeScore, error) {
	for _, peerInList := range scores {
		if peerInList.NodeIP == peer {
			return peerInList, nil
		}

	}
	return NodeScore{}, errors.New("Peer not in peers list")
}

func contactNodeForBarter(peer string, msg string) string {
	conn, err := net.Dial("tcp", peer+PORT)
	utils.ErrorHandler(err)

	defer conn.Close()

	_, err = io.WriteString(conn, msg)

	response := bufio.NewReader(conn)

	responseString, err := response.ReadString('\n')

	utils.ErrorHandler(err)

	return responseString
}

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

func ElectStorageNodes(peerScores []NodeScore, numberOfNodes int) []string {
	// TODO change so we also eelct

	electedNodesScores := []NodeScore{}
	for _, peerScore := range peerScores {

		if len(electedNodesScores) < numberOfNodes {

			electedNodesScores = append(electedNodesScores, peerScore)

		} else {
			for index, currentlyElectedNode := range electedNodesScores {

				if currentlyElectedNode.Score < peerScore.Score {
					electedNodesScores[index] = peerScore
					break
				}
			}
		}

	}

	electedNodes := []string{}

	for _, electedNodeScore := range electedNodesScores {
		electedNodes = append(electedNodes, electedNodeScore.NodeIP)
	}

	return electedNodes
}

func CheckCIDValidity(storageRequest functions.StorageRequest) {
	/*
		Check if : CID is valid and exists
	*/
}

func CheckFileSizeValidity(storageRequest functions.StorageRequest) {
	/*
		Check if fileSize announced in storage request is declared honestly
	*/
}

func CheckEnoughSpace(storageRequest functions.StorageRequest, currentStorageSpace float64) bool {
	/*
		Check if self has enough space to store the file
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
		In our case for now we will consider that if the storage is refused,
		then the tolerance needs to go up
	*/

	fileSize := storageRequest.fileSize

	delta := fileSize / NodeTotalStorageSpace

	increaseTolerance(delta)

}

func craftNewRq(storageRequest functions.StorageRequest) functions.StorageRequest {
	/*
		Function to craft a new better suited request aftet it was refused
	*/

}

func increaseTolerance(delta float64) {
	/*
		Function to increase tolerance
	*/

	AcceptanceToleranceMutex.Lock()
	AcceptanceTolerance += delta
	AcceptanceToleranceMutex.Unlock()

}

func decreaseTolerance(delta float64) {
	/*
		Function to decrease tolerance
	*/

	AcceptanceToleranceMutex.Lock()
	AcceptanceTolerance -= delta
	AcceptanceToleranceMutex.Unlock()

}
