package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

func main() {
	// specifying logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// command-line flags for getting server address
	serverAddr := flag.String("addr", "127.0.0.1", "UDP server IP address")
	serverPort := flag.String("port", "", "UDP server port number")

	// specify a custom error handling
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
		fmt.Println("\t-addr\tinput your server address, defaults to 127.0.0.1")
		fmt.Println("\t-port\tspecify the port number")
	}

	flag.Parse()

	if *serverPort == "" {
		flag.Usage()
		os.Exit(1)
	}

	serverSocket := *serverAddr + ":" + *serverPort

	// resolve the server address
	addr, err := net.ResolveUDPAddr("udp", serverSocket)
	if err != nil {
		logger.Error("Failed to resolve UDP server address", "error", err)
		os.Exit(1)
	}

	// create connection to send packets to UDP server
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		logger.Error("Failed to dial UDP server address", "error", err)
		os.Exit(1)
	}

	// close connection
	defer conn.Close()

	logger.Info("UDP client initialized a connection", slog.String("server_address", addr.String()))

	// read input from user - from Stdin
	reader := bufio.NewReader(os.Stdin)

	// loop to send message
	for {
		fmt.Print("<<-------------------------<<>>--------------------------->>"
		fmt.Print("Enter message to send (or type 'quit'): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Error reading input string", "error", err)
			continue
		}

		message := strings.TrimSpace(input)

		if strings.EqualFold(message, "quit") {
			logger.Info("Exiting client", slog.String("message", "quit"))
			break
		}

		// write the message into the connection
		_, err = conn.Write([]byte(message))
		if err != nil {
			logger.Error("Failed to write message into connection", "message", message, "error", err)
			continue
		}
		logger.Info("Message sent!", "message", message)

		// awaiting a response from server
		buf := make([]byte, 1024)

		// read response from server
		n, err := conn.Read(buf)
		if err != nil {
			logger.Error("Error reading response", "error", err)
			continue
		}
		response := string(buf[:n])
		logger.Info("Received response", "response", response)
	}
}
