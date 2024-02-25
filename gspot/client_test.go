package gspot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientDo(t *testing.T) {
	t.Parallel()

	const wantRes, WantErrs = 100, 0
	var (
		gotRes atomic.Int64

		server = newTestServer(t, func(_ http.ResponseWriter, _ *http.Request) {
			gotRes.Add(1)
		})

		req = newRequest(t, http.MethodGet, server.URL)
	)

	var c Client

	ttl := c.Do(context.Background(), req, wantRes)
	if got := gotRes.Load(); got != wantRes {
		t.Errorf("Res: %d; wanted %d", got, wantRes)
	}
	if got := ttl.Requests; got != wantRes {
		t.Errorf("Reqs: %d; wanted %d", got, wantRes)
	}
	if got := ttl.Errors; got != WantErrs {
		t.Errorf("Errs %d; wanted %d", got, WantErrs)
	}
}

func newRequest(t testing.TB, httpMethod, url string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(httpMethod, url, http.NoBody)
	if err != nil {
		t.Fatalf("Request err=%q; wanted nil", err)
	}
	return req
}

func newTestServer(t testing.TB, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(handler)
	t.Cleanup(s.Close)
	return s
}
