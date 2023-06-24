package utils

import (
	"log"
	"net"
	"strconv"
	"strings"
)

type Input struct {
	Action   string
	Ip       string
	Port     []string
	Username string
	Time     bool
	TimeOut  float64
}

func Format(args []string) Input {

	input := Input{}

	input.Time = false

	for index, element := range args[0:] {
		switch element[1:] {
		case "l":
			input.Action = "listen"
		case "h":

			if len(args) < index+2 {
				log.Fatalln("Error: Host not defined")
			} else if net.ParseIP(args[index+1]) == nil {
				log.Fatalln("Error: Host addres invalid")
			} else {
				input.Ip = args[index+1]
			}
		case "p":

			if len(args) < index+2 {
				log.Fatalln("Error: Port not defined")
			} else if _, err := strconv.Atoi(args[index+1]); err != nil {
				log.Fatalln("Error: Port number invalid")
			} else {

				if InSlice(args, "-l") {
					for i := 1; len(args) > index+i && !strings.Contains(args[index+i], "-"); i++ {
						input.Port = append(input.Port, string(args[index+i]))
					}
				} else {
					input.Port = []string{args[index+1]}
				}
			}

		case "u":
			input.Username = args[index+1]
		case "t":
			input.Time = true
		case "s":

			if len(args) < index+2 {
				log.Fatalln("Error: Slowmode value not defined")
			}

			num, err := strconv.ParseFloat(args[index+1], 64)

			if err != nil {
				log.Fatalln("Error: Slowmode value invalid")
			}
			input.TimeOut = num * 1000

			log.Println(input.TimeOut)
		}
	}

	if !InSlice(args, "-l") {
		input.Action = "connect"
	}

	if input.Username == "" {
		input.Username = "anonymous"
	}

	if InSlice(args, "--test") {
		input.Action = "test"
	}

	if InSlice(args, "--help") {
		input.Action = "help"
	}

	return input
}
