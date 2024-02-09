package main

import (
	"fmt"
	"os"
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
	f := &flags{}
	if err := f.parse(); err != nil {
		os.Exit(1)
	}
	fmt.Println(getBannerText())
	fmt.Printf("Making %d requests to %s with concurrency set to %d.\n", f.n, f.url, f.c)

}
