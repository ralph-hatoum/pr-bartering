package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../utils"
)

func main() {
	address := "localhost"
	port := "8080"

	serverAddress := address + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		peers := []string{"134.214.43.12", "134.214.43.13", "134.214.43.14", "134.214.43.15"}
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
