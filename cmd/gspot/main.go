package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/marktlinn/gspot/gspot"
)

const (
	bannerText = `
   ___________ ____  ____  ______
  / ____/ ___// __ \/ __ \/_  __/
 / / __ \__ \/ /_/ / / / / / /   
/ /_/ /___/ / ____/ /_/ / / /    
\____//____/_/    \____/ /_/     
																
`
)

func getBannerText() string { return bannerText[1:] }

func main() {
	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		os.Exit(1)
	}
}

func run(flg *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		// defaults if no flags present.
		n: 100,
		c: runtime.NumCPU(),
	}
	if err := f.parse(flg, args); err != nil {
		return err
	}
	fmt.Fprintln(out, getBannerText())
	fmt.Fprintf(out, "Making %d requests to %s with concurrency set to %d.\n", f.n, f.url, f.c)

	var ttl gspot.Result
	ttl.Merge(&gspot.Result{
		Bytes:    1024,
		Status:   http.StatusOK,
		Duration: time.Second,
	})
	ttl.Merge(&gspot.Result{
		Bytes:    1024,
		Status:   http.StatusOK,
		Duration: time.Second,
	})
	ttl.Merge(&gspot.Result{
		Status:   http.StatusConflict,
		Duration: time.Second,
	})
	ttl.Merge(&gspot.Result{
		Bytes:    556,
		Status:   http.StatusOK,
		Duration: 2 * time.Second,
	})
	ttl.Finalise(2 * time.Second)
	ttl.Fprint(out)

	return nil
}
