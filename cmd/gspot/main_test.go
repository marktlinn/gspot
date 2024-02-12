package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

type testEnv struct {
	args           string
	stdout, stderr bytes.Buffer
}

func (e *testEnv) run() error {
	flg := flag.NewFlagSet("gspot", flag.ContinueOnError)
	flg.SetOutput(&e.stderr)
	return run(flg, strings.Fields(e.args), &e.stdout)
}

func TestRun(t *testing.T) {
	t.Parallel()
	happyPath := map[string]struct{ in, out string }{
		"url_as_solo_argument": {
			"http://test.com",
			fmt.Sprintf("Making 100 requests to http://test.com with concurrency set to %d.\n", runtime.NumCPU()),
		},
		"url_https": {
			"https://test.com",
			fmt.Sprintf("Making 100 requests to https://test.com with concurrency set to %d.\n", runtime.NumCPU()),
		},
		"n_c_given": {
			"-n=1000 -c=10  http://test.com",
			"Making 1000 requests to http://test.com with concurrency set to 10.\n",
		},
	}
	unhappyPath := map[string]string{
		"no_url":           "",
		"c_greater_than_n": "-c=100 -n=10 https://example.com",
		"incorrect_scheme": "wss://example.com",
		"no_host":          "https://",
		"incorrect_n":      "-n=hello https://example.com",
		"incorrect_c":      "-c=hello https://example.com",
		"negative_nums":    "-c=-10 -n=-5 http://example.com",
		"negative_c":       "-c=-10 -n=5 http://example.com",
		"negative_":        "-c=10 -n=-5 http://example.com",
	}
	for name, tt := range happyPath {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tt.in}
			if err := e.run(); err != nil {
				t.Fatalf("got %q;\nwanted nil", err)
			}
			if out := e.stdout.String(); !strings.Contains(out, tt.out) {
				t.Errorf("got %s;\nwanted %q", out, tt.out)
			}
		})
	}
	for name, tt := range unhappyPath {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tt}
			if e.run() == nil {
				t.Fatalf("got nil; wanted err")
			}
			if e.stderr.Len() == 0 {
				t.Fatal("stderr = 0 bytes;wanted >0")
			}
		})
	}
}
