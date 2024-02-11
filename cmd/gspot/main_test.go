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
		"n_c_given": {
			"-n=1000 -c=10  http://test.com",
			"Making 1000 requests to http://test.com with concurrency set to 10.\n",
		},
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
}
