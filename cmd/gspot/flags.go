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
	flag.StringVar(&f.url, "url", "", "Sets the url which will be targeted")
	flag.IntVar(&f.n, "n", 50, "Sets the number of requests that will be sent to the url in total")
	flag.IntVar(&f.c, "c", runtime.NumCPU(), "set the concurrency level i.e. how many requests will be sent concurrently")

	flag.Parse()
	return nil
}
