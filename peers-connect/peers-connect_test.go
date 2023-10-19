package peersconnect

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestListenPeersRequestsTCP(t *testing.T) {

}

func TestHandleConnection(t *testing.T) {

}

func TestMessageDiscriminator(t *testing.T) {
	message := []byte("StoRqQmQsx5srHm6cAdN3N7qgfPjRDy3ioumsRDWRZBYQv6ornU0.012695312")

	originalOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	messageDiscriminator(message)

	w.Close()
	os.Stdout = originalOutput
	var buf bytes.Buffer
	io.Copy(&buf, r)

	expected := "Recevied message :  StoRqQmQsx5srHm6cAdN3N7qgfPjRDy3ioumsRDWRZBYQv6ornU0.0126953125\nReceived storage request\nCID :  QmQsx5srHm6cAdN3N7qgfPjRDy3ioumsRDWRZBYQv6ornU\nFile Size :  0.0126953125"
	if buf.String() != expected {
		t.Errorf("Expected: %s, Got: %s", expected, buf.String())
	}

}

func TestHandleStorageRequest(t *testing.T) {

}
