// FFMPEG

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func forwardToRTMP(ffmpegBin string, source string, rtmpURL string, debug bool) {
	args := make([]string, 1)

	args[0] = ffmpegBin

	args = append(args, "-re")

	args = append(args, "-protocol_whitelist", "file,sdp,udp,rtp")

	// INPUT
	args = append(args, "-f", "sdp", "-i", source)

	// DESTINATION
	args = append(args, "-f", "flv", rtmpURL)

	cmd := exec.Command(ffmpegBin)
	cmd.Args = args

	if debug {
		cmd.Stderr = os.Stderr
		fmt.Println("Running command: " + cmd.String())
	}

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error: ffmpeg program failed: " + err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
