package core

import (
	"fmt"
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
	data, err := os.ReadFile(filePath)
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
func Sum(dataset string) (map[uint64]string, error) {
	ch := make(chan Pair)
	allSums := make(map[uint64]string)

	entries, err := os.ReadDir(dataset)
	if err != nil {
		log.Printf("Directory '%s' does not exist!\n", dataset)
		return nil, err
	}

	for _, path := range entries {
		fullPath := filepath.Join(dataset, path.Name())
		go sum(fullPath, ch)
	}

	for i := 0; i < len(entries); i++ {
		v := <-ch
		path := v.path
		index := v.sum
		allSums[index] = path
	}

	return allSums, nil
}
