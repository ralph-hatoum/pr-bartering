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
	Scenario  string
	Variation float64
}

var DecreasingBehavior = []ScoreVariationScenario{{Scenario: "failedTestTimeout", Variation: 0.5}, {Scenario: "failedTestWrongAns", Variation: 0.7}}

var IncreasingBehavior = []ScoreVariationScenario{{Scenario: "passedTest", Variation: 0.2}}

func RequestTest(CID string, filesAtPeers []storagerequests.FilesAtPeers, scores []bartering.NodeScore) {

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

	answer := computeExpectedAnswer(CID)
	buffer := []byte(answer)
	conn.Write(buffer)

}

func ContactPeerForTest(CID string, peer string, scores []bartering.NodeScore) {
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
	for _, file := range fileList {
		if file == CID {
			return true
		}
	}
	return false
}

func findScoreVariation(variations []ScoreVariationScenario, scenario string) (float64, error) {
	for _, variation := range variations {
		if variation.Scenario == scenario {
			return variation.Variation, nil
		}
	}
	return 0.0, errors.New("scenario " + scenario + " not found")
}

func computeExpectedAnswer(CID string) string {

	contentString := api_ipfs.CatIPFS(CID)
	contentBytes := []byte(contentString)
	hasher := sha256.New()

	hasher.Write(contentBytes)
	proofResult := hasher.Sum(nil)

	return string(proofResult)
}

func decreaseScore(peer string, scenario string, scores []bartering.NodeScore) {

	decreaseAmount, err := findScoreVariation(DecreasingBehavior, scenario)
	utils.ErrorHandler(err)

	for _, peerScore := range scores {
		if peerScore.NodeIP == peer {
			peerScore.Score -= decreaseAmount
		}
	}

}

func increaseScore(peer string, scenario string, scores []bartering.NodeScore) {

	increaseAmount, err := findScoreVariation(IncreasingBehavior, scenario)
	utils.ErrorHandler(err)

	for _, peerScore := range scores {
		if peerScore.NodeIP == peer {
			peerScore.Score += increaseAmount
		}
	}

}

func checkAnswer(answer string, CID string) bool {
	expectedAnswer := computeExpectedAnswer(CID)
	return expectedAnswer == answer
}
