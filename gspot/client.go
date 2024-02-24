package gspot

import (
	"context"
	"net/http"
	"time"
)

// Cleint encapsulates the Concurency level (C)
// and the throttle limit for requests per second (RPS)
type Client struct {
	RPS int
	C   int
}

// a helper function that sends n requests.
func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})

	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.C))
	}

	var ttl Result
	for ; n > 0; n-- {
		select {
		case <-ctx.Done():
			return &ttl
		default:
			ttl.Merge(Send(r))
		}
	}
	return &ttl
}

// Sends n HTTPS requests and retuns the aggregated result once all have completed.
func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	ttl := c.do(ctx, r, n)
	return ttl.Finalise(time.Since(t))
}
