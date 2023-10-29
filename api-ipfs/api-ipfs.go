package api_ipfs

/*
Functions to interact with the IPFS Daemon
*/

import (
	"os/exec"
	"strings"

	"bartering/utils"
)

func UploadToIPFS(path string) string {
	/*
		To add a file to IPFS
		Arguments : path to file as a string
		Returns : CID (Content Identifier) as a string
	*/

	cmd := "ipfs"
	cmdArgs := []string{"add", path}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	utils.ErrorHandler(err)

	CID := strings.Split(string(cmdOutput), " ")[1]

	return CID

}

func PinToIPFS(cid string) string {
	/*
		To pin a file to IPFS
		Arguments : CID (Content Identifier) of the file as a string
		Returns : output of the pin command as a string

		By default, the "ipfs add" command automatically pins the file, therefore this
		function does not need to be called if the UploadToIPFS function was called
	*/

	cmd := "ipfs"
	cmdArgs := []string{"pin", "add", cid}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	utils.ErrorHandler(err)

	return string(cmdOutput)

}

func UnpinIPFS(cid string) string {
	/*
		To unpin a file to IPFS
		Arguments : CID (Content Identifier) of the file as a string
		Returns : output of the unpin command as a string
	*/

	cmd := "ipfs"
	cmdArgs := []string{"pin", "rm", cid}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	utils.ErrorHandler(err)

	return string(cmdOutput)
}
