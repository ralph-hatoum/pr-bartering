package main

import (
	api_ipfs "bartering/api-ipfs"
	"fmt"
)

func main() {
	// fmt.Println("Starting response receiver")
	fmt.Println(api_ipfs.CatIPFS("QmV9tSDx9UiPeWExXEeH6aoDvmihvx6jD5eLb4jbTaKGps"))
}

