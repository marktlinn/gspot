package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// defines the interface for the flags expected via the cli while using gspot.
type flags struct {
	url  string
	n, c int
}

// parses passed cli flags.
type parseFunc func(string) error

func (f *flags) urlFlag(p *string) parseFunc {
	return func(s string) error {
		_, err := url.Parse(s)
		*p = s
		return err
	}
}

func (f *flags) intFlag(p *int) parseFunc {
	return func(s string) (err error) {
		*p, err = strconv.Atoi(s)
		return err
	}
}

func (f *flags) parse() (err error) {
	parsedArgs := map[string]parseFunc{
		"url": f.urlFlag(&f.url),
		"n":   f.intFlag(&f.n),
		"c":   f.intFlag(&f.c),
	}
	for _, arg := range os.Args[1:] {
		n, v, ok := strings.Cut(arg, "=")
		if !ok {
			continue
		}
		parsed, ok := parsedArgs[strings.TrimPrefix(n, "-")]
		if !ok {
			continue
		}
		if err = parsed(v); err != nil {
			err = fmt.Errorf("inavlid value %q passed for flag %s, error: %w", v, n, err)
			break
		}
	}
	return err
}
