package main

import (
	"log"
	"os"

	"github.com/adzsx/g-wire/pkg/netcli"
	"github.com/adzsx/g-wire/pkg/utils"
)

func main() {
	args := os.Args

	input, err := utils.Format(args)

	log.SetFlags(0)

	if err != nil {
		log.Printf("Input Error: %v", err)
		os.Exit(1)
	}

	if input.Action == "listen" {
		netcli.Listen(input)
	} else if input.Action == "connect" {
		netcli.Connect(input)
	}

}
