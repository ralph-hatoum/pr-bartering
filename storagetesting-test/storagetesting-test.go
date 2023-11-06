package main

import (
	storagetesting "bartering/storage-testing"
	"fmt"
)

func main() {
	fmt.Println("Starting response receiver")
	storagetesting.ContactPeerForTest("CID", "127.0.0.1")
}
