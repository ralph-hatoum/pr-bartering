package bootstrapconnect

/*
Functions to interact with the bootstrap node of the network
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"bartering/utils"
)

func GetPeersFromBootstrapTCP(IP string, port string) string {

	/*
		Function to get peers from the bootstrap node via TCP
		Arguments : IP of bootsrap as string, port as string
		Returns : bootstrap's response as string
	*/

	serverAddress := IP + ":" + port
	conn, err := net.Dial("tcp", serverAddress)
	utils.ErrorHandler(err)

	defer conn.Close()

	messageToBootstrap := "hello\n"

	_, err = io.WriteString(conn, messageToBootstrap) // INCREASE NBMSG COUNTER
	utils.ErrorHandler(err)

	boostrapResponseReader := bufio.NewReader(conn)

	boostrapResponse, err := boostrapResponseReader.ReadString('\n')
	utils.ErrorHandler(err)

	return boostrapResponse
}

func GetPeersFromBootstrapHTTP(IP string, port string) []string {

	/*
		Function to get peers from the bootstrap node via HTTP
		Arguments : IP of bootsrap as string, port as string
		Returns : bootstrap's response as string
	*/

	bootstrapUrl := IP + ":" + port

	bootstrapResponse, err := http.Get("http://" + bootstrapUrl)
	utils.ErrorHandler(err)

	defer bootstrapResponse.Body.Close()

	if bootstrapResponse.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code:", bootstrapResponse.StatusCode)
		panic(-1)
	}

	bootstrapResponseBody, err := ioutil.ReadAll(bootstrapResponse.Body)
	utils.ErrorHandler(err)

	var peers []string

	err = json.Unmarshal(bootstrapResponseBody, &peers)

	utils.ErrorHandler(err)

	return peers

}

func AnnounceSelfToBootstrap(IP string, port string) {

	/*
		Function to call the bootstrap to annouce self and add IP to the IPs that will be announced by the bootstrap
		Arguments : IP as string, port as string
	*/

}
