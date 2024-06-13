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

func StoreKCopiesOnNetwork(peerScores []datastructures.NodeScore, K int, storageRequest datastructures.StorageRequest, port string, bytesAtPeers []datastructures.PeerStorageUse, fulfilledRequests *[]datastructures.FulfilledRequest, scoreDecreaseRefStoReq float64) int {

	okRqs := 0
    usedPeers := make(map[string]bool)
    tries := 0

    for tries < 3 {
        peersToRequest, err := ElectStorageNodes(peerScores, K, usedPeers)
        if err != nil {
            fmt.Println(err)
            return 0
        }

        for _, peer := range peersToRequest {
            if usedPeers[peer] {
                continue
            }
            response := RequestStorageFromPeer(peer, storageRequest, port, bytesAtPeers, peerScores, fulfilledRequests, scoreDecreaseRefStoReq)
            if response == "OK\n" {
                okRqs++
				fmt.Printf("The value of counter is %d.\n ",okRqs)
                usedPeers[peer] = true
                if okRqs == K {
                    fmt.Println("Reached required number of copies")
                    return okRqs
                }
            }
        }
        fmt.Println("Could not reach number of copies ... choosing new nodes")
        tries++
    }

    fmt.Println("Could not reach desired number of copies - only got ", okRqs)
    return okRqs

}

func RemovePeerFromPeers(peerScores []datastructures.NodeScore, peerToRm string) []datastructures.NodeScore {
	for index, peer := range peerScores {
		if peer.NodeIP == peerToRm {
			peerScores = append(peerScores[:index], peerScores[index+1:]...)
		}
	}
	return peerScores
}

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
	//  TODO FIX
	*requests = append(*requests, request)
}

func updateFulfilledRequests(CID string, peer string, fulfilledRequests *[]datastructures.FulfilledRequest) {

	/*
		Function to add a fulfilled request from CID and peer name to fulfilled requests
		Arguments : CID as string, peer id as string, fulfilled requests list pointer
	*/

	fmt.Println("Updating fulfilled requests")

	fmt.Println("before : ", fulfilledRequests)
	newRequest := buildFulfilledRequestObject(CID, peer)

	addFulFilledRequestToFulfilledRequests(newRequest, fulfilledRequests)
	fmt.Println("after : ", fulfilledRequests)

}

func RequestStorageFromPeer(peer string, storageRequest datastructures.StorageRequest, port string, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, fulfilledRequests *[]datastructures.FulfilledRequest, scoreDecreaseRefStoReq float64) string {
	fmt.Println("Requesting storage from peer:", peer)

    storageRqMessage := BuildStorageRequestMessage(storageRequest)

    // Attempt to establish a connection with the peer
    conn, err := net.Dial("tcp", peer+":"+port)
    if err != nil {
        fmt.Printf("Failed to connect to peer %s: %v\n", peer, err)
        return "Failed to connect"
    }
    defer conn.Close() // Ensure the connection is closed after function execution

    // Send the storage request message to the peer
    _, err = io.WriteString(conn, storageRqMessage)
    if err != nil {
        fmt.Printf("Failed to send storage request to peer %s: %v\n", peer, err)
        return "Failed to send"
    }

    // Read the response from the peer
    response, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        fmt.Printf("Failed to read response from peer %s: %v\n", peer, err)
        return "Failed to read"
    }

    fmt.Println("Response from peer:", response)

    // Process the response
    if response == "OK\n" {
        fmt.Printf("Peer %s stored file with CID %s successfully.\n", peer, storageRequest.CID)
        updateBytesAtPeers(bytesAtPeers, peer, storageRequest)
        updateFulfilledRequests(storageRequest.CID, peer, fulfilledRequests)
        return "OK\n"
    } else if response == "KO\n" {
        fmt.Printf("Storage request refused by peer %s, decreasing score.\n", peer)
        updatePeerScoreRefusedRq(scores, peer, scoreDecreaseRefStoReq)
        return "KO\n"
    }

    return "Unknown response"
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
		AppendStorageRequestToDeletionQueue(requestAccepted, deletionQueue)
		fmt.Println("deletion queue", deletionQueue)
	} else {
		fmt.Println("Request ", request, " not valid, not storing ! ")

		messageToPeer = "KO\n"
	}

	_, err = io.WriteString(conn, messageToPeer)

	utils.ErrorHandler(err)

}

