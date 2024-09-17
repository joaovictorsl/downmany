package core

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Pair struct {
	sum  int64
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

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}
	ch <- Pair{int64(_sum), filePath}
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	ch := make(chan Pair)
	sums := make(map[int][]string)
	for _, path := range os.Args[1:] {
		go sum(path, ch)
	}

	for i := 0; i < 11; i++ {
		v := <-ch
		path := v.path
		index := int(v.sum)
		sums[index] = append(sums[index], path)
	}

	for sum, files := range sums {
		fmt.Printf("Sum %d: %v\n", sum, files)
	}
}
