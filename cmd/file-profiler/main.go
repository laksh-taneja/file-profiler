package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/laksh-taneja/file-profiler/internal/processor"
	"github.com/laksh-taneja/file-profiler/internal/reporter"
	"github.com/laksh-taneja/file-profiler/internal/utils"
)

func main() {
	lines := flag.Bool("l", false, "count lines in file(s)")
	longLine := flag.Bool("long", false, "increase buffer size")
	words := flag.Bool("w", true, "count words in file(s)")
	hash := flag.Bool("hash", true, "calculate hash of the file(s)")
	worker := flag.Int("workers", 3, "no. of concurrent workers")
	output := flag.String("o", "./default-results.json", "select the output file")
	flag.Parse()

	n := len(flag.Args())

	if n < 1 {
		log.Fatalf("Incorrect usage: fprofiler -<flags> [file]...")
	}

	if exists, err := utils.FileExists(*output); exists {
		if err != nil {
			log.Fatalf("Error accessing %v: %v", *output, err)
		}
		fmt.Printf("File %s already exists. Overwrite? (y/n): ", *output)
		var cont string
		for {
			fmt.Scan(&cont)
			lowerCont := strings.ToLower(cont)
			if lowerCont == "y" {
				break
			}
			if lowerCont == "n" {
				fmt.Println("Exiting...")
				os.Exit(0)
			}
			fmt.Print("Please enter 'y' or 'n': ")
		}
	}

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer outputFile.Close()

	jobs, results := processor.StartPool(*worker)

	log.Printf("Processing %v file(s)", n)
	go func() {
		for _, f := range flag.Args() {
			inpfd := processor.FileDimension{
				Filename:   f,
				CountLines: *lines,
				LongLines:  *longLine,
				CountWords: *words,
				DoHash:     *hash,
			}
			jobs <- inpfd
		}
		close(jobs)
	}()

	reporter.StreamToJSON(results, outputFile, n)
}
