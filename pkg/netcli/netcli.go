package netcli

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/adzsx/g-wire/pkg/utils"
)

var (
	wg sync.WaitGroup
)

func Connect(input utils.Input) {
	log.SetFlags(0)
	if input.Time {
		log.SetFlags(log.Ltime)
	}

	// Connect to host
	var conn, err = net.Dial("tcp", input.Ip+":"+input.Port[0])

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		log.Fatalln("Connection refused by destination")
		os.Exit(1)
	}

	log.Println("Connected to", input.Ip+":"+input.Port[0])

	// Receive Data
	go func() {
		log.SetFlags(0)
		if input.Time {
			log.SetFlags(log.Ltime)
		}

		for {
			time.Sleep(time.Second * 1)
			// Scan line until \n
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					log.Fatalln("Connection closed by remote host")
					os.Exit(1)
				}
				log.Fatalln("Error reading data:", err.Error())

			}

			log.Print(data)

		}

	}()

	// Send data
	func() {
		log.SetFlags(0)
		if input.Time {
			log.SetFlags(log.Ltime)
		}

		for {
			time.Sleep(time.Second * 1)
			reader := bufio.NewReader(os.Stdin)

			// attach username
			text := input.Username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			conn.Write([]byte(text))
		}
	}()

}

// Listener for multiple ports
func Listen(input utils.Input) {
	var message = [][]string{}
	// Set up listener for every port in range
	for _, port := range input.Port {

		// wg = WaitGroup (Variable to wait until variable hits 0)
		wg.Add(1)
		go lnPort(port, input.Username, len(input.Port), &message, input.Time)
	}

	// Wait till wg is 0
	wg.Wait()
}

// Listener for individual port
func lnPort(port string, username string, nPorts int, message *[][]string, timestamp bool) {
	log.SetFlags(0)
	if timestamp {
		log.SetFlags(log.Ltime)
	}

	//Listen and connect
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		log.Fatalln("Permission denied.\nTry again with root or take a port above 1023")
		wg.Done()
		return
	}

	log.Println("Listening on port", port)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection:", err.Error())
		wg.Done()
		return
	}

	log.Printf("Remote Addr: %v\nLocal Addr: %v\n", conn.RemoteAddr(), conn.LocalAddr())

	// Read data
	go func() {
		for {
			time.Sleep(time.Second * 1)
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					log.Fatalf("Connection on port %v closed by remote host", port)
					wg.Done()
					return
				}
				log.Fatalln("Error reading data:", err.Error())

			}
			log.Print(data)

			if nPorts > 1 {
				for i := 0; i < nPorts-1; i++ {
					*message = append(*message, []string{utils.FilterIp(conn.RemoteAddr().String()), data})
				}
			}
		}

	}()

	// Send data from input
	go func() {
		for {
			time.Sleep(time.Second * 1)
			reader := bufio.NewReader(os.Stdin)

			// attach username
			text := username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			conn.Write([]byte(text))
		}
	}()

	//send data from other clients
	if nPorts > 1 {
		go func() {
			for {
				time.Sleep(time.Second * 1)
				if len(*message) > 0 {
					for index, element := range *message {
						if element[0] != port {
							log.Println(element[1])
							utils.Remove(*message, index)
						}
					}
				}
			}
		}()
	}
}
