package bartering

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	// "bartering/functions"
	"bartering/utils"
)

type NodeScore struct {
	NodeIP string
	Score  float64
}

type PeerStorageUse struct {
	NodeIP        string
	StorageAtNode float64
}

type StorageRequest struct {
	/*
		Data structure to represent storage requests ; consist of the CID of a file and its size
	*/

	CID      string
	fileSize float64
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
	/*
		Function to barter the storage ratio
		Arguments : IP of peer as string
	*/

	newRatio := calculateNewRatio(currentRatio)

	barterMessage := "BarReq" + strconv.FormatFloat(newRatio, 'f', -1, 64)

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

func RespondToBarterMsg(barterMsg string, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore) {
	/*
		Function to answer a barter request
		Arguments : message received as a string, the peer who sent it as a string, the storage space on the node as float64,
		the bytes stored at each peer as a list of PeerStorageUse objects, the scores of peers as a list of NodeScore objects
	*/

	barterMsg_ratioRq := barterMsg[5:8]
	fmt.Println("Ratio received : ", barterMsg_ratioRq)
	barterMsg_ratio, err := strconv.ParseFloat(barterMsg_ratioRq, 64)
	utils.ErrorHandler(err)

	// shouldRatioBeAccepted(barterMsg_ratio, peer, storageSpace, bytesAtPeers, scores)

	if shouldRatioBeAccepted(barterMsg_ratio, peer, storageSpace, bytesAtPeers, scores) {
		// send "OK" to node
		fmt.Println("New ratio is accepted -- sending OK to other peer")
	} else {
		// formulate new ratio proposition
		fmt.Println("New ratio should not be accepted -- sending another proposition to the peer")
		newRatio := formulateBarterResponse(peer, scores, storageSpace, bytesAtPeers)
		fmt.Println("New ratio :", newRatio)
	}

}

func formulateBarterResponse(peer string, scores []NodeScore, storageSpace float64, bytesAtPeers []PeerStorageUse) float64 {

	maxPossible := calculateMaxAcceptableRatio(peer, scores, storageSpace, bytesAtPeers)
	return maxPossible
}

func calculateNewRatio(ratio float64) float64 {
	/*
		Function to  calculate the new ratio to use
		Arguments : the current ratio as float 64
		Returns : new ratio as float64
	*/

	return ratio * (1 + RatioIncreaseRate)
}

func updateRatio(ratio float64) {

}

func shouldRatioBeAccepted(ratio float64, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore) bool {
	/*
		Function to decided based off score and current storage space if the barter request can be accepted
		Arguments : current ratio as float64, the peer id as string, the current storage space as float64, the bytes stored at peers
		as a list of PeerStorageUse objects, the scores as a list of NodeScore objects
		Returns : boolean
	*/

	return isRatioTolerableGivenStorageSpace(peer, ratio, storageSpace, bytesAtPeers) && (ratio < calculateMaxAcceptableRatio(peer, scores, storageSpace, bytesAtPeers))
}

func shouldResponseRatioBeAccepted(ratio float64) bool {
	/*
		Function to decide whether the counter barter made by the other peer should be accepted or not
		Arguments : the proposed ratio as float64
		Returns : boolean
	*/

	return true
}

func isRatioTolerableGivenStorageSpace(peer string, ratio float64, storageSpace float64, bytesAtPeers []PeerStorageUse) bool {
	/*
		Function to decided if the propoed ratio is tolerable given the current storage space on the node
		Arguments : id of the peer as string, current ratio as float64, storage space on the node as float64, bytes stored at the peer
		as a PeerStorageUse object list
		Returns : boolean
	*/

	peerStorageUse, err := findPeerStorageUse(peer, bytesAtPeers)
	utils.ErrorHandler(err)

	return peerStorageUse.StorageAtNode*ratio < storageSpace
}

func findPeerStorageUse(peer string, bytesAtPeers []PeerStorageUse) (PeerStorageUse, error) {
	/*
		Function to find the storage used by self at a peer
		Arguments : peer id as string, bytes stored at all peers as a PeerStorageUse objects list
		Returns : PeerStorageUse object and nil if peer is found, empty PeerStorageUse object and error if peer not found
	*/
	for _, peerStorageUse := range bytesAtPeers {
		if peerStorageUse.NodeIP == peer {
			return peerStorageUse, nil
		}
	}

	return PeerStorageUse{}, errors.New("Peer not found")
}

func calculateMaxAcceptableRatio(peer string, scores []NodeScore, storageSpace float64, bytesAtPeers []PeerStorageUse) float64 {
	/*
		Function to calculate the maximum acceptable ratio given a peer's score
		Arguments : IP of peer as string, peer scores as NodeScore object list
		Return : max acceptabel ratio as float64
	*/

	peerScore, err := fincPeerScore(peer, scores)
	utils.ErrorHandler(err)

	ratio := factorAcceptableRatio * peerScore.Score

	if !isRatioTolerableGivenStorageSpace(peer, ratio, storageSpace, bytesAtPeers) {
		storageUsed, err := findPeerStorageUse(peer, bytesAtPeers)
		utils.ErrorHandler(err)
		ratio = storageSpace / storageUsed.StorageAtNode
		// maybe : ratio := ratio *peerScore.Score ?
	}

	return ratio
}

func fincPeerScore(peer string, scores []NodeScore) (NodeScore, error) {
	/*
		Function to find a peer's score
		Arguments : IP of peer as string
	*/
	for _, peerInList := range scores {
		if peerInList.NodeIP == peer {
			return peerInList, nil
		}

	}
	return NodeScore{}, errors.New("Peer not in peers list")
}

func contactNodeForBarter(peer string, msg string) string {
	/*
		Function to setup tcp connection to contact node to barter ratio
		Arguments : IP of peer as string, message to send as a string
		Returns : string of the peer's response
	*/
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
	/*
		Function to elect nodes to whom self will send storage requests
		Arguments : IP of peer as string, number of nodes as int
		Returns : list of strings containing IPs of nodes to contact
	*/

	// TODO change so we also elect low score nodes to give the opportunity to raise the score

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

func CheckCIDValidity(storageRequest StorageRequest) {
	/*
		Check if : CID is valid and exists
	*/
}

func CheckFileSizeValidity(storageRequest StorageRequest) {
	/*
		Check if fileSize announced in storage request is declared honestly
	*/
}

func CheckEnoughSpace(storageRequest StorageRequest, currentStorageSpace float64) bool {
	/*
		Check if self has enough space to store the file
		Arguments : storage request of type StorageRequest, current storage space used as float64
		Returns : boolean
	*/

	if storageRequest.fileSize+float64(currentStorageSpace) < float64(NodeTotalStorageSpace) {
		return true
	}
	return false
}

func dealWithRefusedRequest(storageRequest StorageRequest) {
	/*
		Function to deal with a refused storage request
		In our case for now we will consider that if the storage is refused,
		then the tolerance needs to go up
	*/

	fileSize := storageRequest.fileSize

	delta := fileSize / float64(NodeTotalStorageSpace)

	increaseTolerance(delta)

}

func craftNewRq(storageRequest StorageRequest) StorageRequest {
	/*
		Function to craft a new better suited request aftet it was refused
	*/
	return StorageRequest{}

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
