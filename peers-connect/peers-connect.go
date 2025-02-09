package peersconnect

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"bartering/bartering-api"
	datastructures "bartering/data-structures"
	"bartering/utils"
)

func ListenPeersRequestsTCPFailure(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, failureMutex *sync.Mutex, msgCounter *int, storageRequestChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		TCP server to receive messages from peers
		To merge with ListenPeersRequestsTCP once tested and working
	*/

	listener, err := net.Listen("tcp", ":"+port)

	fmt.Println(ratiosForPeers)

	utils.ErrorHandler(err)

	defer listener.Close()
	for {
		failureMutex.Lock()
		conn, _ := listener.Accept()
		go handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratiosAtPeers, bytesForPeers, storedForPeers, factorAcceptableRatio, deletienQueue, storageRequestChannel, testRequestsChannel)
		failureMutex.Unlock()
	}
}

func ListenPeersRequestsTCP(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, storageRequestsChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		TCP server to receive messages from peers
	*/

	listener, err := net.Listen("tcp", ":"+port)

	utils.ErrorHandler(err)

	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratiosAtPeers, bytesForPeers, storedForPeers, factorAcceptableRatio, deletienQueue, storageRequestsChannel, testRequestsChannel)
	}
}

func handleConnection(conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, storageRequestsChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer conn.Close()

	buffer := make([]byte, 63)

	conn.Read(buffer)
	MessageDiscriminator(buffer, conn, nodeStorage, bytesAtPeers, scores, ratios, bytesForPeers, storedForPeers, factorAcceptableRatio, deletionQueue, storageRequestsChannel, testRequestsChannel)
}

func MessageDiscriminator(buffer []byte, conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, storageRequestsChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/

	bufferString := string(buffer)
	messageType := bufferString[:5]

	if messageType == "StoRq" {
		fmt.Println("Received storage request")
		// storagerequests.HandleStorageRequest(bufferString, conn, bytesForPeers, storedForPeers)
		storageRequest, err := buildStorageRequest(bufferString)
		if err == nil {
			queueMessage := datastructures.StorageRequestQueueMessage{StorageRequest: storageRequest, Conn: conn}
			storageRequestsChannel <- queueMessage
		}
	} else if messageType == "BarRq" {
		remoteAddr := conn.RemoteAddr()
		ip, _, err := net.SplitHostPort(remoteAddr.String())
		utils.ErrorHandler(err)
		fmt.Println("Received bartering request from peer", ip)
		bartering.RespondToBarterMsg(bufferString, ip, nodeStorage, bytesAtPeers, scores, conn, ratios, factorAcceptableRatio)
	} else if messageType == "TesRq" {
		CID := bufferString[5 : len(bufferString)-1]
		fmt.Println("Recieved test request for file ", CID)
		testRequest := datastructures.TestRequestQueueMessage{CID: CID, Conn: conn}
		// storagetesting.HandleTest(CID, conn)
		testRequestsChannel <- testRequest
	} else {
		fmt.Println("Unrecognized message : ", bufferString)
	}
}

func buildStorageRequest(bufferString string) (datastructures.StorageRequest, error) {

	CID := bufferString[5:51]

	fileSize := bufferString[51:]
	fileSize = strings.Split(fileSize, "\n")[0]
	fileSizeFloat, err := strconv.ParseFloat(fileSize, 64)

	if err != nil {
		return datastructures.StorageRequest{}, fmt.Errorf("could not parse file size ; storage request invalid")
	}

	return datastructures.StorageRequest{FileSize: fileSizeFloat, CID: CID}, nil

}
