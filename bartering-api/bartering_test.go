package bartering

import (
	datastructures "bartering/data-structures"
	"testing"
)

// TODO Write tests

func TestFindNodeRatio(t *testing.T) {
	ratios := []datastructures.NodeRatio{
		{NodeIP: "peer1", Ratio: 0.5},
		{NodeIP: "peer2", Ratio: 0.8},
	}

	ratio, _ := FindNodeRatio(ratios, "peer1")
	if ratio != 0.5 {
		t.Errorf("Expected ratio 0.5, got %v", ratio)
	}
}

func TestUpdatePeerRatio(t *testing.T) {
	ratios := []datastructures.NodeRatio{
		{NodeIP: "peer1", Ratio: 0.5},
		{NodeIP: "peer2", Ratio: 0.8},
	}
	updatePeerRatio(ratios, "peer2", 0.9)
	updatedRatio, _ := FindNodeRatio(ratios, "peer2")
	if updatedRatio != 0.9 {
		t.Errorf("Expected ratio: 0.9, got: %v", updatedRatio)
	}
}

func TestInitNodeScores(t *testing.T) {

}

func TestElectStorageNodes(t *testing.T) {

}

func TestCheckCIDValidity(t *testing.T) {

}

func TestCheckFileSizeValidity(t *testing.T) {

}

func TestCheckEnoughSpace(t *testing.T) {

}

func TestDealWithRefusedRequest(t *testing.T) {

}

func TestIncreaseTolerance(t *testing.T) {

}

func TestDecreaseTolerance(t *testing.T) {

}

func TestShouldReqBeAccepted(t *testing.T) {

}
