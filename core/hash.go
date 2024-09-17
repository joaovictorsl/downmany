package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Pair struct {
	sum  uint64
	path string
}

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, ch chan Pair) {
	data, err := readFile(filePath)
	if err != nil {
		panic(err)
	}

	_sum := uint64(0)
	for _, b := range data {
		_sum += uint64(b)
	}
	ch <- Pair{_sum, filePath}
}

// print the totalSum for all files and the files with equal sum
func Sum() (map[uint64]string, error) {
	ch := make(chan Pair)
	allSums := make(map[uint64]string)

	dataset := "./dataset"
	entries, err := os.ReadDir(dataset)
	if err != nil {
		log.Println("Directory './dataset' does not exist!")
		return nil, err
	}

	for _, path := range entries {
		fullPath := filepath.Join(dataset, path.Name())
		go sum(fullPath, ch)
	}

	for i := 0; i < 10; i++ {
		v := <-ch
		path := v.path
		index := v.sum
		allSums[index] = path
	}
	return allSums, nil
}
