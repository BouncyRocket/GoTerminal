package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/inancgumus/screen"
)

type Client struct {
	conn     net.Conn
	username string
}

var (
	clients = make(map[net.Conn]*Client)
	mutex   sync.Mutex
)

var ipserver = "global"
var portserver = "global"
var passwordserver = "global"

func main() {
	for {
		screen.Clear()
		showMenuserver()
		option := readOptionserver()
		switch option {
		case 1:
			fmt.Println("Settings")
			settingsMenuserver()
		case 2:
			fmt.Println("Run")
			Run()
		case 3:
			fmt.Println("Exit")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func showMenuserver() {
	screen.Clear()
	fmt.Println("=====")
	fmt.Println("Menu")
	fmt.Println("")
	fmt.Println("Settings [1]")
	fmt.Println("Run [2]")
	fmt.Println("Exit [3]")
	fmt.Print("\n>>: ")
}

func readOptionserver() int {
	var option int
	fmt.Scanln(&option)
	return option
}

func settingsMenuserver() {
	screen.Clear()
	fmt.Println("=====")
	fmt.Println("Settings")
	fmt.Println("")
	fmt.Println("Ip [1]")
	fmt.Println("Port [2]")
	fmt.Println("Password [3]")
	fmt.Println("Back [0]")
	fmt.Print("\n>>: ")
	option := readOptionserver()
	switch option {
	case 1:
		fmt.Print("Enter IP: ")
		fmt.Scanln(&ipserver)
		settingsMenuserver()
	case 2:
		fmt.Print("Enter Port: ")
		fmt.Scanln(&portserver)
		screen.Clear()
		settingsMenuserver()
	case 3:
		fmt.Print("Enter Password: ")
		fmt.Scanln(&passwordserver)
		screen.Clear()
		settingsMenuserver()
	case 0:
		return
	default:
		fmt.Println("Invalid option. Please try again.")
		settingsMenuserver()
	}
}

func Run() {
	listener, err := net.Listen("tcp", ipserver+":"+portserver)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on " + ipserver + ":" + portserver)

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

	if password != passwordserver {
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

		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "?") {
			// Ignore the command if it starts with "?#"
			if strings.HasPrefix(message, "?#") {
				continue
			}

			// Handle other commands here
			if strings.HasPrefix(message, "?kick ") {
				username := strings.TrimSpace(strings.TrimPrefix(message, "?kick "))
				kickUser(username)
			}
			// Add more commands as needed
		}

		// Broadcast regular messages to clients
		formattedMessage := fmt.Sprintf("[%s]: %s", client.username, message)
		broadcastMessage(formattedMessage, conn)
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
			fmt.Println(message + "\n")
			if err != nil {
				fmt.Printf("Error broadcasting message to client [%s]: %s\n", client.username, err)
			}
		}
	}
}

func kickUser(username string) {
	mutex.Lock()
	defer mutex.Unlock()

	for conn, client := range clients {
		if client.username == username {
			conn.Close()
			delete(clients, conn)
			broadcastMessage(fmt.Sprintf("User %s has been kicked by the server", username), nil)
			fmt.Printf("User %s has been kicked by the server\n", username)
			return
		}
	}

	fmt.Printf("User %s not found\n", username)
	broadcastMessage(fmt.Sprintf("User %s not found", username), nil)
}
