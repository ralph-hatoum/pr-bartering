package metriccollector

import (
	"bartering/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// Metrics to keep track of : number of msg sent, number of tests performed, energy consumption incurred,
// number of confirmed replicas for each CID, time where nb of confirmed replicas under SLA,
//  CIDs stored, data (nb of bytes probably) stored on node, data (nb of bytes probably) stored at peers, ratio of bytes stored at over bytes stored for
// time to confirm storage

// Architecture : nodes collect metrics and "publish" results on an HTTP endpoint
// Data can be retrieved through Prometheus but also through any other way throught the HTTP endpoint

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

func PublishResultsHTTP(httpPort string) {
	address := "localhost"
	port := httpPort

	serverAddress := address + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		peers := []string{"127.0.0.1"}
		jsonResponse, err := json.Marshal(peers)
		utils.ErrorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	})

	err := http.ListenAndServe(serverAddress, nil)
	utils.ErrorHandler(err)

	fmt.Println("Bootstrap server listening on port 8080")

}
