package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn     net.Conn
	username string
}

var (
	clients = make(map[net.Conn]*Client)
	mutex   sync.Mutex
)

const (
	ServerAddress = "localhost:8888" // Update this with your server address / your Ip here
	Password      = "mypassword"     // Update this with your desired password / for clients joining the server
)

func main() {
	listener, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on " + ServerAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	password = strings.TrimSpace(password)

	if password != Password {
		_, err := conn.Write([]byte("Invalid password. Disconnecting...\n"))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
		fmt.Println("Invalid password")
		return
	}

	fmt.Println("Password accepted. New client connected:", conn.RemoteAddr())

	// Ask the client for a username
	_, err = conn.Write([]byte("Enter your username: "))
	if err != nil {
		fmt.Println("Error writing to client:", err)
		return
	}

	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username = strings.TrimSpace(username)

	client := &Client{
		conn:     conn,
		username: username,
	}

	// Add the client to the clients map
	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	fmt.Printf("New client joined: %s\n", client.username)

	// Notify all clients about the new client's connection
	broadcastMessage(fmt.Sprintf("Allowed connection from %s", conn.RemoteAddr()), nil)

	// Handle client messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message from client:", err)
			break
		}

		message = fmt.Sprintf("[%s]: %s", client.username, strings.TrimSpace(message))
		broadcastMessage(message, conn)
	}

	// Remove the client from the clients map and notify others
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	fmt.Printf("Client disconnected: %s\n", client.username)
	broadcastMessage(fmt.Sprintf("Client disconnected: %s", client.username), nil)
}

func broadcastMessage(message string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn, client := range clients {
		if conn != sender {
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Printf("Error broadcasting message to client [%s]: %s\n", client.username, err)
			}
		}
	}
}
