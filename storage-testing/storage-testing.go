package storagetesting

import (
	api_ipfs "bartering/api-ipfs"
	"bartering/bartering-api"
	storagerequests "bartering/storage-requests"
	"bartering/utils"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var PORT = "8081"

var TIMER_TIMEOUT_SEC = 5

type ScoreVariationScenario struct {
	/*
		We setup scenarios such as "failedTestTimeout" or "FailedTestWrongAns" to decrease or increase scores
		differently given the situation
	*/
	Scenario  string
	Variation float64
}

var DecreasingBehavior = []ScoreVariationScenario{{Scenario: "failedTestTimeout", Variation: 0.5}, {Scenario: "failedTestWrongAns", Variation: 0.7}}

var IncreasingBehavior = []ScoreVariationScenario{{Scenario: "passedTest", Variation: 0.2}}

func RequestTest(CID string, filesAtPeers []storagerequests.FilesAtPeers, scores []bartering.NodeScore) {
	/*
		Function to request tests on a file stored at peers
		Arguments : CID (Content Identifier) of the file as a string, array of FilesAtPeers objects, array of NodeScore objects
	*/

	storers, err := findStorers(CID, filesAtPeers)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, storer := range storers {
		// maybe parallelize ?
		ContactPeerForTest(CID, storer, scores)
	}

}

func HandleTest(CID string, conn net.Conn) {
	/*
		Function to perform tests upon recieving a test request
		Arguments : CID as a string, connection as net.Conn
	*/

	answer := computeExpectedAnswer(CID)
	buffer := []byte(answer)
	conn.Write(buffer)

}

func ContactPeerForTest(CID string, peer string, scores []bartering.NodeScore) {
	/*
		Function to contact a peer to ask for a test, check answer and update score accordingly
		Arguments : CID of file to test a string, peer IP as string, scores as array of NodeScore objects
	*/

	conn, err := net.Dial("tcp", peer+":"+PORT)
	utils.ErrorHandler(err)

	defer conn.Close()

	message := "TesRq" + CID

	_, err = io.WriteString(conn, message)

	utils.ErrorHandler(err)

	responseChannel := make(chan string)

	go handleResponse(responseChannel, conn)

	timer := time.NewTimer(time.Duration(TIMER_TIMEOUT_SEC) * time.Second)

	select {
	case <-timer.C:
		fmt.Println("Timeout: No response received.")
		// Here, score should be decreased as no response was received
		decreaseScore(peer, "failedTestTimeout", scores)
	case response := <-responseChannel:
		fmt.Println("Response received.")
		// Here, response was received, it should be checked if the response is correct or wrong to decide how score should evolve
		fmt.Println(response)
		if checkAnswer(response, CID) {
			increaseScore(peer, "passedTest", scores)
		} else {
			decreaseScore(peer, "failedTestWrongAns", scores)
		}
	}
}

func handleResponse(responseChannel chan<- string, conn net.Conn) {
	/*
		Function to handle a response recieved when requesting a test from a peer
		Arguments : string chanel, connection as net.Conn
	*/

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	response := string(buffer[:n])
	responseChannel <- response
}

func findStorers(CID string, filesAtPeers []storagerequests.FilesAtPeers) ([]string, error) {
	/*
		Function to find peers storing a file for self
		Arguments : CID as string, array of FilesAtPeers objects
		Returns : list of string of peers storing given CID and nil if CID is indeed
		found at other peers, empty list and error if no peer is storing CID
	*/

	storers := []string{}
	for _, peerFiles := range filesAtPeers {
		if lookForFile(CID, peerFiles.Files) {
			storers = append(storers, peerFiles.Peer)
		}
	}

	if len(storers) == 0 {
		return storers, errors.New("no peers storing file with CID " + CID)
	} else {
		return storers, nil
	}
}

func lookForFile(CID string, fileList []string) bool {
	/*
		Function to look for a file in a file list (used in findStorers to check if a peer is storing a file of given CID)
		Arguments : CID of file as string, CID list as list of strings
		Return : True if file is found False otherwise
	*/

	for _, file := range fileList {
		if file == CID {
			return true
		}
	}
	return false
}

func findScoreVariation(variations []ScoreVariationScenario, scenario string) (float64, error) {
	/*
		Function to find how much a score should be decreased or increased given the situation
		Arguments : array of ScoreVariationScenario objects, scenario as a string
		Returns : score variation as float64 and nil if scenario is found, 0.0 and error if the given scenario is not found
	*/

	for _, variation := range variations {
		if variation.Scenario == scenario {
			return variation.Variation, nil
		}
	}
	return 0.0, errors.New("scenario " + scenario + " not found")
}

func computeExpectedAnswer(CID string) string {
	/*
		Given a CID, we compute the answer to a test (for now simple SHA256 hash but this will need to implement filecoin proof)
		Arguments : CID of file as string
		Returns : proof result as string
	*/

	contentString := api_ipfs.CatIPFS(CID)
	contentBytes := []byte(contentString)
	hasher := sha256.New()

	hasher.Write(contentBytes)
	proofResult := hasher.Sum(nil)

	return string(proofResult)
}

/* TODO : unifiy decrease and increase functions into a single update function, and
also unify decreasing and increasing behavior dics into one update doc with signed float64 values
*/
func decreaseScore(peer string, scenario string, scores []bartering.NodeScore) {
	/*
		Given a scenario, decrease a peer's score accordingly
		Arguments : peer IP as string, scenario as a string, scores as NodeScore objects
	*/

	decreaseAmount, err := findScoreVariation(DecreasingBehavior, scenario)
	utils.ErrorHandler(err)

	for _, peerScore := range scores {
		if peerScore.NodeIP == peer {
			peerScore.Score -= decreaseAmount
		}
	}

}

func increaseScore(peer string, scenario string, scores []bartering.NodeScore) {
	/*
		Given a scenario, increase a peer's score accordingly
		Arguments : peer IP as string, scenario as a string, scores as NodeScore objects
	*/

	increaseAmount, err := findScoreVariation(IncreasingBehavior, scenario)
	utils.ErrorHandler(err)

	for _, peerScore := range scores {
		if peerScore.NodeIP == peer {
			peerScore.Score += increaseAmount
		}
	}

}

func checkAnswer(answer string, CID string) bool {
	/*
		Check if the received answer to a test is valid
		Arguments : answer recieved as a string, CID of the file to test a string
	*/

	expectedAnswer := computeExpectedAnswer(CID)
	return expectedAnswer == answer
}
