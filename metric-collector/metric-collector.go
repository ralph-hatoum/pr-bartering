package metriccollector

// Not sure about all of this - should probabluy think it through before starting to code ...

func IncreaseCounter(counter *int) {
	*counter += 1
}

func IncreaseByteCounter(counter *int, bytes int) {
	*counter += bytes
}

func DecreaseByteCounter(counter *int, bytes int) {
	*counter -= bytes
}

func ComputeTotalEnergyConsumption(testCounter int, msgCounter int, singleTestEnergy float64, singleMsgEnergy float64) (float64, float64, float64) {
	// Function to compute energy consumption incurred by testing and network
	// requires singleTestEnergy in watts, singleMsgEnergy in watts

	networkEnergy := float64(msgCounter) * singleMsgEnergy
	testEnergy := float64(testCounter) * singleTestEnergy

	totalEnergy := networkEnergy + testEnergy

	return networkEnergy, testEnergy, totalEnergy

}
