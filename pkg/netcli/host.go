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
	"sync"
	"time"

	"github.com/adzsx/gwire/pkg/crypt"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	auto      bool
	wg        sync.WaitGroup
	sent      int
	colorList []string = []string{"\033[92m", "\033[94m", "\033[95m", "\033[96m"}
)

// Set up listener for each port on list
func HostSetup(input utils.Input) {
	// Global slice for distributing messages
	var message = [][]string{}

	if input.Enc == "auto" {
		auto = true
		var err error
		utils.Print("Generating password\n", 2)
		input.Enc, err = crypt.GenPasswd()
		utils.Err(err, true)
	}
	// Set up listener for every port in range
	for _, port := range input.Port {

		// wg = WaitGroup (Variable to wait until variable hits 0)
		wg.Add(1)

		go connSetup(input, string(port), &message)

	}

	// Wait untill wg is 0
	wg.Wait()

	defer os.Exit(0)
}

func connSetup(input utils.Input, port string, message *[][]string) {

	conn := listen(input, port)

	if auto {
		err := InitConn(input, conn)
		if err != nil {
			log.Println(err)
			return
		}
	}

	utils.Print("Setup finished\n", 1)

	go host(input, conn, port, message)

}

// Set up connection to specific port
func listen(input utils.Input, port string) net.Conn {

	log.SetFlags(0)
	if input.Time {
		log.SetFlags(log.Ltime)
	}

	//Listen and connect
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		log.Fatalln("Permission denied.\nTry again with root or take a port above 1023")
		wg.Done()
		os.Exit(0)
	}

	utils.Print("Listening on port "+port, 1)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection:", err.Error())
		wg.Done()
		return conn
	}

	utils.Print("Connected to "+conn.RemoteAddr().String(), 0)

	return conn
}

// Listener loop for individual port
func host(input utils.Input, conn net.Conn, port string, message *[][]string) {

	utils.Print("Started host routine", 3)

	fmt.Println()

	// Read data
	var receivedHost []string
	go func() {
		for {
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			data := string(buffer[:bytes])
			receivedHost = append(receivedHost, data)

			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("Connection on port %v closed", port)
					wg.Done()
					return
				} else {
					log.Fatalln("Error reading data:", err.Error())
				}

			}

			if len(input.Port) > 1 {
				*message = append(*message, []string{utils.FilterChar(conn.LocalAddr().String(), ":", false), data})
			}
		}

	}()

	// Function for printing messages
	go func() {
		var data string
		for {
			time.Sleep(time.Millisecond * time.Duration(input.TimeOut))

			if len(receivedHost) != 0 {
				fmt.Print("\x1b[s")
				fmt.Print("\033[1A\033[999D\033[K")
				if len([]byte(input.Enc)) != 0 {
					data = crypt.DecryptAES(receivedHost[0], []byte(input.Enc)) + "\n"
				} else {
					data = receivedHost[0] + "\n"
				}

				color := utils.GetRandomString(colorList, utils.FilterChar(data, ">", true))
				fmt.Print(color)
				fmt.Print(data)

				fmt.Print("\033[0m\x1b[u\033[B")

				receivedHost = utils.Remove(receivedHost, receivedHost[0])
			}
		}
	}()

	// Send data from input
	go func() {
		reader := bufio.NewReader(os.Stdin)
		log.SetFlags(0)
		if input.Time {
			log.SetFlags(log.Ltime)
		}

		for {

			// attach username
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			fmt.Print("\033[F")
			fmt.Print("\033[0K")
			fmt.Print("\033[F")
			fmt.Print("\033[0K")
			fmt.Print("\033[37m")

			fmt.Println(text)

			// Move back down, print in white
			fmt.Print("\033[2B")
			fmt.Print("\033[37m")

			if len(input.Port) > 1 {

				if len(input.Enc) != 0 {

					*message = append(*message, []string{"0", crypt.EncryptAES(text, []byte(input.Enc))})
				} else {
					*message = append(*message, []string{"0", text})
				}

				sent = -1

			} else {

				if len([]byte(input.Enc)) != 0 {
					conn.Write([]byte(crypt.EncryptAES(text, []byte(input.Enc))))
				} else {
					conn.Write([]byte(text))
				}
			}
		}
	}()

	//send data from other clients
	if len(input.Port) > 1 {
		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				if len(*message) > 0 {

					for _, element := range *message {
						if element[0] != utils.FilterChar(conn.LocalAddr().String(), ":", false) {
							conn.Write([]byte(element[1]))
							sent += 1
						}
						time.Sleep(time.Millisecond * 50)
					}
					if sent == len(input.Port)-1 {
						*message = [][]string{}
						sent = 0
					}

				}
			}
		}()
	}
}

func InitConn(input utils.Input, conn net.Conn) error {
	// Make buffer for receiving RSA public key
	utils.Print("Waiting for RSA key from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)
	buffer := make([]byte, 4096)
	bytes, err := conn.Read(buffer)
	if err != nil {
		wg.Done()
		return errors.New("connection closed")
	}
	sentPublicKey := buffer[:bytes]

	// Convert bytes back to public key
	publicKey, err := x509.ParsePKCS1PublicKey(sentPublicKey)

	if err != nil {
		wg.Done()
		return errors.New("received data not RSA publickey")
	}

	utils.Print("Publickey received from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)

	// Send encrypted AES key over connection
	utils.Print("Sending Password", 2)
	encKey := crypt.EncryptRSA(*publicKey, []byte(input.Enc))
	conn.Write(encKey)

	utils.Print("Waiting for password confirmation", 2)

	buffer = make([]byte, 512)
	bytes, err = conn.Read(buffer)
	utils.Err(err, true)
	data := string(buffer[:bytes])

	utils.Print("Received password confirmation", 2)

	if crypt.DecryptAES(data, []byte(input.Enc)) != input.Enc {
		conn.Write([]byte("wrong password"))
		return errors.New("wrong password")
	}

	conn.Write([]byte(crypt.EncryptAES("success", []byte(input.Enc))))

	return nil
}
