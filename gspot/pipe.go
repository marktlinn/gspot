package gspot

import "net/http"

// generates requests from the passed function n times, sending the results to the specified (send-only) channel.
func Produce(out chan<- *http.Request, n int, fn func() *http.Request) {
	for ; n > 0; n-- {
		out <- fn()
	}
}

// runs Produce inside a goroutine enabling concurrency.
func produce(n int, fn func() *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Produce(out, n, fn)
	}()
	return out
}
