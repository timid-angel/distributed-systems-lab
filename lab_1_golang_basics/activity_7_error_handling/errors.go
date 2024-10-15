package main

import (
	"fmt"
	"os"
)

// function to open a file, returns the error if there is any
func readFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	fmt.Println("File opened successfully:", fileName)
	return nil
}

func main() {
	err := readFile("test.txt")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File read successfully.")
	}
}
