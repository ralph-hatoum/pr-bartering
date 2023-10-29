package peersconnect

import (
	"fmt"
	"net"

	"bartering/bartering-api"
	"bartering/utils"
)

var PORT = "8081"

var NodeStorageSpace = 400988.45 * 1000

var bytesAtPeers = []bartering.PeerStorageUse{{NodeIP: "127.0.0.1", StorageAtNode: 400988.45}}

var scores = []bartering.NodeScore{{NodeIP: "127.0.0.1", Score: 10.0}}

func ListenPeersRequestsTCP() {
	/*
		TCP server to receive messages from peers
	*/
	listener, err := net.Listen("tcp", ":"+PORT)

	utils.ErrorHandler(err)

	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer conn.Close()

	buffer := make([]byte, 63)

	conn.Read(buffer)

	fmt.Println("Recevied message : ", string(buffer))

	MessageDiscriminator(buffer, conn)

}

func MessageDiscriminator(buffer []byte, conn net.Conn) {
	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/
	bufferString := string(buffer)

	messageType := bufferString[:5]

	fmt.Println(messageType)

	if messageType == "StoRq" {

		handleStorageRequest(bufferString)

	} else if messageType == "BarRq" {

		// fmt.Println("Received bartering request")
		// Maybe check if peer is known ???

		remoteAddr := conn.RemoteAddr()
		ip, _, err := net.SplitHostPort(remoteAddr.String())
		utils.ErrorHandler(err)

		fmt.Println("Received bartering request from peer", ip)

		bartering.RespondToBarterMsg(bufferString, ip, float64(NodeStorageSpace), bytesAtPeers, scores)

	} else {
		fmt.Println("Unrecognized message")
	}
}

func handleStorageRequest(bufferString string) {
	/*
		Function to handle a storage message type message
		Arguments : buffer received through a tcp connection, as a string
	*/

	fmt.Println("Received storage request")
	CID := bufferString[5:51]
	fmt.Println("CID : ", CID)
	fileSize := bufferString[51:]
	fmt.Println("File Size : ", fileSize)
}
