package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"bartering/utils"
	"bufio"
	"strings"
)

/*
	Code for bootstrap node - here a simple HTTP server
	On 8082, bootstrap returns a list of peers as a list of ip addresses as strings
*/



func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Println("Missing bootstrap IP")
		panic(-1)
	}

	fmt.Println("-- BOOTSTRAP NODE --")

	address := args[1]
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

		// Identify the client's IP address
		clientIP := strings.Split(r.RemoteAddr, ":")[0]

		// Filter out the client's IP from the list of peers if present
		filteredPeers := []string{}
		for _, peer := range peers {
			if peer != clientIP {
				filteredPeers = append(filteredPeers, peer)
			}
		}

		jsonResponse, err := json.Marshal(filteredPeers)
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
		peerAddress := scanner.Text()

		// Remove quotes and parentheses
		peerAddress = strings.Trim(peerAddress, "\"()")
		

		
		peers = append(peers, peerAddress)
		
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return []string{}, err
	}

	return peers, nil

}
