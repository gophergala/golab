// Package main of GoLab is the main package defining the entry point
// and which compiles into the GoLab executable.
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala/golab/ctrl"
	"log"
	"net/http"
	"runtime"
)

var port int

// processFlags registers flags, parses them and validates them.
// Returns nil if everything is ok, else an error.
func processFlags() error {
	flag.IntVar(&port, "port", 60148, "Port to start the UI web server on; valid range: 0..65535")

	flag.Parse()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is outside of range 0..65535", port)
	}

	return nil
}

// main is the entry point of GoLab.
// Processes the command line flags, initializes the game engine
// and starts the UI webserver.
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := processFlags(); err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	ctrl.InitNew()

	fmt.Printf("Starting GoLab webserver on port %d...\n", port)
	fmt.Printf("Open http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
