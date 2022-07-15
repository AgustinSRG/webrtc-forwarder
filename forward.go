// Code to forward the track

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

func createForwardSDPFile(fileName string, videoPort int, audioPort int) string {

	nl := "\n"

	sdpFileContents := "v=0" + nl +
		"o=- 0 0 IN IP4 127.0.0.1" + nl +
		"s=Pion WebRTC" + nl +
		"c=IN IP4 127.0.0.1" + nl +
		"t=0 0" + nl +
		"m=audio " + fmt.Sprint(audioPort) + " RTP/AVP 111" + nl +
		"a=rtpmap:111 OPUS/48000/2" + nl +
		"m=video " + fmt.Sprint(videoPort) + " RTP/AVP 96" + nl +
		"a=rtpmap:96 VP8/90000"

	err := os.WriteFile(fileName, []byte(sdpFileContents), 0644)
	if err != nil {
		panic(err)
	}

	return fileName
}

func forwardTrack(track *webrtc.TrackRemote, port int) {
	// Set payload type
	var payloadType uint8

	if track.Kind() == webrtc.RTPCodecTypeVideo {
		payloadType = 96
	} else {
		payloadType = 111
	}

	// Create a local addr
	var laddr *net.UDPAddr
	var err error = nil

	if laddr, err = net.ResolveUDPAddr("udp", "127.0.0.1:"); err != nil {
		panic(err)
	}

	// Create remote addr
	var raddr *net.UDPAddr
	if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
		panic(err)
	}

	// Dial udp
	var conn *net.UDPConn
	if conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
		panic(err)
	}
	defer func(conn net.PacketConn) {
		if closeErr := conn.Close(); closeErr != nil {
			panic(closeErr)
		}
	}(conn)

	b := make([]byte, 1500)
	rtpPacket := &rtp.Packet{}
	for {
		// Read
		n, _, readErr := track.Read(b)
		if readErr != nil {
			return
		}

		// Unmarshal the packet and update the PayloadType
		if err = rtpPacket.Unmarshal(b[:n]); err != nil {
			panic(err)
		}
		rtpPacket.PayloadType = payloadType

		// Marshal into original buffer with updated PayloadType
		if n, err = rtpPacket.MarshalTo(b); err != nil {
			panic(err)
		}

		// Write
		if _, err = conn.Write(b[:n]); err != nil {
			// For this particular example, third party applications usually timeout after a short
			// amount of time during which the user doesn't have enough time to provide the answer
			// to the browser.
			// That's why, for this particular example, the user first needs to provide the answer
			// to the browser then open the third party application. Therefore we must not kill
			// the forward on "connection refused" errors
			if opError, ok := err.(*net.OpError); ok && opError.Err.Error() == "write: connection refused" {
				continue
			} else {
				return
			}
		}
	}
}
