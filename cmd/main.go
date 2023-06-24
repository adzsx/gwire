package main

import (
	"log"
	"os"

	"github.com/adzsx/g-wire/pkg/netcli"
	"github.com/adzsx/g-wire/pkg/utils"
)

var (
	help = `
gwire usage:
	gwire [flags]

Flags:
	-h 			help message
	-l			listen
	-p port			Open connection on [port]
	-h host			Connect to [host] Ip
	-u username		[username] is didsplayed for other users
	-t			enable timestamps
	`
)

func main() {
	log.SetFlags(0)

	args := os.Args

	if utils.InSlice(args, "--help") {
		log.Print(help)
		os.Exit(0)
	}

	if len(args) < 3 {
		log.Fatalln("Enter --help for help")
	}

	input, err := utils.Format(args)

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
