package processor

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func processStream(reader io.Reader, f FileDimension) (int, int, error) {
	var lines, words int
	scanner := bufio.NewScanner(reader)
	if f.LongLines {
		const maxCapacity = 10 * 1024 * 1024
		buf := make([]byte, 1024*1024)
		scanner.Buffer(buf, maxCapacity)
	}
	for scanner.Scan() {
		if f.CountLines {
			lines++
		}
		if f.CountWords {
			lineText := scanner.Text()
			words += len(strings.Fields(lineText))
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	return lines, words, nil
}

func Analyze(fd FileDimension) (*FileResult, error) {
	log.Printf("Processing file %v", fd.Filename)

	file, err := os.Open(fd.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words int
	var lines int
	var checksum string

	hasher := sha256.New()
	var reader io.Reader = file

	if fd.DoHash {
		reader = io.TeeReader(file, hasher)
	}
	if fd.CountLines || fd.CountWords {
		lines, words, err = processStream(reader, fd)
	} else if fd.DoHash {
		io.Copy(io.Discard, reader)
	}
	if fd.DoHash {
		checksum = hex.EncodeToString(hasher.Sum(nil))
	}

	return &FileResult{WordCount: words, LineCount: lines, CheckSum: checksum}, nil
}
