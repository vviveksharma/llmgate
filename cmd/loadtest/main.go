package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	addr := flag.String("addr", "http://localhost:8080", "gateway base URL")
	n := flag.Int("n", 100, "number of requests")
	c := flag.Int("c", 10, "concurrency")
	flag.Parse()

	client := &http.Client{Timeout: 30 * time.Second}
	body := []byte(`{"model":"mock","messages":[{"role":"user","content":"hi"}]}`)

	var wg sync.WaitGroup
	sem := make(chan struct{}, *c)
	var ok, fail int64
	var mu sync.Mutex

	for i := 0; i < *n; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			req, _ := http.NewRequest(http.MethodPost, *addr+"/v1/chat/completions", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-API-Key", "dev-key")
			resp, err := client.Do(req)
			mu.Lock()
			defer mu.Unlock()
			if err != nil || resp.StatusCode >= 400 {
				fail++
				return
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			ok++
		}()
	}
	wg.Wait()
	fmt.Printf("loadtest complete: ok=%d fail=%d\n", ok, fail)
}
