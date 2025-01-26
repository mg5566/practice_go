package main

import (
	"runtime"

	"example.com/golang_test/web_crawler/crawler"
)

func main() {
	numCPUS := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUS)

	crawler.Crawler("https://github.com/mg5566?tab=following")
}
