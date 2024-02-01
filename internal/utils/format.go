package utils

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	verbose int
)

type Input struct {
	Action   string
	Ip       string
	Port     []string
	Username string
	Time     bool
	TimeOut  float64
	Enc      string
}

func Format(args []string) Input {

	if len(args) < 2 {
		log.Println("Enter --help for help")
		os.Exit(0)
	}

	input := Input{}
	input.Enc = "auto"
	input.TimeOut = 1000

	for index, element := range args[0:] {
		switch element[1:] {
		case "nfo":
			input.Action = "info"
			return input
		case "l", "-listen":
			input.Action = "listen"

		case "h", "-host":
			if len(args) > index+1 && net.ParseIP(args[index+1]) != nil {
				input.Ip = args[index+1]
			} else {
				input.Ip = "scan"
			}

		case "p", "-port":
			if len(args) < index+2 {
				Err(errors.New("port not defined"), true)
			} else if _, err := strconv.Atoi(args[index+1]); err != nil {
				Err(errors.New("invalid port"), true)
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
			if len(args) < index+2 || args[index+1][0:1] == "-" {
				input.Enc = "auto"
			} else if len(args[index+1]) != 32 {
				Err(errors.New("password has to be 32 characters"), true)
			} else if len(args[index+1]) == 32 {
				input.Enc = args[index+1]
			}
		case "d", "-no-encryption":
			input.Enc = ""
			Print("No encryption", 2)

		case "u", "-username":
			input.Username = args[index+1]

		case "t", "-time":
			input.Time = true

		case "s", "-slowmode":
			if InSlice(args, "-l") {
				if len(args) < index+2 {
					Err(errors.New("slowmode value not defined"), true)
				}

				num, err := strconv.ParseFloat(args[index+1], 64)

				if err != nil {
					Err(errors.New("slowmode value invalid"), true)
				}

				input.TimeOut = num * 1000
			}

		case "v", "-verbose":
			if len(args) < index+2 || args[index+1][:1] == "-" {
				verbose = 1
			} else {
				var err error
				verbose, err = strconv.Atoi(args[index+1])
				Err(err, false)

				if verbose < 3 {
					verbose = 3
				}

			}

		case "-debug":
			verbose = 3
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

	if input.TimeOut < 100 {
		input.TimeOut = 100
	}
	return input
}
