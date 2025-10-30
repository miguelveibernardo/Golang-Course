package main

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	const totalRequest = 500
	const concurrency = 20

	var wg sync.WaitGroup
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second}

	sem := make(chan struct{}, concurrency)
	for i := 0; i < totalRequest; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(n int) {
			defer wg.Done()
			defer func() { <-sem }()

			data := url.Values{}
			data.Set("description", fmt.Sprintf("load task %d", n))
			resp, err := client.PostForm("http://localhost:8080/create", data)
			if err != nil {
				fmt.Println("Request error", err)
				return
			}
			resp.Body.Close()
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d requests in %v\n", totalRequest, duration)
	fmt.Printf("Average per request: %v\n", duration/time.Duration(totalRequest))
}
