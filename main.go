package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const HOST = "http://127.0.0.1:8000/api/ftps/"

var (
	HOST_PASSWORD = "admin"
	HOST_USERNAME = "admin"
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

func logInfo(conn net.Conn, username string, password string, pwned bool) {
	data := map[string]interface{}{
		"remoteAddr":     conn.RemoteAddr().String(),
		"username":       username,
		"password":       password,
		"client_version": "1.0",
		"pwned":          pwned,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Create a new HTTP POST request with JSON body
	req, err := http.NewRequest("POST", HOST, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Username:", username, "Password:", password, "removeAddress:", conn.RemoteAddr().String())
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Welcome message
	conn.Write([]byte("220 Welcome to FTP Server\r\n"))

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
			conn.Write([]byte("503 Bad Credentials.\r\n"))
			pwned := HOST_USERNAME == username && HOST_PASSWORD == password
			logInfo(conn, username, password, pwned)
			return

		} else if strings.HasPrefix(cmd, "QUIT") {
			conn.Write([]byte("221 Goodbye.\r\n"))
			return

		} else {
			conn.Write([]byte("502 Command not implemented.\r\n"))
		}
	}
}
