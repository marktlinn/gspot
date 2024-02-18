package gspot

import (
	"fmt"
	"net/http"
	"time"
)

// Represents the function called by the Client.Do method, which sends an HTTP request and returns the perf results.
type SendFunc func(*http.Request) *Result

// Sends an HTTP request and returns the the perf result.
func Send(r *http.Request) *Result {
	t := time.Now()

	fmt.Printf("req: %s\n", r.URL)
	time.Sleep(100 * time.Millisecond)

	return &Result{
		Duration: time.Since(t),
		Bytes:    100,
		Status:   http.StatusOK,
	}
}
