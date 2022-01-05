package main

import (
	"flag"

	"github.com/adrianolmedo/go-restapi-practice/internal/app"
)

func main() {
	// Configuration.
	c := flag.String("c", "config.json", "Load configuration from `FILE`")
	flag.Parse()

	app.Run(*c)
}
