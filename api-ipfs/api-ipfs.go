package api_ipfs


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



