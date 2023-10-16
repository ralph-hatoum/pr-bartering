package api_ipfs

import (
	"fmt"
	"os/exec"
)

func uploadToIPFS(path string) string {
	cmd := "ipfs"
	args := []string{"add", path}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)

}

func pinToIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "add", cid}

	output, err := exec.Command(cmd, args...).Output()

	errorHandler(err)

	return string(output)

}

func unpinIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "rm", cid}

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
