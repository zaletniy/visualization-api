package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"

	"visualization-api/pkg/logging"
)

var (
	logRotate  *log.RotateWriter
	version    = "UNDEFINED"
	gitVersion = "UNDEFINED"

	//app level flags
	versionParam = flag.Bool("version", false, "Prints version information")
)

func cleanupOnExit() {
	// this function is used to perform all cleanup on application exit
	// such as file descriptor close
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigc
		log.Logger.Info("Caught signal '", s, "' shutting down")
		// close global descriptor
		logRotate.Lock.Lock()
		defer logRotate.Lock.Unlock()
		err := logRotate.Fp.Close()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}()
}

func main() {

	flag.Parse()

	if *versionParam {
		fmt.Printf("auth_proxy version %s %s \n", version, gitVersion)
		os.Exit(0)
	}

	cleanupOnExit()
}
