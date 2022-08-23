// Main

package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// Program entry point
func main() {
	// Read env vars
	ffmpegPath := os.Getenv("FFMPEG_PATH")

	if ffmpegPath == "" {
		ffmpegPath = "/usr/bin/ffmpeg"
	}

	// Read arguments
	args := os.Args

	if len(args) < 2 {
		if len(args) > 1 && (args[1] == "--help" || args[1] == "-h") {
			printHelp()
		} else if len(args) > 1 && (args[1] == "--version" || args[1] == "-v") {
			printVersion()
		} else {
			printHelp()
		}
		return
	}

	debug := false
	authToken := ""
	authSecret := ""
	portAudio := 0
	portVideo := 0
	sdpFile := ""
	forwardMode := ""
	forwardParam := ""

	source := ""

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "--debug" {
			debug = true
		} else if arg == "--ffmpeg-path" {
			if i == len(args)-3 {
				fmt.Println("The option '--ffmpeg-path' requires a value")
				os.Exit(1)
			}
			ffmpegPath = args[i+1]
			i++
		} else if arg == "--input" || arg == "-i" {
			if i == len(args)-3 {
				fmt.Println("The option '--input' requires a value")
				os.Exit(1)
			}
			source = args[i+1]
			i++
		} else if arg == "--sdp-file" || arg == "-sdp" {
			if i == len(args)-3 {
				fmt.Println("The option '--sdp-file' requires a value")
				os.Exit(1)
			}
			sdpFile = args[i+1]
			i++
		} else if arg == "--forward-mode" || arg == "-fm" {
			if i == len(args)-3 {
				fmt.Println("The option '--forward-mode' requires a value")
				os.Exit(1)
			}
			forwardMode = args[i+1]
			i++
		} else if arg == "--auth" || arg == "-a" {
			if i == len(args)-3 {
				fmt.Println("The option '--auth' requires a value")
				os.Exit(1)
			}
			authToken = args[i+1]
			i++
		} else if arg == "--video-port" || arg == "-vp" {
			if i == len(args)-3 {
				fmt.Println("The option '--video-port' requires a value")
				os.Exit(1)
			}
			vp, err := strconv.Atoi(args[i+1])
			if err != nil || vp <= 0 {
				fmt.Println("The option '--video-port' requires a numeric value")
				os.Exit(1)
			}
			portVideo = vp
			i++
		} else if arg == "--audio-port" || arg == "-ap" {
			if i == len(args)-3 {
				fmt.Println("The option '--audio-port' requires a value")
				os.Exit(1)
			}
			ap, err := strconv.Atoi(args[i+1])
			if err != nil || ap <= 0 {
				fmt.Println("The option '--audio-port' requires a numeric value")
				os.Exit(1)
			}
			portAudio = ap
			i++
		} else if arg == "--secret" || arg == "-s" {
			if i == len(args)-3 {
				fmt.Println("The option '--secret' requires a value")
				os.Exit(1)
			}
			authSecret = args[i+1]
			i++
		}
	}

	if source == "" {
		fmt.Println("Missing required option: --input")
		os.Exit(1)
	}

	if portVideo == 0 {
		fmt.Println("Missing required option: --video-port")
		os.Exit(1)
	}

	if portAudio == 0 {
		fmt.Println("Missing required option: --video-audio")
		os.Exit(1)
	}

	if portAudio == portVideo {
		fmt.Println("Port for video cannot be the same as the port for audio")
		os.Exit(1)
	}

	if sdpFile == "" {
		fmt.Println("Missing required option: --sdp-file")
		os.Exit(1)
	}

	if forwardMode == "" {
		fmt.Println("Missing required option: --forward-mode")
		os.Exit(1)
	}

	if forwardMode != "TEST" && forwardMode != "RTMP" && forwardMode != "CUSTOM" {
		fmt.Println("Invalid forward mode: " + forwardMode)
		os.Exit(1)
	}

	if forwardMode == "RTMP" {
		forwardParam = os.Getenv("RTMP_FORWARD_URL")
		uSource, err := url.Parse(forwardParam)
		if err != nil || (uSource.Scheme != "rtmp" && uSource.Scheme != "rtmps") {
			fmt.Println("Invalid RTMP URL provided. Please set RTMP_FORWARD_URL to a valid URL when usinmg RTMP forward mode.")
			os.Exit(1)
		}
	} else if forwardMode == "CUSTOM" {
		forwardParam = os.Getenv("CUSTOM_FORWARD_COMMAND")
		if forwardParam == "" {
			fmt.Println("Please set CUSTOM_FORWARD_COMMAND when using CUSTOM forward mode.")
			os.Exit(1)
		}
	}

	uSource, err := url.Parse(source)
	if err != nil || (uSource.Scheme != "ws" && uSource.Scheme != "wss") {
		fmt.Println("The source is not a valid websocket URL")
		os.Exit(1)
	}

	protocolSource := uSource.Scheme
	hostSource := uSource.Host
	streamIdSource := ""

	if len(uSource.Path) > 0 {
		streamIdSource = uSource.Path[1:]
	} else {
		fmt.Println("The source URL must contain the stream ID. Example: ws://localhost/stream-id")
		os.Exit(1)
	}

	wsURLSource := url.URL{
		Scheme: protocolSource,
		Host:   hostSource,
		Path:   "/ws",
	}

	if authSecret != "" {
		authToken = generateToken(authSecret, streamIdSource)
	}

	if _, err := os.Stat(ffmpegPath); err != nil {
		fmt.Println("Error: Could not find 'ffmpeg' at specified location: " + ffmpegPath)
		os.Exit(1)
	}

	runProcess(wsURLSource, streamIdSource, ProcessOptions{
		debug:        debug,
		portAudio:    portAudio,
		portVideo:    portVideo,
		sdpFile:      sdpFile,
		ffmpeg:       ffmpegPath,
		forwardMode:  forwardMode,
		forwardParam: forwardParam,
		authToken:    authToken,
	})
}

func printHelp() {
	fmt.Println("Usage: webrtc-forwarder [OPTIONS]")
	fmt.Println("    OPTIONS:")
	fmt.Println("        --help, -h                              Prints command line options.")
	fmt.Println("        --version, -v                           Prints version.")
	fmt.Println("        --debug                                 Enables debug mode.")
	fmt.Println("        --input, -i <SOURCE>                    Input WebRTC stream. Example: ws(s)://host:port/stream-id")
	fmt.Println("        --sdp-file, -sdp <file>                 File where to print the SDP description.")
	fmt.Println("        --forward-mode, -fm <MODE>              Forward mode can be: TEST, RTMP or CUSTOM.")
	fmt.Println("        --video-port, -vp <port>                Sets the port for video packets.")
	fmt.Println("        --audio-port, -ap <port>                Sets the port for audio packets.")
	fmt.Println("        --ffmpeg-path <path>                    Sets FFMpeg path.")
	fmt.Println("        --auth, -a <auth-token>                 Sets authentication token for the source.")
	fmt.Println("        --secret, -s <secret>                   Sets secret to generate authentication tokens.")
	fmt.Println("    FORWARD MODES:")
	fmt.Println("        --forward-mode TEST                     Creates the SDP file and does nothing else. For testing.")
	fmt.Println("        --forward-mode RTMP                     Forwards the RTC stream to RTMP. Set RTMP_FORWARD_URL env variable.")
	fmt.Println("        --forward-mode CUSTOM                   Runs a custom command to forward the stream. Set CUSTOM_FORWARD_COMMAND env variable.")
}

func printVersion() {
	fmt.Println("webrtc-forwarder 1.0.0")
}
