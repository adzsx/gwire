package utils

import (
	"errors"
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

	if input.Ip == "" && input.Action != "listen" {
		missing = append(missing, "host")
	}

	if len(input.Port) == 0 {
		missing = append(missing, "port")
	}

	if len(missing) == 1 {
		return errors.New("missing value for: " + strings.Join(missing, ", "))
	} else if len(missing) > 1 {
		return errors.New("missing values for: " + strings.Join(missing, ", "))
	}

	return nil
}
