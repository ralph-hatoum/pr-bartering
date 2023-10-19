package peersconnect

import (
	"fmt"
	"net"

	"../utils"
)

var PORT = "8081"

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

	messageDiscriminator(buffer)

}

func messageDiscriminator(buffer []byte) {
	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/
	bufferString := string(buffer)

	messageType := bufferString[:5]

	if messageType == "StoRq" {
		handleStorageRequest(bufferString)
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
