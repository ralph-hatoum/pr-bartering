package metriccollector

import (
	"testing"
)

func TestExtractFailureModel(t *testing.T) {

	msgCounter, testCounter := InitiateCounters()

	IncreaseCounter(&msgCounter)
	IncreaseCounter(&testCounter)

	if msgCounter != 1 || testCounter != 1 {
		t.Errorf("counters not increased properly")
	}
}

func TestComputeTotalEnergyConsumption(t *testing.T) {

	singleTestEnergy := 3.0
	singleMsgEnergy := 5.0

	msgCounter := 200
	testCounter := 75

	networkEnergy, testEnergy, totalEnergy := ComputeTotalEnergyConsumption(testCounter, msgCounter, singleTestEnergy, singleMsgEnergy)

	if networkEnergy != 1000.0 || testEnergy != 225.0 || totalEnergy != 1225.0 {
		t.Errorf("energy not computed correctly")
	}

}
