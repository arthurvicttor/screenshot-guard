package main

import (
	"flag"
	"log"
	"os"
)

const version = "0.1"

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ltime)

	logger.Printf("Screenshot Guard v%s", version)

	if *verbose {
		logger.Println("Verbose mode enabled")
	}

	logger.Println("Guard started")

}