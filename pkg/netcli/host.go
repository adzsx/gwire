package netcli

import (
	"bufio"
	"crypto/x509"
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
	wg sync.WaitGroup
)

// Set up listener for each port on list
func HostSetup(input utils.Input) {
	// Global slice for distributing messages
	var message = [][]string{}

	if input.Enc == "auto" {
		var err error
		input.Enc, err = crypt.GenPasswd(32)
		utils.Err(err)
	}
	// Set up listener for every port in range
	for _, port := range input.Port {

		// wg = WaitGroup (Variable to wait until variable hits 0)
		wg.Add(1)
		conn := listen(input, port)
		if input.Enc == "auto" {
			hostRSA(input, conn)
		}
		go host(input, conn, port, &message)
	}

	// Wait untill wg is 0
	wg.Wait()
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

	log.Println("Listening on port", port)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection:", err.Error())
		wg.Done()
		return conn
	}

	log.Printf("Connected to %v", conn.LocalAddr())

	return conn
}

// Listener loop for individual port
func host(input utils.Input, conn net.Conn, port string, message *[][]string) {

	// Read data
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
					log.Fatalf("Connection on port %v closed", port)
					wg.Done()
					return
				}
				log.Fatalln("Error reading data:", err.Error())

			}
			if len([]byte(input.Enc)) != 0 {
				log.Print(crypt.DecryptAES(data, []byte(input.Enc)))
				// log.Print(crypt.DecryptAES(data, input.Key))
			} else {
				log.Print(data)
			}

			// if len(input.Port) > 1 {
			// 	for i := 0; i < len(input.Port)-1; i++ {
			// 		*message = append(*message, []string{utils.FilterPort(conn.LocalAddr().String()), data})
			// 	}
			// }
		}

	}()

	// Send data from input
	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(input.TimeOut))
			reader := bufio.NewReader(os.Stdin)

			// attach username
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			if len(text) > 16834 {
				log.Println("Message cant be over 16834 characters long")
			}

			if len(input.Port) > 1 {
				for i := 0; i < len(input.Port); i++ {
					*message = append(*message, []string{"0", text})
				}
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
				time.Sleep(time.Millisecond * time.Duration(input.TimeOut))
				if len(*message) > 0 {
					for index, element := range *message {
						if element[0] != utils.FilterPort(conn.LocalAddr().String()) {
							conn.Write([]byte(element[1]))
							*message = utils.Remove(*message, index)
						}
					}
				}
			}
		}()
	}
}

func hostRSA(input utils.Input, conn net.Conn) {

	// Make buffer for receiving RSA public key

	utils.VPrint("Waiting for RSA public from "+utils.FilterIp(conn.LocalAddr().String())+"\n", 2)
	buffer := make([]byte, 4096)
	bytes, err := conn.Read(buffer)
	utils.Err(err)
	sentPublicKey := buffer[:bytes]

	utils.VPrint("Publickey received from "+utils.FilterIp(conn.LocalAddr().String())+"\n", 1)

	// Convert bytes back to public key
	publicKey, err := x509.ParsePKCS1PublicKey(sentPublicKey)
	utils.Err(err)

	// Send encrypted AES key over connection
	encKey := crypt.EncryptRSA(*publicKey, []byte(input.Enc))
	conn.Write(encKey)
}
