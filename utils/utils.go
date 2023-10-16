package utils

import (
	"fmt"
	"os"
)

func ErrorHandler(err error) {
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		panic(0)
	}
}

func ListPrint(list []string) {
	for _, element := range list {
		fmt.Print(element + " ")
	}
}

func GetFileSize(path string) float64 {
	// Returns file size in KB

	fileInfo, err := os.Stat(path)
	ErrorHandler(err)
	fileSize := fileInfo.Size()

	return float64(fileSize) / 1024.0
}
