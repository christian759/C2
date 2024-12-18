package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func handle_client(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	command := bufio.NewReader(os.Stdin)

	for {
		text, _ := command.ReadString('\n')

		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Command sent to client:", text)

		// Set a timeout for the client response (e.g., 10 seconds)
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		// Read a fixed number of bytes (e.g., 1024 bytes at once)
		buf := make([]byte, 1048576)
		n, err := reader.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Timeout occurred waiting for client response")
			} else {
				fmt.Println(err)
			}
		}

		// Convert the buffer to string and print it
		fullMessage := string(buf[:n])
		fmt.Println("Client Response: ", fullMessage)
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Panic(err)
	}

	defer listener.Close()

	fmt.Println("Server is listening on port 8000...")

	// Accept incoming connections
	for {
		// net.KeepAliveConfig(true, 300000, 30, 3000)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("New Client Connected: ", conn.RemoteAddr())

		// Handle each connection in a new goroutine
		go handle_client(conn)

	}
}
