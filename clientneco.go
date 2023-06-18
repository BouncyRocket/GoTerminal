package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ServerAddr = "192.168.3.20:8888" // Update this with your server address / put the server ip here / the server you wanta join
	Password   = "(dont change this this is not for the password of the server but if remove this it crashes)"
)

func main() {
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Enter password: ")
	password, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	password = strings.TrimSpace(password)

	if password != Password {
		fmt.Println("Invalid password")
		return
	}

	fmt.Println("Password accepted. You are now connected to the server.")

	// Start listening for user input
	go listenForInput(conn)

	// Read and display messages from the server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading from server:", scanner.Err())
	}
}

func listenForInput(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		message := scanner.Text()

		_, err := writer.WriteString(message + "\n")
		if err != nil {
			fmt.Println("Error writing to server:", err)
			break
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer:", err)
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading from input:", scanner.Err())
	}
}
