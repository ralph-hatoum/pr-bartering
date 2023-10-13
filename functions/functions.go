package functions

import (
	"fmt"
	"os/exec"
	"strings"
)

func Store(path string) {
	fmt.Println("Store")

	upload_command_result := uploadToIPFS(path)

	CID := strings.Split(upload_command_result, " ")[1]

	fmt.Println(CID)

	// TODO add to storage pool
	// TODO build storage request
	// TODO propagate to network
}

func uploadToIPFS(path string) string {
	cmd := "ipfs"
	args := []string{"add", path}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)

}

func errorHandler(err error) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		panic(0)
	}
}
