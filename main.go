package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Start listening on port 21
	listener, err := net.Listen("tcp", ":21")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("FTP Server Listening on :21")

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		fmt.Println("New connection accepted.")

		// Handle connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Welcome message
    conn.Write([]byte("220 Welcome to My FTP Server\r\n"))

    // Read commands from client
    buf := make([]byte, 1024)
    var username string
    for {
        n, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading:", err.Error())
            return
        }
        cmd := strings.TrimSpace(string(buf[:n]))
        fmt.Println("Received command:", cmd)

        // Respond based on the received command
        if strings.HasPrefix(cmd, "USER") {
            parts := strings.Split(cmd, " ")
            if len(parts) < 2 {
                conn.Write([]byte("501 Syntax error in parameters or arguments.\r\n"))
                continue
            }
            username = parts[1]
            conn.Write([]byte("331 User name okay, need password.\r\n"))

        } else if strings.HasPrefix(cmd, "PASS") {
            if username == "" {
                conn.Write([]byte("503 Bad sequence of commands. Send USER first.\r\n"))
                continue
            }
            parts := strings.Split(cmd, " ")
            if len(parts) < 2 {
                conn.Write([]byte("501 Syntax error in parameters or arguments.\r\n"))
                continue
            }
            password := parts[1]
            conn.Write([]byte("230 User logged in, proceed.\r\n"))
            fmt.Println("Username:", username, "Password:", password)
            return

        } else if strings.HasPrefix(cmd, "QUIT") {
            conn.Write([]byte("221 Goodbye.\r\n"))
            return

        } else {
            conn.Write([]byte("502 Command not implemented.\r\n"))
        }
    }
}

