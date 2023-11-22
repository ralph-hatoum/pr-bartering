package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bartering/utils"
)

/*
	Code for bootstrap node - here a simple HTTP server
	On 8082, bootstrap returns a list of peers as a list of ip addresses as strings
*/

func main() {

	fmt.Println("-- BOOTSTRAP NODE --")

	address := "localhost"
	port := "8082"

	fmt.Println("Listening on port ", port)

	serverAddress := address + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-- PEER CONNECTION -- HANDLING CONNECTION --")
		peers := []string{"134.214.43.12", "134.214.43.13", "134.214.43.14", "134.214.43.15", "127.0.0.1"}
		jsonResponse, err := json.Marshal(peers)
		utils.ErrorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

		fmt.Println("-- PEER CONNECTION HANDLED SUCCESFULLY --")
	})

	err := http.ListenAndServe(serverAddress, nil)
	utils.ErrorHandler(err)

	fmt.Println("Bootstrap server listening on port 8080")

}
