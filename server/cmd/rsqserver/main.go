package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

// Version
const VERSION = "0.0.1"

var (
	port         int
	database     string
	reset        bool
	workers      int
	versionOnly  bool
	pollInterval int // polling interval for workers in seconds
)

func init() {
	flag.IntVarP(&port, "port", "p", 3333, "port number to serve")
	flag.BoolVarP(&reset, "reset", "r", false, "wipe and reset the database")
	flag.StringVarP(&database, "database", "d", "models.db", "path and name of database to store model results")
	flag.IntVarP(&workers, "workers", "w", 0, "number of workers, set to negative number for no workers to be activated on initialization")
	flag.BoolVarP(&versionOnly, "version", "v", false, "print the version")
	flag.IntVarP(&pollInterval, "pollInterval", "i", 2, "polling interval for the workers to check for new queued models, default 2 seconds")
	flag.Parse()
}

func main() {
	if versionOnly {
		fmt.Print(VERSION)
		os.Exit(0)
	}

}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
