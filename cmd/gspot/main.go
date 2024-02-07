package main

import "fmt"

const (
	bannerText = `
   ___________ ____  ____  ______
  / ____/ ___// __ \/ __ \/_  __/
 / / __ \__ \/ /_/ / / / / / /   
/ /_/ /___/ / ____/ /_/ / / /    
\____//____/_/    \____/ /_/     
																
`

	usageText = `
Usage:
	-url
		URL of the HTTP server you want to make requests against 
		(Required value)
	-n
		Number of requests you want to make
	-c
		The number of requests to be executed concurrently
	`
)

func getBannerText() string { return bannerText[1:] }
func getUsageText() string  { return usageText[1:] }

func main() {
	fmt.Println(getBannerText())
	fmt.Println(getUsageText())
}
