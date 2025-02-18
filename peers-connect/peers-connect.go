package peersconnect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	datastructures "bartering/data-structures"
)

func ListenPeersRequestsTCP(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, storageRequestsChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		TCP server to receive messages from peers
	*/

	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println("Could not start peer listener :", err)
		return
	}

	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handleConnection(conn, storageRequestsChannel, testRequestsChannel)
	}
}

func handleConnection(conn net.Conn, storageRequestsChannel chan datastructures.StorageRequestQueueMessage, testRequestsChannel chan datastructures.TestRequestQueueMessage) {

	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer conn.Close()

	buffer := make([]byte, 64)

	_, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Could not read from connection :", err)
		return
	}

	msgType := buffer[0]
	payload := buffer[1:]

	switch msgType {
	case 0x01:
		fmt.Println("Received storage request")
		cid := payload[:33]
		fileSize := payload[33:]
		var fileSizeFloat float64
		fileSizeBuf := bytes.NewReader((fileSize))
		err = binary.Read(fileSizeBuf, binary.BigEndian, &fileSizeFloat)
		if err != nil {
			fmt.Println("Could not read file size from buffer :", err)
			return
		}
		storageRequest := datastructures.StorageRequest{FileSize: fileSizeFloat, CID: string(cid)}
		storageRequestsChannel <- datastructures.StorageRequestQueueMessage{StorageRequest: storageRequest, Conn: conn}

	case 0x02:
		fmt.Println("Received test request")
		cid := payload[:33]
		testRequestsChannel <- datastructures.TestRequestQueueMessage{CID: string(cid), Conn: conn}
	}
}
