package functions

import (
	"bartering/utils"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	result := utils.GetFileSize("../test-data/test.txt")
	if result != 0.0126953125 {
		t.Errorf("Expected 0.0126953125, but got %f", result)
	}
}

// func TestListPrint(t *testing.T) {
// 	list_to_test := []string{"hello", "hi"}

// }
