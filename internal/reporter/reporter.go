package reporter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/laksh-taneja/file-profiler/internal/processor"
)

func StreamToJSON(results chan *processor.FileResult, outputFile io.Writer, total int) {
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	writer.WriteString("[\n")
	first := true
	count := 0
	for res := range results {
		count++
		percentage := (float64(count) / float64(total)) * 100
		// \r moves the cursor back to the start of the line
		fmt.Fprintf(os.Stderr, "\rProgress: [%-20s] %.1f%% (%d/%d)",
			strings.Repeat("=", (count*20)/total),
			percentage, count, total)
		if !first {
			writer.WriteString(",\n")
		}
		b, err := json.MarshalIndent(res, "	", "	")
		if err != nil {
			log.Printf("failed to marshal: %v", err)
			continue
		}
		writer.Write(b)
		first = false
	}
	writer.WriteString("\n]")
}
