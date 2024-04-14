package configextractor

import (
	"testing"
	"fmt"
)

func TestConfigExtractor(t *testing.T) {
	config := ConfigExtractor("../config.yaml")
	fmt.Println("Port :", config.Port)
	port := fmt.Sprint(config.Port)
	fmt.Println("Port as string :",port)
	ConfigPrinter(config)
}