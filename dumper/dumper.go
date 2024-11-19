package dumper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Datastructure struct {
	Name  string
	Value any
}

func Dumper(datastructures []Datastructure) {
	address := "0.0.0.0"
	port := "8083"

	serverAddr := address + ":" + port

	http.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {
		response, err := responseBuilder(datastructures)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	})

	err := http.ListenAndServe(serverAddr, nil)

	if err != nil {
		fmt.Println("Could not start dumper")
	}
}

func responseBuilder(datastructures []Datastructure) ([]byte, error) {
	output := make(map[string]any)

	for _, datastructure := range datastructures {
		output[datastructure.Name] = datastructure.Value
	}

	jsonResponse, err := json.Marshal(output)

	if err != nil {
		fmt.Println("Could not build /dump response")
		return nil, fmt.Errorf("could not build /dump resp")
	}

	return jsonResponse, nil
}
