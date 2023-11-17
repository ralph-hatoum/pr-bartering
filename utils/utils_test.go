package utils

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestListPrint(t *testing.T) {
	list := []string{"hello", "goodbye"}
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ListPrint(list)

	w.Close()
	os.Stdout = originalStdout

	capturedOutput, _ := ioutil.ReadAll(r)

	expected := "hello goodbye"

	if !strings.Contains(string(capturedOutput), expected) {
		t.Errorf("Expected output: %s, got: %s", expected, capturedOutput)
	}

}
