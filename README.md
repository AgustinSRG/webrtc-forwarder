# WebRTC Forwarder

Utility for [webrtc-cdn](https://github.com/AgustinSRG/webrtc-cdn) to forward WebRTC to other protocol, like RTMP.

It uses [FFMpeg](https://ffmpeg.org/) for the protocol remuxing, and the [pion/webrtc](https://github.com/pion/webrtc) for WebRTC connectivity.

## Compilation

In order to install dependencies, type:

```
go get github.com/AgustinSRG/webrtc-forwarder
```

To compile the code type:

```
go build
```

The build command will create a binary in the currenct directory, called `webrtc-forwarder`, or `webrtc-forwarder.exe` if you are using Windows.

## Usage

You can use the program from the command line:

```
webrtc-forwarder [OPTIONS]
```

### OPTIONS (Required)

Here is a list of the required options:

| Option | Description |
|---|---|
| `--input, -i <input-url>` | Sets the input URL. Example: `ws://localhost/stream-id` |
| `--video-port, -vp <port>` | Port to forward video RTP packets. |
| `--audio-port, -ap <port>` | Port to forward audio RTP packets. |
| `--sdp-file, -sdp <file.sdp>` | File to use to forward the stream. After the connection is stablished, you can use this file as an input of FFMPEG. |
| `--forward-mode, -fm <mode>` | Forward mode, check the section below for mode details. |

### Forward modes

The available forward modes are the following:

| Mode | Description |
|---|---|
| `TEST` | Just setups the SDP file and lets you test it by yourself. |
| `RTMP` | Forwards to RTMP using the envirinment variable `RTMP_FORWARD_URL`. Example: `rtmp://live.twitch.tv/app/$STREAM_KEY` |
| `CUSTOM` | Run a custom command to forward or process the stream. The command must be set in `CUSTOM_FORWARD_COMMAND` environment variable. |

### OPTIONS (Optional)

Here is a list of the rest of the options:

| Option | Description |
|---|---|
| `--help, -h` | Shows the command line options |
| `--version, -v` | Shows the version |
| `--debug` | Enables debug mode (prints more messages) |
| `--ffmpeg-path <path>` | Sets the FFMpeg path. By default is `/usr/bin/ffmpeg`. You can also change it with the environment variable `FFMPEG_PATH` |
| `--auth, -a <auth-token>` | Sets auth token for the source. |
| `--secret, -s <secret>` | Provides secret to generate authentication tokens. |

## WebRTC options

You can configure WebRTC configuration options with environment variables:

| Variable Name | Description |
|---|---|
| STUN_SERVER | STUN server URL. Example: `stun:stun.l.google.com:19302` |
| TURN_SERVER | TURN server URL. Set if the server is behind NAT. Example: `turn:turn.example.com:3478` |
| TURN_USERNAME | Username for the TURN server. |
| TURN_PASSWORD | Credential for the TURN server. |
