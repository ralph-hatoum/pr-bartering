package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bartering/utils"
	"bufio"
)

/*
	Code for bootstrap node - here a simple HTTP server
	On 8082, bootstrap returns a list of peers as a list of ip addresses as strings
*/

func main() {

	// args := os.Args

	// if len(args) != 2 {
	// 	fmt.Println("Missing bootstrap IP")
	// 	panic(-1)
	// }

	fmt.Println("-- BOOTSTRAP NODE --")

	address := "0.0.0.0"
	port := "8082"

	fmt.Println("Listening on port ", port)

	// Build IP peers array

	peers, err := BuildPeersIPlist("./ips.txt")

	if err != nil {
		fmt.Println("Error building peer IP list")
		return
	}

	serverAddress := address + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-- PEER CONNECTION -- HANDLING CONNECTION --")
		// peers := []string{"134.214.202.223", "134.214.202.224"}
		jsonResponse, err := json.Marshal(peers)
		utils.ErrorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

		fmt.Println("-- PEER CONNECTION HANDLED SUCCESFULLY --")
	})

	err = http.ListenAndServe(serverAddress, nil)
	utils.ErrorHandler(err)

}

func BuildPeersIPlist(path string) ([]string, error) {

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return []string{}, err
	}

	defer file.Close()

	var peers []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		peers = append(peers, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return []string{}, err
	}

	return peers, nil

}
