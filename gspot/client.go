package gspot

import (
	"net/http"
	"time"
)

type Client struct {
	// TODO
}

// a helper function that sends n requests.
func (c *Client) do(r *http.Request, n int) *Result {
	var ttl Result
	for ; n > 0; n-- {
		ttl.Merge(Send(r))
	}
	return &ttl
}

// Sends n HTTPS requests and retuns the aggregated result once all have completed.
func (c *Client) Do(r *http.Request, n int) *Result {
	// TODO: add concurrency
	t := time.Now()
	ttl := c.do(r, n)
	return ttl.Finalise(time.Since(t))
}
