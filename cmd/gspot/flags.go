package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// defines the interface for the flags expected via the cli while using gspot.
type flags struct {
	url  string
	n, c int
}

// defines a numeric interface for positive numbers.
type num int

// converts a pointer to an int into a pointer to a num.
func toNum(p *int) *num {
	return (*num)(p)
}

func (n *num) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("num should be a positive int")
	}
	*n = num(v)
	return err
}

func (n *num) String() string {
	return strconv.Itoa(int(*n))
}

func (f *flags) parse() (err error) {
	flag.StringVar(&f.url, "url", "", "Sets the url which will be targeted")
	flag.IntVar(&f.n, "n", 50, "Sets the number of requests that will be sent to the url in total")
	flag.IntVar(&f.c, "c", runtime.NumCPU(), "set the concurrency level i.e. how many requests will be sent concurrently")

	flag.Parse()
	if err := f.validateFlag(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		return err
	}
	return nil
}

// validates the flags passed to the cli.
func (f *flags) validateFlag() error {
	if err := validateURL(f.url); err != nil {
		return fmt.Errorf("invalid -url: %q", err)
	}
	if f.c > f.n {
		return fmt.Errorf("-c=%d should be less than or equal to -n=%d", f.c, f.n)
	}
	return nil
}

func validateURL(s string) error {
	url, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("-url is required")
	case validateScheme(&url.Scheme):
		err = errors.New("scheme must be `http://` or `https://`")
	case url.Host == "":
		err = errors.New("host is missing")
	}
	return err
}

// validates the url scheme is either http or https.
func validateScheme(s *string) bool {
	return *s != "http" && *s != "https"
}
