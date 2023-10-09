package netcli

import (
	"bufio"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/adzsx/gwire/pkg/crypt"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	conn net.Conn
	err  error
)

// Function connects to host with TCP
func ClientSetup(input utils.Input) {
	log.SetFlags(0)
	if input.Time {
		log.SetFlags(log.Ltime)
	}

	// Connect to host

	if input.Ip != "scan" {
		conn, err = net.Dial("tcp", input.Ip+":"+input.Port[0])
	} else {
		// Scan every host in network for open port
		hosts, err := GetHosts(Subnet())

		utils.Err(err, true)
		input.Ip, conn = ScanRange(hosts, input.Port[0])
	}

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		utils.Err(errors.New("connection refused by destination"), true)
	}

	utils.Print("Connected to "+input.Ip+":"+input.Port[0]+"\n", 0)

	if input.Enc == "auto" {
		input, err = initClient(input, conn)
		utils.Err(err, true)
	}

	utils.Print("Setup finished", 1)
	client(input, conn)
}

// Function for ongoing connection
func client(input utils.Input, conn net.Conn) {
	// Receive Data

	utils.Print("Started client routine", 3)

	// Starting ui
	utils.Ansi("\033[S")
	utils.Ansi("\033G")

	// Receive Data
	var received []string

	go func() {
		for {
			//Read data
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			received = append(received, string(buffer[:bytes]))

			if err != nil {
				if err.Error() == "EOF" {
					utils.Print("Connection closed by remote host", 0)
					os.Exit(0)
				}
				log.Fatalln("Error reading data:", err.Error())

			}

		}
	}()

	// Function for printing received data
	go func() {
		var data string

		for {
			time.Sleep(time.Millisecond * time.Duration(input.TimeOut))

			if len(received) != 0 {

				utils.Ansi("\x1b[s")
				utils.Ansi("\033[1A\033[999D\033[K")
				utils.Ansi("\033[96m")

				if len([]byte(input.Enc)) != 0 {
					data = crypt.DecryptAES(received[0], []byte(input.Enc)) + "\n"
				} else {
					data = received[0] + "\n"
				}

				color := utils.GetRandomString(colorList, utils.FilterChar(data, ">", true))
				fmt.Print(color)
				fmt.Print(data)

				utils.Ansi("\033[0m\x1b[u\033[2A")

				received = utils.Remove(received, received[0])

			}

		}

	}()

	// Send data
	func() {
		reader := bufio.NewReader(os.Stdin)

		log.SetFlags(0)
		if input.Time {
			log.SetFlags(log.Ltime)
		}

		for {
			time.Sleep(time.Millisecond * 100)

			text := input.Username + "> "

			inp, _ := reader.ReadString('\n')
			text += inp

			// Move up one line, Clear it. Again. Print in blue
			utils.Ansi("\033[F\033[0K\033[F\033[0K\033[37m")

			fmt.Println(text)

			// Move back down, print in white
			utils.Ansi("\033[2B\033[0m")

			if len(text) > 16384 {
				log.Println("Message cant be over 16384 characters long")
				break
			}

			if len([]byte(input.Enc)) != 0 {
				conn.Write([]byte(crypt.EncryptAES(text, []byte(input.Enc))))
			} else {
				conn.Write([]byte(text))
			}

		}
	}()

}

// Func for setting up RSA encryption for the clientcs
func initClient(input utils.Input, conn net.Conn) (utils.Input, error) {

	utils.Print("Generating RSA Keys", 1)
	var rsaKeys = crypt.GenKeys()

	byteKey := x509.MarshalPKCS1PublicKey(&rsaKeys.PublicKey)

	utils.Print("Sending Public Key", 2)
	conn.Write(byteKey)

	// Wait for host to send password
	utils.Print("Waiting for response", 2)
	buffer := make([]byte, 512)
	bytes, err := conn.Read(buffer)
	if err != nil {
		log.Println("Connection unexpectedly closed. Aborting Setup")
		return input, errors.New("connection failed")
	}
	data := string(buffer[:bytes])

	passwd := crypt.DecryptRSA(rsaKeys, []byte(data))

	input.Enc = passwd

	utils.Print("Password received", 1)
	utils.Print("Seinding Password confirmation package", 2)

	conn.Write([]byte(crypt.EncryptAES(input.Enc, []byte(input.Enc))))

	utils.Print("Waiting for Control package", 2)

	buffer = make([]byte, 512)
	bytes, err = conn.Read(buffer)
	utils.Err(err, true)
	data = string(buffer[:bytes])
	data = crypt.DecryptAES(data, []byte(input.Enc))

	utils.Print("Received Control Package", 2)

	if string(data) == "wrong password" {
		log.Println("Wrong password. Aborting connection")
		return input, errors.New("wrong password")
	}

	return input, nil
}
