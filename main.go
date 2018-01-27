package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var shouldProfile = flag.Bool("profile", false, "Enable CPU profiling")
var profilePath = flag.String("profilePath", "/tmp/chessProfile.txt", "Location of CPU profile")

func main() {
	flag.Parse()

	// Enable profiling on command line flag.
	if *shouldProfile {
		startProfile(*profilePath)
		defer stopProfile()
	}

	startEngine()
}

func startProfile(filename string) {
	profile, err := os.Create(filename)

	if err != nil {
		log.Fatal("Could not create profile: ", err)
	}
	err = pprof.StartCPUProfile(profile)

	if err != nil {
		log.Fatal("Could not start profile: ", err)
	}
}

func stopProfile() {
	pprof.StopCPUProfile()
}
