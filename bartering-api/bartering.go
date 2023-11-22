package bartering

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

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

type NodeRatio struct {
	NodeIP string
	Ratio  float64
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

var PORT = "8084"

// var currentRatio float64

func InitiateBarter(peer string, ratios []NodeRatio) error {
	/*
		Function to barter the storage ratio
		Arguments : IP of peer as string
	*/

	currentRatio, err := FindNodeRatio(ratios, peer)

	if err != nil {
		return errors.New(err.Error())
	}

	newRatio := calculateNewRatio(currentRatio)

	barterMessage := "BarRq" + strconv.FormatFloat(newRatio, 'f', -1, 64)

	response := contactNodeForBarter(peer, barterMessage)

	if response == "OK\n" {
		// update that ratio value
		updatePeerRatio(ratios, peer, newRatio)
	} else {
		// in this case we have received a response to our barter message, we have to deal w it
		ratio, err := strconv.ParseFloat(response[:len(response)-1], 64)
		utils.ErrorHandler(err)
		updatePeerRatio(ratios, peer, ratio)
	}
	return nil
}

func RespondToBarterMsg(barterMsg string, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore, conn net.Conn, ratios []NodeRatio) {
	/*
		Function to answer a barter request
		Arguments : message received as a string, the peer who sent it as a string, the storage space on the node as float64,
		the bytes stored at each peer as a list of PeerStorageUse objects, the scores of peers as a list of NodeScore objects
	*/

	barterMsg_ratioRq := barterMsg[5:8]
	fmt.Println("Ratio received : ", barterMsg_ratioRq)
	barterMsg_ratio, err := strconv.ParseFloat(barterMsg_ratioRq, 64)
	fmt.Println("conversion OK")
	utils.ErrorHandler(err)

	if shouldRatioBeAccepted(barterMsg_ratio, peer, storageSpace, bytesAtPeers, scores) {
		fmt.Println("New ratio is accepted -- sending OK to other peer")
		_, err := io.WriteString(conn, "OK\n")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Sent OK message to peer")
		}
		updatePeerRatio(ratios, peer, barterMsg_ratio)
		fmt.Println(ratios)

	} else {
		// formulate new ratio proposition
		fmt.Println("New ratio should not be accepted -- sending another proposition to the peer")
		newRatio := formulateBarterResponse(peer, scores, storageSpace, bytesAtPeers)
		fmt.Println("New ratio :", newRatio)
		// toSend := "BarAn" + fmt.Sprintf("%f", newRatio)
		toSend := fmt.Sprintf("%f\n", newRatio)
		_, err = io.WriteString(conn, toSend)
		utils.ErrorHandler(err)
		updatePeerRatio(ratios, peer, newRatio)
	}

}

func updatePeerRatio(ratios []NodeRatio, peer string, newRatio float64) {
	for index, nodeRatio := range ratios {
		if nodeRatio.NodeIP == peer {
			ratios[index].Ratio = newRatio
		}
	}
}

func FindNodeRatio(ratios []NodeRatio, peer string) (float64, error) {
	/*
		Function to find a peer's current storage ratio
		Arguments : list of NodeRatio objects, peer ip as string
		Returns : storage ratio as float64 and nil if no error, 0 and error otherwise
	*/

	for _, nodeRatio := range ratios {
		if nodeRatio.NodeIP == peer {
			return nodeRatio.Ratio, nil
		}
	}

	return NodeRatio{}.Ratio, errors.New("peer not found when looking for ratio")

}

func formulateBarterResponse(peer string, scores []NodeScore, storageSpace float64, bytesAtPeers []PeerStorageUse) float64 {
	/*
		Function to counter barter in case the other node's proposition is not acceptable
		Arguments : peer id as string, list of NodeScore objects, node storage space as float64, list of PeerStorageUse objects
		Returns : new ratio as float64
	*/

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

func shouldRatioBeAccepted(ratio float64, peer string, storageSpace float64, bytesAtPeers []PeerStorageUse, scores []NodeScore) bool {
	/*
		Function to decided based off score and current storage space if the barter request can be accepted
		Arguments : current ratio as float64, the peer id as string, the current storage space as float64, the bytes stored at peers
		as a list of PeerStorageUse objects, the scores as a list of NodeScore objects
		Returns : boolean
	*/

	currentStorage, err := findPeerStorageUse(peer, bytesAtPeers)
	utils.ErrorHandler(err)
	if currentStorage.StorageAtNode == 0.0 {
		return true
	}

	return (isRatioTolerableGivenStorageSpace(peer, ratio, storageSpace, bytesAtPeers) && (ratio < calculateMaxAcceptableRatio(peer, scores, storageSpace, bytesAtPeers)))
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

	return PeerStorageUse{}, errors.New("peer not found")
}

func calculateMaxAcceptableRatio(peer string, scores []NodeScore, storageSpace float64, bytesAtPeers []PeerStorageUse) float64 {
	/*
		Function to calculate the maximum acceptable ratio given a peer's score
		Arguments : IP of peer as string, peer scores as NodeScore object list
		Return : max acceptabel ratio as float64
	*/

	peerScore, err := findPeerScore(peer, scores)
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

func findPeerScore(peer string, scores []NodeScore) (NodeScore, error) {
	/*
		Function to find a peer's score
		Arguments : IP of peer as string
	*/
	for _, peerInList := range scores {
		if peerInList.NodeIP == peer {
			return peerInList, nil
		}

	}
	return NodeScore{}, errors.New("peer not in peers list")
}

func contactNodeForBarter(peer string, msg string) string {
	/*
		Function to setup tcp connection to contact node to barter ratio
		Arguments : IP of peer as string, message to send as a string
		Returns : string of the peer's response
	*/
	conn, err := net.Dial("tcp", peer+":"+PORT)
	utils.ErrorHandler(err)

	defer conn.Close()

	_, err = io.WriteString(conn, msg)

	utils.ErrorHandler(err)
	fmt.Println("barter message sent")
	time.Sleep(2 * time.Second)
	response := bufio.NewReader(conn)
	fmt.Println("Response received!")
	responseString, err := response.ReadString('\n')
	fmt.Println(responseString)
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
