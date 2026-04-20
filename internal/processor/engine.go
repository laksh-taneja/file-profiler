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

func Analyze(f FileDimension) (*FileResult, error) {
	log.Printf("Processing file %v", f.Filename)

	file, err := os.Open(f.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words int
	var lines int
	var checksum string

	hasher := sha256.New()
	var reader io.Reader = file

	if f.DoHash {
		reader = io.TeeReader(file, hasher)
	}
	if f.CountLines || f.CountWords {
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
		if err = scanner.Err(); err != nil {
			return nil, err
		}
	} else if f.DoHash {
		io.Copy(io.Discard, reader)
	}
	if f.DoHash {
		checksum = hex.EncodeToString(hasher.Sum(nil))
	}

	return &FileResult{WordCount: words, LineCount: lines, CheckSum: checksum}, nil
}
