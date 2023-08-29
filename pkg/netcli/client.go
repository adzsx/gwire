package netcli

import (
	"bufio"
	"crypto/x509"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
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
		log.Fatalln("Connection refused by destination")
		os.Exit(0)
	}

	log.Printf("Connected to %v", input.Ip+":"+input.Port[0]+"\n")

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
	go func() {
		for {

			time.Sleep(time.Millisecond * time.Duration(input.TimeOut))

			//Read data
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			data := string(buffer[:bytes])

			if err != nil {
				if err.Error() == "EOF" {
					log.Fatalln("Connection closed by remote host")
					os.Exit(0)
				}
				log.Fatalln("Error reading data:", err.Error())

			}

			if len([]byte(input.Enc)) != 0 {
				log.Print(crypt.DecryptAES(data, []byte(input.Enc)))
			} else {
				log.Print(data)
			}

		}
	}()

	// Send data
	func() {
		log.SetFlags(0)
		if input.Time {
			log.SetFlags(log.Ltime)
		}

		for {
			time.Sleep(time.Millisecond * time.Duration(input.TimeOut))
			reader := bufio.NewReader(os.Stdin)

			// attach username
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

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

	timeout, err := strconv.Atoi(string(data))
	utils.Err(err, true)

	input.TimeOut = float64(timeout)

	return input, nil
}
