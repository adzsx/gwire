package netcli

import (
	"bufio"
	"crypto/x509"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/adzsx/gwire/pkg/crypt"
	"github.com/adzsx/gwire/pkg/utils"
)

// Function connects to host with TCP
func ClientSetup(input utils.Input) {
	log.SetFlags(0)
	if input.Time {
		log.SetFlags(log.Ltime)
	}

	// Connect to host
	var conn, err = net.Dial("tcp", input.Ip+":"+input.Port[0])

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		log.Fatalln("Connection refused by destination")
		os.Exit(0)
	}

	log.Println("Connected to", input.Ip+":"+input.Port[0]+"\n")

	if input.Enc == "auto" {
		input.Enc = clientRSA(input, conn)
	}

	utils.VPrint("Setup finished\n")
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
func clientRSA(input utils.Input, conn net.Conn) string {

	utils.VPrint("Generating RSA Keys")
	var rsaKeys = crypt.GenKeys()

	byteKey := x509.MarshalPKCS1PublicKey(&rsaKeys.PublicKey)

	utils.VPrint("Sending Public Key")
	conn.Write(byteKey)

	// Wait for host to send password
	utils.VPrint("Waiting for response...")
	buffer := make([]byte, 512)
	bytes, err := conn.Read(buffer)
	utils.Err(err)
	data := buffer[:bytes]

	passwd := crypt.DecryptRSA(rsaKeys, data)

	utils.VPrint("Received Password\n")

	return passwd
}
