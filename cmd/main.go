package main

import (
	"os"

	"github.com/adzsx/g-wire/pkg/netcli"
	"github.com/adzsx/g-wire/pkg/utils"
)

func main() {
	args := os.Args

	input := utils.Format(args)

	if input.Action == "listen" {
		netcli.Listen(input)
	} else if input.Action == "connect" {
		netcli.Connect(input)
	}
}
