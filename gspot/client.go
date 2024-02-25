package gspot

import (
	"context"
	"net/http"
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

// a helper function that sends n requests.
// Idle connections are closed once the performance data has been registered.
func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})

	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
	}

	var (
		ttl    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	for res := range split(p, c.C, c.send(client)) {
		ttl.Merge(res)
	}
	return &ttl
}

// Sends n HTTPS requests and returns the aggregated result once all have completed.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
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
			MaxIdleConnsPerHost: c.C,
		},
	}
}
