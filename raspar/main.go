package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoWild/raspar/crawler"
)

func main() {
	//flags
	url := flag.String("url", "", "the URL/website to be searched")
	flag.Parse()

	if *url == "" {
		fmt.Println("URL string is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// web
	data := crawler.GetDataAndResponse(*url)
	fmt.Println(data)
}
