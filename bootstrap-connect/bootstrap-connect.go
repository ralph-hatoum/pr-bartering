package bootstrapconnect

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"../utils"
)

func GetPeersFromBootstrapTCP(IP string, port string) {

	serverAddress := IP + ":" + port
	conn, err := net.Dial("tcp", serverAddress)
	utils.ErrorHandler(err)

	defer conn.Close()

	message := "hello\n"

	_, err = io.WriteString(conn, message)
	utils.ErrorHandler(err)

	fmt.Println("Called bootstrap, awaiting response")

	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	utils.ErrorHandler(err)
	fmt.Println(response)

}

func GetPeersFromBootstrapHTTP(IP string, port string) []string {
	serverUrl := IP + ":" + port

	response, err := http.Get("http://" + serverUrl)
	utils.ErrorHandler(err)

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code:", response.StatusCode)
		panic(-1)
	}

	body, err := ioutil.ReadAll(response.Body)
	utils.ErrorHandler(err)

	var peers []string

	err = json.Unmarshal(body, &peers)

	utils.ErrorHandler(err)

	return peers

}
