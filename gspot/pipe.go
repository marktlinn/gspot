package gspot

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// generates requests from the passed function n times, sending the results to the specified (send-only) channel.
func Produce(ctx context.Context, out chan<- *http.Request, n int, fn func() *http.Request) {
	for ; n > 0; n-- {
		select {
		case <-ctx.Done():
			return
		case out <- fn():
		}
	}
}

// runs Produce inside a goroutine.
func produce(ctx context.Context, n int, fn func() *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Produce(ctx, out, n, fn)
	}()
	return out
}

// sets a delay on messages between channels i.e. between what comes in and what goes out.
func Throttle(in <-chan *http.Request, out chan<- *http.Request, delay time.Duration) {
	t := time.NewTicker(delay)
	defer t.Stop()

	for r := range in {
		<-t.C
		out <- r
	}
}

// runs Throttle inside a goroutine.
func throttle(in <-chan *http.Request, delay time.Duration) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Throttle(in, out, delay)
	}()
	return out
}

// Splits the pipline up, running the fn func with the data received from in, the results being channeled to out.
func Split(in <-chan *http.Request, out chan<- *Result, c int, fn SendFunc) {
	send := func() {
		for r := range in {
			out <- fn(r)
		}
	}

	var wg sync.WaitGroup
	wg.Add(c)
	for ; c > 0; c-- {
		go func() {
			defer wg.Done()
			send()
		}()
	}
	wg.Wait()
}

// runs the Split function.
func split(in <-chan *http.Request, c int, fn SendFunc) <-chan *Result {
	out := make(chan *Result)
	go func() {
		defer close(out)
		Split(in, out, c, fn)
	}()
	return out
}
