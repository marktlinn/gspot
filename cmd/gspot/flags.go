package main

import (
	"flag"
	"runtime"
)

// defines the interface for the flags expected via the cli while using gspot.
type flags struct {
	url  string
	n, c int
}

func (f *flags) parse() (err error) {
	// TODO: finish flad 5.
	url := flag.String("url", "", "Sets the url which will be targeted")
	n := flag.Int("n", 50, "Sets the number of requests that will be sent to the url in total")
	c := flag.Int("c", runtime.NumCPU(), "set the concurrency level i.e. how many requests will be sent concurrently")

	flag.Parse()
	f.url = *url
	f.n = *n
	f.c = *c
	return nil
}
