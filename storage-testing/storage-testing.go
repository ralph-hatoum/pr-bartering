package storagetesting

import (
	api_ipfs "bartering/api-ipfs"
	datastructures "bartering/data-structures"
	storagerequests "bartering/storage-requests"
	"bartering/utils"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func PeriodicTests(fulfilledRequests *[]datastructures.FulfilledRequest, scores []datastructures.NodeScore, timerTimeoutSec float64, port string, testingPeriod float64, DecreasingBehavior []datastructures.ScoreVariationScenario, IncreasingBehavior []datastructures.ScoreVariationScenario, bytesAtPeers []datastructures.PeerStorageUse, scoreDecreaseRefStoReq float64) {

	/*
		Function to requests tests periodically from peers storing our data
		Arguments : FulfilledRequest array, NodeScore array
	*/

	fmt.Println("Periodic Tester started!")

	for {
		fmt.Println("before sleep")
		time.Sleep(time.Duration(testingPeriod) * time.Second)
		fmt.Println("after sleep")
		fmt.Println(fulfilledRequests)
		for _, fulfilledRequest := range *fulfilledRequests {
			if len(*fulfilledRequests) == 0 {
				fmt.Println("No tests to do")
			}
			testResult := ContactPeerForTest(fulfilledRequest.CID, fulfilledRequest.Peer, scores, timerTimeoutSec, port, DecreasingBehavior, IncreasingBehavior)
			if !testResult {
				// Could not confirm storage ; need to request storage from other node
				fmt.Println("requesting storage from other node ... ")
				stoReq := datastructures.StorageRequest{CID: fulfilledRequest.CID, FileSize: fulfilledRequest.FileSize}
				peersToRq := storagerequests.RemovePeerFromPeers(scores, fulfilledRequest.Peer)
				storagerequests.StoreKCopiesOnNetwork(peersToRq, 1, stoReq, port, bytesAtPeers, fulfilledRequests, scoreDecreaseRefStoReq)
			}
		}
	}
}

func RequestTest(CID string, filesAtPeers []datastructures.FilesAtPeers, scores []datastructures.NodeScore, timerTimeoutSec float64, port string, DecreasingBehavior []datastructures.ScoreVariationScenario, IncreasingBehavior []datastructures.ScoreVariationScenario) {

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
		ContactPeerForTest(CID, storer, scores, timerTimeoutSec, port, DecreasingBehavior, IncreasingBehavior)
	}

}

func HandleTest(CID string, conn net.Conn) {

	/*
		Function to perform tests upon recieving a test request
		Arguments : CID as a string, connection as net.Conn
	*/

	answer := computeExpectedAnswer(CID)
	fmt.Println("Proof computed : ", answer)
	buffer := []byte(answer)
	conn.Write(buffer) // INCREASE NBMSG COUNTER

}

func ContactPeerForTest(CID string, peer string, scores []datastructures.NodeScore, timerTimeoutSec float64, port string, DecreasingBehavior []datastructures.ScoreVariationScenario, IncreasingBehavior []datastructures.ScoreVariationScenario) bool {
    conn, err := net.Dial("tcp", peer+":"+port)
    utils.ErrorHandler(err)
    defer conn.Close()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    message := "TesRq" + CID
    _, err = io.WriteString(conn, message)
    utils.ErrorHandler(err)

    responseChannel := make(chan string)
    var wg sync.WaitGroup
    wg.Add(1)

    go handleResponse(ctx, &wg, responseChannel, conn)

    timer := time.NewTimer(time.Duration(timerTimeoutSec) * time.Second)
    defer timer.Stop()

    defer wg.Wait()  // Ensures `wg.Wait()` is called before function exit

    select {
    case <-timer.C:
        fmt.Println("Timeout: No response received.")
        decreaseScore(peer, "failedTestTimeout", scores, DecreasingBehavior)
        cancel() // Cancel the context to signal handleResponse
        return false
    case response := <-responseChannel:
        if checkAnswer(response, CID) {
            fmt.Println("test passed")
            increaseScore(peer, "passedTest", scores, IncreasingBehavior)
            return true
        } else {
            fmt.Println("test not passed")
            decreaseScore(peer, "failedTestWrongAns", scores, DecreasingBehavior)
            return false
        }
    }
}


func handleResponse(ctx context.Context, wg *sync.WaitGroup, responseChannel chan<- string, conn net.Conn) {
	defer wg.Done()
	defer close(responseChannel)

	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, exiting handleResponse")
			return
		default:
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading response:", err)
				return
			}
			response := string(buffer[:n])
			responseChannel <- response
			return // Successfully read response, exit goroutine
		}
	}
}

func findStorers(CID string, filesAtPeers []datastructures.FilesAtPeers) ([]string, error) {

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

func findScoreVariation(variations []datastructures.ScoreVariationScenario, scenario string) (float64, error) {

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

func computeExpectedAnswer(CID string) []byte {

	/*
		Given a CID, we compute the answer to a test (for now simple SHA256 hash but this will need to implement filecoin proof)
		Arguments : CID of file as string
		Returns : proof result as string
	*/

	CID = CID[:46]
	contentString := api_ipfs.CatIPFS(CID)
	contentBytes := []byte(contentString)
	hasher := sha256.New()

	hasher.Write(contentBytes)
	proofResult := hasher.Sum(nil)

	return proofResult
}

/* TODO : unifiy decrease and increase functions into a single update function, and
also unify decreasing and increasing behavior dics into one update doc with signed float64 values
*/

func decreaseScore(peer string, scenario string, scores []datastructures.NodeScore, DecreasingBehavior []datastructures.ScoreVariationScenario) {

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

func increaseScore(peer string, scenario string, scores []datastructures.NodeScore, IncreasingBehavior []datastructures.ScoreVariationScenario) {

	/*
		Given a scenario, increase a peer's score accordingly
		Arguments : peer IP as string, scenario as a string, scores as NodeScore objects
	*/

	increaseAmount, err := findScoreVariation(IncreasingBehavior, scenario)
	utils.ErrorHandler(err)

	for index, peerScore := range scores {
		if peerScore.NodeIP == peer {
			scores[index].Score += increaseAmount
		}
	}

}

func checkAnswer(answer string, CID string) bool {

	/*
		Check if the received answer to a test is valid
		Arguments : answer recieved as a string, CID of the file to test a string
	*/

	expectedAnswer := computeExpectedAnswer(CID)
	return string(expectedAnswer) == answer
}
