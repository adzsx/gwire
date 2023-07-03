package main

import (
	"log"
	"os"

	"github.com/adzsx/gwire/pkg/netcli"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	help = `
gwire usage:
	gwire [flags]

Flags:
	    --help			help message
	-l, --listen			listen
	-p, --port 	[port]		use port [port]
	-h, --host 	[host]		Connect to [host]-(Ip)
	-u, --username 	[username]	[username] is didsplayed for other users
	-t, --time			enable timestamps
	-s, --slowmode	[seconds]	Enable slowmode
	-e, --encrypt	[password]	Encrypt messages with AES
	`

	version = "gwire v1.1"
)

func main() {
	log.SetFlags(0)

	args := os.Args

	input := utils.Format(args)

	if input.Action == "version" {
		log.Println(version)
		os.Exit(0)
	}

	if len(args) < 3 && input.Action != "test" && input.Action != "help" {
		log.Println("Enter --help for help")
	}

	if input.Action == "help" {
		log.Print(help)
		os.Exit(0)
	}

	if input.Action == "listen" {
		netcli.Listen(input)
	} else if input.Action == "connect" {
		netcli.Connect(input)
	}

}
