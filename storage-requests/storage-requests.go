package storagerequests

import (
	api_ipfs "bartering/api-ipfs"
	datastructures "bartering/data-structures"
	"bartering/utils"
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

func BuildStorageRequestMessage(storageRequest datastructures.StorageRequest) string {

	/*
		Function to build a storage request message from a StorageRequest object
		Arguments : StorageRequest object
		Output : string
	*/

	fileSizeString := fmt.Sprintf("%.*f", 10, storageRequest.FileSize)
	storageRequestMessage := "StoRq" + storageRequest.CID + fileSizeString

	return storageRequestMessage
}

func buildFulfilledRequestObject(CID string, peer string) datastructures.FulfilledRequest {

	/*
		Function to build a fulfilled request object
		Arguments : CID as string, peer as string
		Output : FulfilledRequest object
	*/

	fulfilledRequest := datastructures.FulfilledRequest{CID: CID, Peer: peer}

	return fulfilledRequest
}

func addFulFilledRequestToFulfilledRequests(request datastructures.FulfilledRequest, requests *[]datastructures.FulfilledRequest) {

	/*
		Function to add fulfilled request object to fulfilled requests list
		Arguments : fulfilledRequest, fulfilledRequests list
	*/

	*requests = append(*requests, request)
}

func updateFulfilledRequests(CID string, peer string, fulfilledRequests *[]datastructures.FulfilledRequest) {

	/*
		Function to add a fulfilled request from CID and peer name to fulfilled requests
		Arguments : CID as string, peer id as string, fulfilled requests list pointer
	*/

	newRequest := buildFulfilledRequestObject(CID, peer)

	addFulFilledRequestToFulfilledRequests(newRequest, fulfilledRequests)

}

func RequestStorageFromPeer(peer string, storageRequest datastructures.StorageRequest, port string, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, fulfilledRequests *[]datastructures.FulfilledRequest, scoreDecreaseRefStoReq float64) {

	/*
		Function to request storage from a peer
		Arguments : peer id as string, storageRequest object, port to contact peer on as string, PeerStorageUse array, NodeScore array, fulfilledRequests array pointer
	*/

	storageRqMessage := BuildStorageRequestMessage(storageRequest)

	conn, err := net.Dial("tcp", peer+":"+port)

	utils.ErrorHandler(err)

	_, err = io.WriteString(conn, storageRqMessage)

	utils.ErrorHandler(err)

	response := bufio.NewReader(conn)

	responseString, err := response.ReadString('\n')

	utils.ErrorHandler(err)
	fmt.Println(responseString)

	if responseString == "OK\n" {
		fmt.Println("Peer ", peer, " stored file with CID ", storageRequest.CID, " successfully.")
		updateBytesAtPeers(bytesAtPeers, peer, storageRequest)
		updateFulfilledRequests(storageRequest.CID, peer, fulfilledRequests)
	} else if responseString == "KO\n" {
		fmt.Println("Storage refused by node, decreasing score")
		updatePeerScoreRefusedRq(scores, peer, scoreDecreaseRefStoReq)
	}
}

func updatePeerScoreRefusedRq(scores []datastructures.NodeScore, peer string, scoreDecreaseRefStoReq float64) {

	/*
		Function used to update the peer score (decrease it) upon refusing a storage request
		Arguments : NodeScore array, peer id as string
	*/

	for index, peerScore := range scores {
		if peerScore.NodeIP == peer {
			scores[index].Score -= scoreDecreaseRefStoReq
		}
	}
}

func updateBytesAtPeers(bytesAtPeers []datastructures.PeerStorageUse, peer string, storageRequest datastructures.StorageRequest) {

	/*
		Function to update a PeerStorageUse object in a PeerStorageUse
		Arguments : PeerStorageUse array, peer id as string, storageRequest object
	*/

	for index, bytesAtPeer := range bytesAtPeers {
		if bytesAtPeer.NodeIP == peer {
			bytesAtPeers[index].StorageAtNode += storageRequest.FileSize
		}
	}
}

func updateBytesForPeers(bytesForPeers []datastructures.PeerStorageUse, peer string, fileSize float64) {

	/*
		Function to update a PeerStorageUse object in a PeerStorageUse
		Arguments : PeerStorageUse array, peer id as string, storageRequest object
	*/

	for index, bytesForPeer := range bytesForPeers {
		if bytesForPeer.NodeIP == peer {
			bytesForPeers[index].StorageAtNode += fileSize
		}
	}
}

func HandleStorageRequest(bufferString string, conn net.Conn, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest) {

	/*
		Function to handle a storage message type message
		Arguments : buffer received through a tcp connection, as a string, net.Conn object, PeerStorageUse array, pointer to fulfilledRequest array
	*/

	peer := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	var messageToPeer string
	fmt.Println("Received storage request")
	CID := bufferString[5:51]
	fmt.Println("CID : ", CID)
	fileSize := bufferString[51:]
	fileSize = strings.Split(fileSize, "\n")[0]
	fmt.Println("File Size : ", fileSize)

	fileSizeFloat, err := strconv.ParseFloat(fileSize, 64)
	utils.ErrorHandler(err)

	request := datastructures.StorageRequest{FileSize: fileSizeFloat, CID: CID}

	fmt.Println("Storage request : ", request, " ; checking validity ...")

	if CheckRqValidity(request) {
		fmt.Println("Request ", request, " valid, storing ! ")
		fmt.Println("Pinning to IPFS ...")
		api_ipfs.PinToIPFS(CID)
		fmt.Println("File pinned to IPFS!")
		messageToPeer = "OK\n"
		updateBytesForPeers(bytesForPeers, peer, fileSizeFloat)
		updateFulfilledRequests(CID, peer, storedForPeers)
		fmt.Println("stored for peers : ", storedForPeers)
	} else {
		fmt.Println("Request ", request, " not valid, not storing ! ")

		messageToPeer = "KO\n"
	}

	_, err = io.WriteString(conn, messageToPeer)

	utils.ErrorHandler(err)

}

func HandleStorageRequestTimed(bufferString string, conn net.Conn, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, deletionQueue *[]datastructures.StorageRequestTimedAccepted) {

	/*
		SHOULD REPLACE HANDLESTORAGEREQUEST FUNC ONCE EVERYTHING IS DONE, TESTED AND WORKING
		Function to handle a storage message type message
		Arguments : buffer received through a tcp connection, as a string, net.Conn object, PeerStorageUse array, pointer to fulfilledRequest array
	*/

	peer := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	var messageToPeer string
	fmt.Println("Received storage request")
	CID := bufferString[5:51]
	fmt.Println("CID : ", CID)
	DurationMinutes := bufferString[51:]
	DurationMinutes = strings.Split(DurationMinutes, "\n")[0]
	fmt.Println("Request duration (minutes) : ", DurationMinutes)

	DurationMinutesInt, err := strconv.ParseInt(DurationMinutes, 10, 64)
	utils.ErrorHandler(err)

	request := datastructures.StorageRequestTimed{DurationMinutes: DurationMinutesInt, CID: CID}

	fmt.Println("Storage request : ", request, " ; checking validity ...")

	if CheckRqValidityTimed(request) {
		fmt.Println("Request ", request, " valid, storing ! ")
		fmt.Println("Pinning to IPFS ...")
		api_ipfs.PinToIPFS(CID)
		fmt.Println("File pinned to IPFS!")
		messageToPeer = "OK\n"
		// updateBytesForPeers(bytesForPeers, peer, fileSizeFloat)
		updateFulfilledRequests(CID, peer, storedForPeers)
		fmt.Println("stored for peers : ", storedForPeers)
		requestAccepted := buildStorageRequestTimedAcceptedObjectFromStorageRequestTimed(request)
		fmt.Println("accepted  timed request : ", requestAccepted)
		appendStorageRequestToDeletionQueue(requestAccepted, deletionQueue)
		fmt.Println("deletion queue", deletionQueue)
	} else {
		fmt.Println("Request ", request, " not valid, not storing ! ")

		messageToPeer = "KO\n"
	}

	_, err = io.WriteString(conn, messageToPeer)

	utils.ErrorHandler(err)

}

func CheckRqValidityTimed(storageRequest datastructures.StorageRequestTimed) bool {
	deadline := computeDeadlineFromTimedStorageRequest(storageRequest)

	return time.Now().Before(deadline)
}

func ElectStorageNodes(peerScores []datastructures.NodeScore, numberOfNodes int) ([]string, error) {

	/*
		Function to elect nodes to whom self will send storage requests
		Arguments : IP of peer as string, number of nodes as int
		Returns : list of strings containing IPs of nodes to contact
	*/

	// TODO change so we also elect low score nodes to give the opportunity to raise the score

	if numberOfNodes > len(peerScores) {
		return []string{}, errors.New("asking for more peers than we know")
	}

	electedNodesScores := []datastructures.NodeScore{}
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

	return electedNodes, nil
}

func CheckRqValidity(storageRequest datastructures.StorageRequest) bool {

	/*
		Function to decide if a received storageRequest should be accepted or not
	*/

	return true
}

func CheckCIDValidity(storageRequest datastructures.StorageRequest) bool {

	/*
		Check if : CID is valid and exists
		problem : ipfs cat with wrong CID goes through ipfs search which can be very long
		no efficient way to check if CID exists
		however if we work under the hypothesis that peers in our network have preestablished ipfs links
	*/

	return true
}

func CheckFileSizeValidity(storageRequest datastructures.StorageRequest) bool {

	/*
		Check if fileSize announced in storage request is declared honestly
	*/

	return true
}

func CheckEnoughSpace(storageRequest datastructures.StorageRequest, currentStorageSpace float64, NodeTotalStorageSpace float64) bool {

	/*
		Check if self has enough space to store the file
		Arguments : storage request of type StorageRequest, current storage space used as float64
		Returns : boolean
	*/

	return storageRequest.FileSize+float64(currentStorageSpace) < NodeTotalStorageSpace

}

func GarbageCollector(storageDeletionQueue []datastructures.StorageRequestTimedAccepted) {

	/*
		Function to run in background to perform garage collection, aka deal with requests that have expired
		Arguments : storageDeletionQueue (slice of StorageRequestTimedAccepted objects)
	*/

	for {
		if len(storageDeletionQueue) != 0 {
			if storageDeletionQueue[0].Deadline.Before(time.Now()) {
				storageDeletionQueue = garbageCollectionStrategy(storageDeletionQueue)
			}
		}
	}
}

func garbageCollectionStrategy(storageDeletionQueue []datastructures.StorageRequestTimedAccepted) []datastructures.StorageRequestTimedAccepted {

	/*
		Garbage collection strategy
		Essentially, we might not want our node to directly deleted expired requests
		(for example, only delete when no more storage is available to increase availabilty of data)
		Strategy should be defined in this function
		Arguments : storageDeletionQueue (slice of StorageRequestTimedAccepted objects)
	*/
	storageDeletionQueue = storageDeletionQueue[1:]
	return storageDeletionQueue
}

func appendStorageRequestToDeletionQueue(storageRequest datastructures.StorageRequestTimedAccepted, deletionQueue *[]datastructures.StorageRequestTimedAccepted) {

	sortedSlice := *deletionQueue

	if len(sortedSlice) == 0 {
		sortedSlice = append(sortedSlice, storageRequest)
		*deletionQueue = sortedSlice
		return
	}
	i := len(sortedSlice) - 1
	// fmt.Println("lenght of sorted slice : ", len(sortedSlice))
	// fmt.Println("value of i : ", i)
	// fmt.Println("printed element of sortedSlice : ", sortedSlice[i].Deadline)
	for i >= 0 && storageRequest.Deadline.After(sortedSlice[i].Deadline) {
		// sortedSlice[i+1] = sortedSlice[i]
		i--
	}

	fmt.Println("index to add at : ", i)

	sortedSlice[i+1] = storageRequest
	*deletionQueue = sortedSlice

	// fmt.Println("queue before appending : ", queue)

	// if len(queue) == 0 {
	// 	new := append(queue, storageRequest)
	// 	*deletionQueue = new
	// 	return
	// }

	// index := 0
	// for storageRequest.Deadline.After(queue[index].Deadline) {
	// 	index += 1
	// }
	// before := append([]datastructures.StorageRequestTimedAccepted{}, queue[:index]...)
	// after := append([]datastructures.StorageRequestTimedAccepted{}, queue[index:]...)
	// newQueue := append(before, storageRequest)
	// newQueue = append(newQueue, after...)
	// *deletionQueue = newQueue
}

func computeDeadlineFromTimedStorageRequest(storageRequest datastructures.StorageRequestTimed) time.Time {

	timeToAdd := time.Duration(storageRequest.DurationMinutes) * time.Minute

	deadline := time.Now().Add(timeToAdd)

	return deadline
}

func buildStorageRequestTimedAcceptedObjectFromStorageRequestTimed(storageRequest datastructures.StorageRequestTimed) datastructures.StorageRequestTimedAccepted {

	CID := storageRequest.CID
	deadline := computeDeadlineFromTimedStorageRequest(storageRequest)

	return datastructures.StorageRequestTimedAccepted{CID: CID, Deadline: deadline}
}
