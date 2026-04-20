package processor

import (
	"bufio"
	"log"
	"os"
)

func Analyze(filename string, countLines bool) (int, error) {
	log.Printf("Processing file %v", filename)

	file, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	var count int

	scanner := bufio.NewScanner(file)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	for scanner.Scan() {
		count++
	}

	return count, nil
}
