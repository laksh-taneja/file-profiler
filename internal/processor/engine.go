package processor

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/laksh-taneja/file-profiler/internal/utils"
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
	_, err := utils.FileExists(fd.Filename)
	if err != nil {
		return nil, err
	}
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
		_, err = io.Copy(io.Discard, reader)
	}
	if fd.DoHash {
		checksum = hex.EncodeToString(hasher.Sum(nil))
	}

	return &FileResult{Filename: fd.Filename, WordCount: words, LineCount: lines, CheckSum: checksum}, err
}
