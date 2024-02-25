package gspot

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

/*
Client encapsulates:

* C - the Concurency level

* RSP - the throttle limit for requests per second

* Timeout - the maximum time each request should be allowed before returning an error. If not set or set to 0, no Timeout is applied.
*/
type Client struct {
	RPS     int
	C       int
	Timeout time.Duration
}

// When created and passed to the Do functions, Option can be used for setting and changing the Client's default behvaiour.
type Option func(*Client)

// The Timeout option changes the Client's RPS field.
func Timeout(d time.Duration) Option {
	return func(c *Client) { c.Timeout = d }
}

// The Concurrency option changes the Client's concurrency level.
func Concurrency(n int) Option {
	return func(c *Client) { c.C = n }
}

// a helper function that sends n requests.
// Idle connections are closed once the performance data has been registered.
func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})

	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.concurrencyDefault()))
	}

	var (
		ttl    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	for res := range split(p, c.concurrencyDefault(), c.send(client)) {
		ttl.Merge(res)
	}
	return &ttl
}

/*
Exposeses the Client.Do method via a function that can be easily called. Sends n Get requests to the specified url with a concurrency level defaulting to the number of CPUs on the host machine.

Do can take Options (a variadic func) for any values to be set in the Client's fields.
*/
func Do(ctx context.Context, url string, n int, options ...Option) (*Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request error: %w", err)
	}
	var c Client
	for _, option := range options {
		option(&c)
	}
	return c.Do(ctx, r, n), nil
}

// Sends n HTTPS requests and returns the aggregated result once all have completed.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	if c == nil {
		panic("error in gspot.Do(): Client cannot be nil.")
	}
	t := time.Now()
	ttl := c.do(ctx, r, n)
	return ttl.Finalise(time.Since(t))
}

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.concurrencyDefault(),
		},
	}
}

// Sets the concurrency level default to the number of CPU on the host machine if level not explicitly set.
func (c *Client) concurrencyDefault() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}
