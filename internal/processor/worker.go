package processor

import (
	"log"
	"sync"
)

func StartPool(workerCount int) (chan FileDimension, chan *FileResult) {
	jobs := make(chan FileDimension)
	results := make(chan *FileResult)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Go(func() {
			for job := range jobs {
				res, err := Analyze(job)
				if err != nil {
					log.Printf("\nError %v: %v\n", job.Filename, err)
					continue
				}
				results <- res
			}
		})
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	return jobs, results
}
