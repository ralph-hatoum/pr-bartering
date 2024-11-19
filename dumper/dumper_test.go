package dumper

import (
	datastructures "bartering/data-structures"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func init() {
	bytesAtPeers := []datastructures.PeerStorageUse{{NodeIP: "peer1", StorageAtNode: 21.0}}
	bytesForPeers := []datastructures.PeerStorageUse{{NodeIP: "peer2", StorageAtNode: 23.0}}
	fulfilledRequests := []datastructures.FulfilledRequest{{CID: "TestCID", Peer: "peer2", FileSize: 10.0}}
	storagePool := []string{"pool1", "pool2"}
	pendingRequests := []datastructures.StorageRequest{{CID: "TestCID", FileSize: 11.0}}
	peers := []string{"peer1", "peer2", "peer3"}
	scores := []datastructures.NodeScore{{NodeIP: "testIP", Score: 10.0}}
	ratiosForPeers := []datastructures.NodeRatio{{NodeIP: "testIP", Ratio: 1.0}}
	ratiosAtPeers := []datastructures.NodeRatio{{NodeIP: "testIP", Ratio: 2.0}}
	storedForPeers := []datastructures.FulfilledRequest{{CID: "TestCID", Peer: "peer2", FileSize: 10.0}, {CID: "TestCID2", Peer: "peer3", FileSize: 13.0}}

	go Dumper(bytesAtPeers, bytesForPeers, fulfilledRequests, storagePool, pendingRequests, peers, scores, ratiosForPeers, ratiosAtPeers, storedForPeers)
}

func Test_Dumper(t *testing.T) {
	dumperResponse, err := http.Get(fmt.Sprintf("http://%s:%s/%s", "localhost", "8083", "dump"))
	if err != nil {
		t.FailNow()
	}
	defer dumperResponse.Body.Close()

	dumperResponseBody, _ := io.ReadAll(dumperResponse.Body)

	var output map[string]any
	err = json.Unmarshal(dumperResponseBody, &output)
	if err != nil {
		t.FailNow()
	}
	if dumperResponse.StatusCode != 200 {
		t.FailNow()
	}
}
