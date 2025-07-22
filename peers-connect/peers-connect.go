package peersconnect

import (
	"net"

	worker "bartering/worker"

	"go.uber.org/zap"
)

func ListenPeersRequestsTCP(port string, workloadChannel chan worker.Workload, logger *zap.Logger) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("Could not start peer listener :", zap.Error(err))
		return
	}

	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		// TODO: protection : be able to block, set cooldowns, etc
		go handleConnection(conn, workloadChannel, logger)
	}
}

func handleConnection(conn net.Conn, workloadChannel chan worker.Workload, logger *zap.Logger) {
	defer conn.Close()

	buffer := make([]byte, 64)
	_, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Could not read from connection :", zap.Error(err))
		return
	}

	workloadType := buffer[0]
	payloadBytes := buffer[1:33]

	payloadString := string(payloadBytes)

	workload := worker.Workload{Cid: payloadString, WorkloadType: workloadType, PeerIP: conn.RemoteAddr().String()}

	workloadChannel <- workload
}
