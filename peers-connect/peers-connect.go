package peersconnect

import (
	"fmt"
	"net"

	"bartering/bartering-api"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
	"bartering/utils"
)

// var PORT = "8081"

// var NodeStorageSpace = 400988.45 * 1000

// var bytesAtPeers = []bartering.PeerStorageUse{{NodeIP: "127.0.0.1", StorageAtNode: 400988.45}}

// var scores = []bartering.NodeScore{{NodeIP: "127.0.0.1", Score: 10.0}}

func ListenPeersRequestsTCP(port string, nodeStorage float64, bytesAtPeers []bartering.PeerStorageUse, scores []bartering.NodeScore, ratios []bartering.NodeRatio) {
	/*
		TCP server to receive messages from peers
	*/
	listener, err := net.Listen("tcp", ":"+port)

	utils.ErrorHandler(err)

	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratios)
	}
}

func handleConnection(conn net.Conn, nodeStorage float64, bytesAtPeers []bartering.PeerStorageUse, scores []bartering.NodeScore, ratios []bartering.NodeRatio) {
	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer conn.Close()

	buffer := make([]byte, 63)

	conn.Read(buffer)

	MessageDiscriminator(buffer, conn, nodeStorage, bytesAtPeers, scores, ratios)

}

func MessageDiscriminator(buffer []byte, conn net.Conn, nodeStorage float64, bytesAtPeers []bartering.PeerStorageUse, scores []bartering.NodeScore, ratios []bartering.NodeRatio) {
	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/
	bufferString := string(buffer)

	messageType := bufferString[:5]

	if messageType == "StoRq" {

		storagerequests.HandleStorageRequest(bufferString)

	} else if messageType == "BarRq" {

		fmt.Println("Received bartering request")
		// Maybe check if peer is known ???

		remoteAddr := conn.RemoteAddr()
		ip, _, err := net.SplitHostPort(remoteAddr.String())
		utils.ErrorHandler(err)

		fmt.Println("Received bartering request from peer", ip)

		bartering.RespondToBarterMsg(bufferString, ip, nodeStorage, bytesAtPeers, scores, conn, ratios)

	} else if messageType == "TesRq" {
		fmt.Println("Recieved test request")
		CID := bufferString[5 : len(bufferString)-1]
		fmt.Println(CID)
		storagetesting.HandleTest(CID, conn)

	} else {
		fmt.Println("Unrecognized message : ", bufferString)
	}
}
