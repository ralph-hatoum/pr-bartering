package bootstrap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func errorHandler(err error) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		panic(0)
	}
}

func main() {
	address := "localhost"
	port := "8080"

	serverAddress := address + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		peers := []string{"134.214.43.12", "134.214.43.13", "134.214.43.14", "134.214.43.15"}
		jsonResponse, err := json.Marshal(peers)
		errorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	})
	err := http.ListenAndServe(serverAddress, nil)
	errorHandler(err)

	fmt.Println("Bootstrap server listening on port 8080")

}
