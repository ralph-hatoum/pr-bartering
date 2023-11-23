package utils

/*
General purpose, useful functions
*/

import (
	"fmt"
	"os"
)

func ErrorHandler(err error) {

	/*
		To handle errors ; panics with -1 if there is an error
		Arguments : error of type error
	*/

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		panic(-1)
	}
}

func ListPrint(list []string) {

	/*
		To print a string list's elements
		Arguments : list as a string list
	*/

	for _, element := range list {
		fmt.Print(element + " ")
	}
}

func GetFileSize(path string) float64 {

	/*
		Returns file size in KB
		Arguments : path to file as a string
		Returns : file size in KB as float64
	*/
	fileInfo, err := os.Stat(path)
	ErrorHandler(err)
	fileSize := fileInfo.Size()

	return float64(fileSize) / 1024.0
}