func CheckRqValidityTimed(storageRequest datastructures.StorageRequestTimed) bool {
	deadline := ComputeDeadlineFromTimedStorageRequest(storageRequest)

	return time.Now().Before(deadline)
}

func ElectStorageNodes(peerScores []datastructures.NodeScore, numberOfNodes int, usedPeers map[string]bool) ([]string, error) {
    if numberOfNodes > len(peerScores) {
        return nil, errors.New("asking for more peers than we have available")
    }

    var electedNodes []string
    for _, peerScore := range peerScores {
        if len(electedNodes) >= numberOfNodes {
            break
        }
        if !usedPeers[peerScore.NodeIP] {
            electedNodes = append(electedNodes, peerScore.NodeIP)
        }
    }

    if len(electedNodes) < numberOfNodes {
        return nil, fmt.Errorf("not enough available peers to satisfy the request")
    }
    return electedNodes, nil
}

func ElectStorageNodesLowAndHigh(peerScores []datastructures.NodeScore, numberOfNodes int) []datastructures.NodeScore {

	/*
		Function to elect nodes to whom self will send storage requests - low and high scores
		nodeScore list NEEDS to be SORTED for this function to behave correctly
		Arguments : nodeScore list, number of nodes as int
		Returns : list of strings containing IPs of nodes to contact
	*/

	lowScoreProportion := 0.2
	nbNodesChosen := 0
	chosen := []datastructures.NodeScore{}
	for float64(nbNodesChosen) < lowScoreProportion*float64(numberOfNodes) {
		chosen = append(chosen, peerScores[nbNodesChosen])
		nbNodesChosen += 1
		peerScores = peerScores[1:]
	}
	// index := 1
	for nbNodesChosen < numberOfNodes {
		chosen = append(chosen, peerScores[len(peerScores)-1])
		nbNodesChosen += 1
		// fmt.Println(len(peerScores))
		peerScores = peerScores[:len(peerScores)-1]
	}
	return chosen

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
				storageDeletionQueue = GarbageCollectionStrategy(storageDeletionQueue)
			}
		}
	}

}

func GarbageCollectionStrategy(storageDeletionQueue []datastructures.StorageRequestTimedAccepted) []datastructures.StorageRequestTimedAccepted {

	/*
		Garbage collection strategy
		Essentially, we might not want our node to directly deleted expired requests
		(for example, only delete when no more storage is available to increase availabilty of data)
		Strategy should be defined in this function
		Arguments : storageDeletionQueue (slice of StorageRequestTimedAccepted objects)
	*/
	if len(storageDeletionQueue) != 0 {
		storageDeletionQueue = storageDeletionQueue[1:]
	}
	return storageDeletionQueue
}

func AppendStorageRequestToDeletionQueue(storageRequest datastructures.StorageRequestTimedAccepted, deletionQueue *[]datastructures.StorageRequestTimedAccepted) {

	queue := *deletionQueue
	newQueue := AuxInsertInSortedList(storageRequest, queue)
	*deletionQueue = newQueue

}

func AuxInsertInSortedList(storageRequest datastructures.StorageRequestTimedAccepted, queue []datastructures.StorageRequestTimedAccepted) []datastructures.StorageRequestTimedAccepted {

	if len(queue) == 0 {

		queue = append(queue, storageRequest)
		return queue

	} else {

		if storageRequest.Deadline.After(queue[0].Deadline) {

			newQueue := AuxInsertInSortedList(storageRequest, queue[1:])
			head := []datastructures.StorageRequestTimedAccepted{queue[0]}
			newQueue = append(head, newQueue...)
			return newQueue

		} else {

			newQueue := []datastructures.StorageRequestTimedAccepted{storageRequest}
			newQueue = append(newQueue, queue...)
			return newQueue

		}
	}
}

func ComputeDeadlineFromTimedStorageRequest(storageRequest datastructures.StorageRequestTimed) time.Time {

	timeToAdd := time.Duration(storageRequest.DurationMinutes) * time.Minute

	deadline := time.Now().Add(timeToAdd)

	return deadline
}

func buildStorageRequestTimedAcceptedObjectFromStorageRequestTimed(storageRequest datastructures.StorageRequestTimed) datastructures.StorageRequestTimedAccepted {

	CID := storageRequest.CID
	deadline := ComputeDeadlineFromTimedStorageRequest(storageRequest)

	return datastructures.StorageRequestTimedAccepted{CID: CID, Deadline: deadline}
}
