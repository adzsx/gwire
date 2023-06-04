package netcli

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/adzsx/g-wire/pkg/utils"
)

func Connect(input utils.Input) {

	// Connect to host
	var conn, err = net.Dial("tcp", input.Ip+":"+input.Port)

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		fmt.Println("Connection refused by destination")
		os.Exit(1)
	}

	fmt.Println("Connected to", input.Ip+":"+input.Port)

	// Receive Data
	go func() {
		for {
			// Scan line until \n
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Connection stopped by remote host")
					os.Exit(1)
				}
				fmt.Println("Error reading data:", err.Error())

			}
			fmt.Print(data)
		}

	}()

	// Send data
	func() {
		for {
			reader := bufio.NewReader(os.Stdin)

			// attach username
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			conn.Write([]byte(text))
		}
	}()

}

func Listen(input utils.Input) {

	ln, err := net.Listen("tcp", ":"+input.Port)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		fmt.Println("Permission denied.\nTry again with root or take a port above 1023")
		os.Exit(1)
	}

	fmt.Println("Listening on port", input.Port)

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err.Error())
		return
	}

	fmt.Println("Connection established")

	// Read data
	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Connection stopped by remote host")
					os.Exit(1)
				}
				fmt.Println("Error reading data:", err.Error())

			}
			fmt.Print(data)
		}

	}()

	// Send data
	func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			conn.Write([]byte(text))
		}
	}()
}
