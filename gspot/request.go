package gspot

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Represents the function called by the Client.Do method, which sends an HTTP request and returns the perf results.
type SendFunc func(*http.Request) *Result

// Sends an HTTP request and returns the the perf result.
func Send(r *http.Request) *Result {
	t := time.Now()

	var (
		statusCode int
		bytes      int64
	)

	res, err := http.DefaultClient.Do(r)
	if err == nil {
		statusCode = res.StatusCode
		bytes, err = io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}

	fmt.Printf("req: %s\n", r.URL)

	return &Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   statusCode,
		Error:    err,
	}
}
