package utils

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

var (
	//errors
	invalid []string
)

type Input struct {
	Action   string
	Ip       string
	Port     []string
	Username string
}

func Format(args []string) (Input, error) {

	input := Input{}

	for index, element := range args[0:] {
		switch element[1:] {
		case "l":
			input.Action = "listen"
		case "h":

			if len(args) < index+2 {
				invalid = append(invalid, "host")
			} else if net.ParseIP(args[index+1]) == nil {
				invalid = append(invalid, "host")
			} else {
				input.Ip = args[index+1]
			}
		case "p":

			if len(args) < index+2 {
				invalid = append(invalid, "port")
			} else if _, err := strconv.Atoi(args[index+1]); err != nil {
				invalid = append(invalid, "port")
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
		}
	}

	if !InSlice(args, "-l") {
		input.Action = "connect"
	}

	if input.Username == "" {
		input.Username = "anonymous"
	}

	if input.Ip == "" && !InSlice(invalid, "host") {
		invalid = append(invalid, "host")
	}
	if len(input.Port) == 0 && !InSlice(invalid, "port") {
		invalid = append(invalid, "port")
	}

	if len(invalid) > 0 {
		return input, errors.New("Invalid arguments: " + strings.Join(invalid, ", "))
	}

	return input, nil
}
