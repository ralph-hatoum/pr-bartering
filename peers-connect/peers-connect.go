package peersconnect

import (
	"fmt"
	"net"
	"sync"
	"time"

	"bartering/bartering-api"
	datastructures "bartering/data-structures"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
	"bartering/utils"
)

func ListenPeersRequestsTCPFailure(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, failureMutex *sync.Mutex, msgCounter *int) {

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
		go handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratiosAtPeers, bytesForPeers, storedForPeers, factorAcceptableRatio, deletienQueue, msgCounter)
		failureMutex.Unlock()
	}
}

func ListenPeersRequestsTCP(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		TCP server to receive messages from peers
	*/
	maxConnections := 4 // Maximum number of concurrent connections
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		utils.ErrorHandler(err)
		return
	}
	defer listener.Close()

	// Create a channel to limit the number of concurrent connections
	connLimiter := make(chan struct{}, maxConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		connLimiter <- struct{}{} // Block if the limit is reached
		go func(conn net.Conn) {
			defer func() {
				<-connLimiter // Release the slot
				conn.Close()
			}()

			// Set TCP keep-alive
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				tcpConn.SetKeepAlive(true)
				tcpConn.SetKeepAlivePeriod(3 * time.Minute) // Keep-alive period
			}

			handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratiosAtPeers, bytesForPeers, storedForPeers, factorAcceptableRatio, deletienQueue, msgCounter)
		}(conn)
	}
}

func handleConnection(conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer conn.Close()

	buffer := make([]byte, 63)

	conn.Read(buffer)

	MessageDiscriminator(buffer, conn, nodeStorage, bytesAtPeers, scores, ratios, bytesForPeers, storedForPeers, factorAcceptableRatio, deletionQueue, msgCounter)

}

func MessageDiscriminator(buffer []byte, conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/

	bufferString := string(buffer)
	messageType := bufferString[:5]

	if messageType == "StoRq" {

		storagerequests.HandleStorageRequest(bufferString, conn, bytesForPeers, storedForPeers)
		// storagerequests.HandleStorageRequestTimed(bufferString, conn, bytesAtPeers, storedForPeers, deletionQueue)

	} else if messageType == "BarRq" {

		remoteAddr := conn.RemoteAddr()
		ip, _, err := net.SplitHostPort(remoteAddr.String())
		utils.ErrorHandler(err)
		fmt.Println("Received bartering request from peer", ip)
		bartering.RespondToBarterMsg(bufferString, ip, nodeStorage, bytesAtPeers, scores, conn, ratios, factorAcceptableRatio, msgCounter)

	} else if messageType == "TesRq" {

		fmt.Println("Recieved test request")
		CID := bufferString[5 : len(bufferString)-1]
		fmt.Println(CID)
		storagetesting.HandleTest(CID, conn)

	} else {

		fmt.Println("Unrecognized message : ", bufferString)

	}
}
