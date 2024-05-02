package api_ipfs

/*
Functions to interact with the IPFS Daemon
*/

import (
	"fmt"
	"os/exec"
	"strings"
)

func UploadToIPFS(path string) (string, error) {
	/*
		To add a file to IPFS
		Arguments : path to file as a string
		Returns : CID (Content Identifier) as a string
	*/

	cmd := "ipfs"
	cmdArgs := []string{"add", path}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	if err != nil {
		fmt.Println("ERROR : could not upload to IPFS")
		return "", fmt.Errorf("could not upload to IPFS")
	}

	CID := strings.Split(string(cmdOutput), " ")[1]

	return CID, nil

}

func PinToIPFS(cid string) (string, error) {
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

	if err != nil {
		fmt.Println("ERROR : could not pin to IPFS")
		return "", fmt.Errorf("could not pin to IPFS")
	}

	return string(cmdOutput), nil

}

func UnpinIPFS(cid string) (string, error) {
	/*
		To unpin a file to IPFS
		Arguments : CID (Content Identifier) of the file as a string
		Returns : output of the unpin command as a string
	*/

	cmd := "ipfs"
	cmdArgs := []string{"pin", "rm", cid}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	if err != nil {
		fmt.Println("ERROR : could not unpin to IPFS")
		return "", fmt.Errorf("could not unpin to IPFS")
	}

	return string(cmdOutput), nil
}

func CatIPFS(cid string) (string, error) {
	/*
		To cat a file on IPFS (see content)
		Arguments : CID (Content Identifier) of the file as a string
		Returns : output of the unpin command as a string
	*/
	fmt.Println("Calling cat command for CID", cid)
	cmd := "/usr/local/bin/ipfs"
	cmdArgs := []string{"cat", "--timeout=30s", cid}

	cmdOutput, err := exec.Command(cmd, cmdArgs...).Output()

	if err != nil {
		fmt.Println("ERROR : could not cat to IPFS")
		return "", fmt.Errorf("could not cat to IPFS")
	}

	return string(cmdOutput), nil
}
