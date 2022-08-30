// FFMPEG

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	child_process_manager "github.com/AgustinSRG/go-child-process-manager"
)

var (
	forward_lock *sync.Mutex
	forward_proc *os.Process
)

func initForward() {
	forward_lock = &sync.Mutex{}
	forward_proc = nil
}

func setProcess(p *os.Process) {
	forward_lock.Lock()
	defer forward_lock.Unlock()

	forward_proc = p
}

func killProcess() {
	forward_lock.Lock()
	defer forward_lock.Unlock()

	if forward_proc != nil {
		forward_proc.Kill()
		forward_proc = nil
	}

	os.Exit(0)
}

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

	child_process_manager.ConfigureCommand(cmd)

	err := cmd.Start()

	if err != nil {
		fmt.Println("Error: ffmpeg program failed: " + err.Error())
		os.Exit(1)
	}

	child_process_manager.AddChildProcess(cmd.Process)

	setProcess(cmd.Process)

	err = cmd.Wait()

	if err != nil {
		fmt.Println("Error: ffmpeg program failed: " + err.Error())
		os.Exit(1)
	}

	setProcess(nil)

	os.Exit(0)
}

func forwardCustom(customCommand string, debug bool) {
	args := strings.Fields(customCommand)

	var cmd *exec.Cmd

	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}

	if debug {
		cmd.Stderr = os.Stderr
		fmt.Println("Running command: " + cmd.String())
	}

	child_process_manager.ConfigureCommand(cmd)

	err := cmd.Start()

	if err != nil {
		fmt.Println("Error: ffmpeg program failed: " + err.Error())
		os.Exit(1)
	}

	child_process_manager.AddChildProcess(cmd.Process)

	setProcess(cmd.Process)

	err = cmd.Wait()

	if err != nil {
		fmt.Println("Error: ffmpeg program failed: " + err.Error())
		os.Exit(1)
	}

	setProcess(nil)

	os.Exit(0)
}
