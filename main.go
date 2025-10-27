package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// APIRequest models an API call & the result
//
type APIRequest struct {
	ID        int
	Endpoint  string
	Timestamp time.Time
	Response  string
	Duration  time.Duration
	Timeout   bool
}

func randomEndpoint() string {
	endpoints := []string{"/api/foo", "/api/bar", "/api/baz", "/api/qux"}
	return endpoints[rand.Intn(len(endpoints))]
}

func apiCall(req *APIRequest) {
	// Simulate an API request
	// Sleep random 400ms-1200ms
	sleepTime := time.Duration(400+rand.Intn(801)) * time.Millisecond
	start := time.Now()
	// Artificially sleep
	time.Sleep(sleepTime)
	req.Response = "OK"
	req.Duration = time.Since(start)
	// For this mock, always response within duration, but req.Duration records real time
}

func processRequest(req *APIRequest, sem chan struct{}, results chan<- APIRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	// Acquire semaphore slot
	sem <- struct{}{}
	timeout := time.After(1500 * time.Millisecond)
	start := time.Now()
	done := make(chan struct{})

	go func() {
		apiCall(req)
		close(done)
	}()

	select {
	case <-done:
		// Completed on time
		elapsed := time.Since(start)
		if elapsed > 1500*time.Millisecond {
			req.Timeout = true
		}
	case <-timeout:
		// Timeout occurred
		req.Timeout = true
		// Wait for goroutine to finish
		<-done
	}

	// Release semaphore slot
	<-sem
	// Send to results channel
	results <- *req
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const reqCount = 10
	const maxConcurrency = 3
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	results := make(chan APIRequest, reqCount)

	for i := 0; i < reqCount; i++ {
		req := &APIRequest{
			ID:        i + 1,
			Endpoint:  randomEndpoint(),
			Timestamp: time.Now(),
		}
		wg.Add(1)
		go processRequest(req, sem, results, &wg)
	}

	// Wait for all to finish
	wg.Wait()
	close(results)

	fmt.Println("--- API Request Summary ---")
	for res := range results {
		out := fmt.Sprintf("[ID %02d] %-10s  Elapsed: %-7v", res.ID, res.Endpoint, res.Duration)
		if res.Timeout || res.Duration > 1500*time.Millisecond {
			out += "  [timeout]"
		}
		fmt.Println(out)
	}
}