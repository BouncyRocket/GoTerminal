package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/inancgumus/screen"
)

var Ip = "global"
var port = "global"

const password = ""

func main() {
	for {
		screen.Clear()
		showMenu()
		option := readOption()
		switch option {
		case 1:
			fmt.Println("Settings")
			settingsMenu()
		case 2:
			fmt.Println("Join")
			join()
		case 3:
			fmt.Println("Exit")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func showMenu() {
	screen.Clear()
	fmt.Println("=====")
	fmt.Println("Menu")
	fmt.Println("")
	fmt.Println("1. Settings")
	fmt.Println("2. Join")
	fmt.Println("3. Exit")
	fmt.Print("\n>>: ")
}

func readOption() int {
	var option int
	_, err := fmt.Scanln(&option)
	if err != nil {
		log.Println("Error reading option:", err)
	}
	return option
}

func settingsMenu() {
	screen.Clear()
	fmt.Println("=====")
	fmt.Println("Settings")
	fmt.Println("")
	fmt.Println("1. IP")
	fmt.Println("2. Port")
	fmt.Println("0. Back")
	fmt.Print("\n>>: ")
	option := readOption()
	switch option {
	case 1:
		fmt.Print("Enter IP: ")
		fmt.Scanln(&Ip) // Assign directly to the global variable Ip
		screen.Clear()
		settingsMenu()
	case 2:
		fmt.Print("Enter Port: ")
		fmt.Scanln(&port) // Assign directly to the global variable port
		screen.Clear()
		settingsMenu()
	case 0:
		main()
	default:
		fmt.Println("Invalid option. Please try again.")
		screen.Clear()
		settingsMenu()
	}
}

func join() {
	conn, err := net.Dial("tcp", Ip+":"+port)
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

	if password != password {
		fmt.Println("Invalid password")
		return
	}

	fmt.Println("Password accepted. You are now connected to the server.")

	// Start listening for user input
	go listenForInput(conn)

	// Read and display messages from the server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		log.Println("[Server]:", message)
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
			log.Println("Error writing to server:", err)
			break
		}
		err = writer.Flush()
		if err != nil {
			log.Println("Error flushing writer:", err)
			break
		}
	}

	if scanner.Err() != nil {
		log.Println("Error reading from input:", scanner.Err())
	}
}
