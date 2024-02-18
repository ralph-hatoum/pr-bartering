package main

import (
	"fmt"
	"testing"
)

func TestBuildPeersIPlist(t *testing.T) {
	peers, err := BuildPeersIPlist("ips.txt")
	if err != nil {
		fmt.Println("failed test")
		return
	}

	fmt.Println(peers)
}
