package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {
	// specify a logging standard
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// specify a context
	ctx := context.Background()

	// specify a socket to listen on. socket = address:port, defaulting to 0.0.0.0
	addr, err := net.ResolveUDPAddr("udp", ":34567")
	// handle error wrt address resolution
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Failed to resolve UDP address",
			slog.Int("address", 34567),
			slog.String("network", "udp"),
			slog.Any("error", err),
		)
		return
	}

	// start listening for UDP packets on the socket
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"Failed to start UDP listener",
			slog.String("address", addr.String()),
			slog.String("network", addr.Network()),
			slog.Any("error", err),
		)
		// exit if UDP cannot get a listener
		os.Exit(1)
	}

	// close UDP connection
	defer conn.Close()

	// log the starting message
	logger.InfoContext(ctx, "Starting UDP Server on addr: ", slog.String("address", addr.String()))

	// buffer to hold incoming data
	buf := make([]byte, 1024)

	// listen for packets in a loop
	for {
		// generate a unique ID for each request
		requestID := uuid.New().String()

		// read bytes from buffer
		p, remoteAddr, err := conn.ReadFromUDP(buf)
		// record time
		startTime := time.Now()
		if err != nil {
			logger.ErrorContext(ctx, "Failed to read bytes from buffer", slog.Any("error", err))
			continue
		}

		// acknowledge reciept/send something back
		message := fmt.Sprintf("Server received this message: %s\n", strings.ToUpper(string(buf[:p])))
		// send a response back
		_, err = conn.WriteToUDP([]byte(message), remoteAddr)
		if err != nil {
			logger.ErrorContext(
				ctx,
				"Failed to send reply to remote addr",
				slog.String("remote_address", remoteAddr.String()),
				slog.String("network", remoteAddr.Network()),
				slog.String("request_id", requestID),
				slog.Any("error", err),
			)
		}

		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		logger.InfoContext(
			ctx,
			"Processed UDP packet",
			slog.String("request_id", requestID),
			slog.String("remote_address", remoteAddr.String()),
			slog.Int("bytes_received", p),
			slog.Duration("duration_ms", time.Duration(elapsedTime.Milliseconds())),
			slog.Any("data", string(buf[:p])),
		)
	}
}
