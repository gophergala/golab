// Package main of GoLab is the main package defining the entry point
// and which compiles into the GoLab executable.
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala/golab/ctrl"
	"github.com/gophergala/golab/model"
	"github.com/gophergala/golab/view"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// port tells on which port to open the UI web server
var port int

// autoOpen tells if the UI web page should be auto-opened in the users's default browser
var autoOpen bool

// processFlags registers flags, parses them and validates them.
// Returns nil if everything is ok, else an error.
func processFlags() error {
	// General flags
	flag.IntVar(&port, "port", 1234, "Port to start the UI web server on; valid range: 0..65535")
	flag.BoolVar(&autoOpen, "autoOpen", true, "Auto-opens the UI web page in the default browser")

	// Control/Engine flags
	flag.IntVar(&ctrl.LoopDelay, "loopDelay", 50, "Loop delay of the game engine, in milliseconds; valid range: 10..100")
	flag.Float64Var(&model.V, "v", model.BlockSize*2.0, "moving speed of Gopher and the Buddlogs in pixel/sec; valid range: 20..200")

	// View flags
	flag.IntVar(&view.ViewWidth, "viewWidth", 700, "Width of the view image in pixels in the UI web page; valid range: 200..1000")
	flag.IntVar(&view.ViewHeight, "viewHeight", 700, "Height of the view image in pixels in the UI web page; valid range: 200..1000")

	flag.Parse()

	if port < 0 || port > 65535 {
		return fmt.Errorf("port %d is outside of range 0..65535", port)
	}

	if ctrl.LoopDelay < 10 || ctrl.LoopDelay > 100 {
		return fmt.Errorf("loopDelay %d is outside of range 10..100", ctrl.LoopDelay)
	}

	if model.V < 20 || model.V > 200 {
		return fmt.Errorf("v %f is outside of range 20..200", model.V)
	}

	if view.ViewWidth < 200 || view.ViewWidth > 1000 {
		return fmt.Errorf("viewWidth %d is outside of range 200..1000", view.ViewWidth)
	}

	if view.ViewHeight < 200 || view.ViewHeight > 1000 {
		return fmt.Errorf("viewHeight %d is outside of range 200..1000", view.ViewHeight)
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

	ctrl.StartEngine()

	fmt.Printf("Starting GoLab webserver on port %d...\n", port)
	url := fmt.Sprintf("http://localhost:%d/", port)
	if autoOpen {
		fmt.Printf("Opening %s...\n", url)
		if err := open(url); err != nil {
			fmt.Println("Auto-open failed:", err)
			fmt.Printf("Open %s in your browser.\n", url)
		}
	} else {
		fmt.Printf("Auto-open not enabled, open %s in your browser.\n", url)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}
