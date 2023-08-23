package netcli

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/adzsx/gwire/pkg/utils"
)

var (
	sconn   net.Conn
	counter int
	found   bool
)

func ScanRange(ips []string, port string) (string, net.Conn) {
	connChan := make(chan net.Conn)

	for _, element := range ips {
		address := element + ":" + port
		counter++
		go scan(address, connChan)
	}

	for counter > 0 {
		time.Sleep(time.Millisecond * 100)
		if counter == 0 {
			break
		}
	}

	if len(connChan) == 0 {
		log.Println("No host was found")
		os.Exit(0)
	}

	sconn = <-connChan

	ip := sconn.RemoteAddr().String()

	log.Println("Done scanning\nFound: " + sconn.RemoteAddr().String())

	return ip, sconn
}

func scan(address string, connChan chan net.Conn) {
	ping := Ping(utils.FilterChar(address, ":", true))
	if !ping {
		time.Sleep(time.Millisecond)
		counter--
		return
	}
	conn, err := net.Dial("tcp", address)

	if err != nil {
		/* log.Println(counter) */
		time.Sleep(time.Millisecond)
		counter--
		return
	}

	for found {
		if !found {
			break
		}
		time.Sleep(time.Second)
	}

	for {
		if !found {
			found = true
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Found open port on %v\nDo you want to connect? [y/n] ", utils.FilterChar(address, ":", true))
			input, _ := reader.ReadString('\n')

			input = input[0:1]

			if input == "n" || input == "no" {

				counter--
				if counter == 0 {
					return
				}
				reader = bufio.NewReader(os.Stdin)
				fmt.Print("Continue scan anyways? [y/n] ")
				input, _ := reader.ReadString('\n')

				input = input[0:1]

				log.Println()

				found = false

				if input == "n" || input == "no" {
					conn.Close()
					os.Exit(0)
				} else {
					return

				}
			}

			log.Println("Accepted")

			counter = 0

			connChan <- conn
			log.Println("Done")
		}

		time.Sleep(time.Second)

	}

}

func Ping(ip string) bool {

	cmd := exec.Command("ping", "-i", "0.1", "-c", "3", "-w", "1", ip)
	out, _ := cmd.Output()

	output, _ := strconv.Atoi(utils.FilterChar(utils.FilterChar(string(out), ",", false), ",", true)[1:2])

	return output > 1

}
