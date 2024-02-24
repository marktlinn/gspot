package gspot

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Represents the result of a gspot request.
type Result struct {
	Requests int           // the number of requests made.
	Errors   int           // the number of errors (if any) from the reuqest.
	Error    error         // nil if no error raised, else error.
	Status   int           // HTTP status code returned from request.
	Bytes    int64         // the number of bytes downloaded by the request
	RPS      float64       // the number of requests per second.
	Duration time.Duration // the time taken to complete the requests.
	Fastest  time.Duration // the time taken by the fasted request to complete.
	Slowest  time.Duration // time taken for the slowest request to complete.
}

// Merges the Result on which the method is called with the Result passed to the method.
func (r *Result) Merge(o *Result) {
	r.Requests++
	r.Bytes += o.Bytes

	if r.Fastest == 0 || o.Duration < r.Fastest {
		r.Fastest = o.Duration
	}

	if o.Duration > r.Slowest {
		r.Slowest = o.Duration
	}

	switch {
	case o.Error != nil:
		fallthrough
	case o.Status >= http.StatusBadRequest:
		r.Errors++
	}
}

// Calculates the total duration based on all completed requests
func (r *Result) Finalise(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

// Formats and prints the associated results to an io.Writer.
func (r *Result) Fprint(out io.Writer) {
	data := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	data("\nSummary:\n")
	data("\tSuccess: %.0f%%\n", r.success())
	data("\tRequests: %d\n", r.Requests)
	data("\tRPS: %.1f\n", r.RPS)
	data("\tBytes: %d\n", r.Bytes)
	data("\tDuration: %s\n", r.Duration)
	data("\tErrors: %d\n", r.Errors)
	if r.Requests > 1 {
		data("\tFastest : %s\n", round(r.Fastest))
		data("\tSlowest : %s\n", round(r.Slowest))

	}
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}

func (r *Result) success() float64 {
	res, e := float64(r.Requests), float64(r.Errors)
	return (res - e) / res * 100
}

// Returns the result as a string.
func (r *Result) String() string {
	var s strings.Builder
	r.Fprint(&s)
	return s.String()
}
