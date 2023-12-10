package utils

import (
	"errors"
	"net"
	"strings"
)

/*
type Input struct {
	Action   string
	Ip       string
	Port     []string
	Username string
	Time     bool
	TimeOut  float64
	Enc      string
}
*/

func CheckInput(input Input) error {
	var missing []string
	if input.Action == "" {
		missing = append(missing, "action")
	}

	if input.Action == "info" || input.Action == "help" {
		return nil

	}

	if input.Ip == "" && input.Action != "listen" {
		missing = append(missing, "host")
	}

	if len(input.Port) == 0 {
		missing = append(missing, "port")
	}

	if len(missing) == 1 {
		return errors.New("missing value for: " + missing[0])
	} else if len(missing) > 1 {
		return errors.New("missing values for: " + strings.Join(missing, ", "))
	}

	return nil
}

func IP(ip string) bool {
	return net.ParseIP(ip) != nil
}
