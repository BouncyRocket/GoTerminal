package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"io"
	"log"
	"net/http"

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
	fmt.Println("Settings [1]")
	fmt.Println("Join [2]")
	fmt.Println("Exit [3]")
	fmt.Print("\n>>: ")
}

func readOption() int {
	var option int
	fmt.Scanln(&option)
	return option
}

func settingsMenu() {
	screen.Clear()
	fmt.Println("=====")
	fmt.Println("Settings")
	fmt.Println("")
	fmt.Println("Ip [1]")
	fmt.Println("Port [2]")
	fmt.Println("Back [0]")
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
		command := scanner.Text()
		handleCommand(command)
		fmt.Println(">>: ", command)
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

func handleCommand(command string) {
	switch {
	case strings.HasPrefix(command, "cdownload"):
		// Extract the file name from the command
		fileName := strings.TrimSpace(strings.TrimPrefix(command, "cdownload"))

		// Download the file
		err := downloadFile(fileName)
		if err != nil {
			log.Println("Error downloading file:", err)
			return
		}

		fmt.Printf("File '%s' downloaded successfully.\n", fileName)
	}
}

func downloadFile(fileName string) error {
	// Open the file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download the file using HTTP request
	resp, err := http.Get(fileName)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
