package main

import (
	"flag"
	"log"
	"os"

	"github.com/tankbusta/drivehash"
)

var (
	resultDirectory   string
	startingDirectory string
	processingBacklog int
	writerBacklog     int
	ignoreAdmin       bool
)

func init() {
	flag.StringVar(&resultDirectory, "result-dir", "./results", "the directory to save the results")
	flag.StringVar(&startingDirectory, "start-dir", "", "the directory or drive path to walk")
	flag.IntVar(&processingBacklog, "backlog-processing", drivehash.DefaultProcessingBacklog, "The maximum number of filepaths in the processing backlog at any given time")
	flag.IntVar(&writerBacklog, "backlog-writer", drivehash.DefaultWriterBacklog, "The maximum number of hashes in the writer backlog at any given time")
	flag.BoolVar(&ignoreAdmin, "ignore-admin", false, "Ignore the check if admin flag")

	flag.Parse()
}

func main() {
	if startingDirectory == "" {
		flag.Usage()
		os.Exit(1)
	}

	hasher := drivehash.New(processingBacklog, writerBacklog)

	if err := hasher.Start(startingDirectory, resultDirectory, ignoreAdmin); err != nil {
		log.Printf("[ X ] Failed to hash: %s\n", err)
	}
}
