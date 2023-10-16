package api_ipfs

import (
	"os/exec"

	"../utils"
)

func UploadToIPFS(path string) string {
	cmd := "ipfs"
	args := []string{"add", path}

	output, err := exec.Command(cmd, args...).Output()

	utils.ErrorHandler(err)

	return string(output)

}

func PinToIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "add", cid}

	output, err := exec.Command(cmd, args...).Output()

	utils.ErrorHandler(err)

	return string(output)

}

func UnpinIPFS(cid string) string {
	cmd := "ipfs"
	args := []string{"pin", "rm", cid}

	output, err := exec.Command(cmd, args...).Output()

	utils.ErrorHandler(err)

	return string(output)
}
