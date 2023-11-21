package functions

import (
	"bartering/bartering-api"
	"bartering/utils"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	result := utils.GetFileSize("../test-data/test.txt")
	if result != 0.0126953125 {
		t.Errorf("Expected 0.0126953125, but got %f", result)
	}
}

func TestInitiateBytesAtPeers(t *testing.T) {
	peers := []string{"peer1", "peer2"}
	storageAtPeer1 := bartering.PeerStorageUse{NodeIP: "peer1", StorageAtNode: 0.0}
	storageAtPeer2 := bartering.PeerStorageUse{NodeIP: "peer2", StorageAtNode: 0.0}
	result := initiatePeerStorageUseArray(peers, 0.0)
	if result[0] != storageAtPeer1 || result[1] != storageAtPeer2 {
		t.Errorf("BytesAtPeers not initiated correctly")
	}
}

func TestInitiateScores(t *testing.T) {
	peers := []string{"peer1", "peer2"}
	peerScore1 := bartering.NodeScore{NodeIP: "peer1", Score: 10.0}
	peerScore2 := bartering.NodeScore{NodeIP: "peer2", Score: 10.0}

	result := initiateScores(peers, 10.0)

	if result[0] != peerScore1 || result[1] != peerScore2 {
		t.Errorf("Scores not initiated correctly")
	}
}
