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
	Key      []byte
}

func Format(args []string) Input {

	input := Input{}

	input.Time = false

	for index, element := range args[0:] {
		switch element[1:] {
		case "l", "-listen":
			input.Action = "listen"
		case "h", "-host":

			if len(args) < index+2 {
				log.Fatalln("Error: Host not defined")
			} else if net.ParseIP(args[index+1]) == nil {
				log.Fatalln("Error: Host addres invalid")
			} else {
				input.Ip = args[index+1]
			}
		case "p", "-port":

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

		case "e", "-encrypt":
			if len(args) < index+2 {
				log.Fatalln("Error: Key not defined")
			} else if len(args[index+1]) != 32 {
				log.Fatalf("Error: Key invalid key length: %v\nKey has to be 32 characters", len(args[index+1]))
			} else {
				input.Key = []byte(args[index+1])
			}

		case "u", "-username":
			input.Username = args[index+1]
		case "t", "-time":
			input.Time = true
		case "s", "-slowmode":

			if len(args) < index+2 {
				log.Fatalln("Error: Slowmode value not defined")
			}

			num, err := strconv.ParseFloat(args[index+1], 64)

			if err != nil {
				log.Fatalln("Error: Slowmode value invalid")
			}

			input.TimeOut = num * 1000
		}
	}

	if input.Action == "" {
		input.Action = "connect"
	}

	if input.Username == "" {
		input.Username = "anonymous"
	}

	if InSlice(args, "--help") {
		input.Action = "help"
	}

	if InSlice(args, "--version") {
		input.Action = "version"
	}

	return input
}
