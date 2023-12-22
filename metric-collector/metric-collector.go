package metriccollector

func InitiateCounters() (int, int) {
	// Function to iniate counters to grab metrics
	// We need to keep track of : number of messages sent and number of tests performed

	return 0, 0
}

func IncreaseCounter(counter *int) {
	*counter += 1
}

func ComputeTotalEnergyConsumption(testCounter int, msgCounter int, singleTestEnergy float64, singleMsgEnergy float64) (float64, float64, float64) {
	// Function to compute energy consumption incurred by testing and network
	// requires singleTestEnergy in watts, singleMsgEnergy in watts

	networkEnergy := float64(msgCounter) * singleMsgEnergy
	testEnergy := float64(testCounter) * singleTestEnergy

	totalEnergy := networkEnergy + testEnergy

	return networkEnergy, testEnergy, totalEnergy

}
