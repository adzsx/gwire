package main

import (
	"log"
	"os"

	"github.com/adzsx/g-wire/pkg/netcli"
	"github.com/adzsx/g-wire/pkg/utils"
	"github.com/adzsx/g-wire/test"
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

	input := utils.Format(args)

	if len(args) < 3 && input.Action != "test" && input.Action != "help" {
		log.Fatalln("Enter --help for help")
	}

	if input.Action == "help" {
		log.Print(help)
		os.Exit(0)
	}

	if input.Action == "listen" {
		netcli.Listen(input)
	} else if input.Action == "connect" {
		netcli.Connect(input)
	} else if input.Action == "test" {
		test.Test()
	}

}
