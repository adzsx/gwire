package main

import (
	"log"
	"os"

	"github.com/adzsx/gwire/pkg/gui"
	"github.com/adzsx/gwire/pkg/netcli"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	help = `
gwire usage:
	gwire [mode] [flags]

Modes:
	[No input]	Uses mode to connect and chat
	info		Shows network related info


Flags:
	    --help				help message
	-l, --listen				listen
	-p, --port 		[port]		use port [port]
	-h, --host 		[host]		Connect to [host]-(Ip)
	-v, --verbose		[level]		Show some more info levels: 1-3
	-u, --username 		[username]	[username] is didsplayed for other users
	-t, --time				enable timestamps
	-s, --slowmode		[seconds]	Enable slowmode
	-e, --encrypt		[password]	use AES instead of RSA
	-d, --no-encryption			Do not use encryption (Not recommended)
	`

	version = "gwire v1.4.0"

	input utils.Input

	graphics bool
)

func main() {
	log.SetFlags(0)
	args := os.Args

	if len(args) == 1 {
		gui.GUI(version)
	} else {
		input = utils.Format(args)

		if input.Action == "version" {
			log.Println(version)
			os.Exit(0)
		}

		err := utils.CheckInput(input)
		utils.Err(err, true)

		if input.Action == "help" {
			log.Print(help)
			os.Exit(0)
		}

		if input.Action == "listen" {
			netcli.HostSetup(input, graphics)

		} else if input.Action == "connect" {
			netcli.ClientSetup(input, graphics)
		} else if input.Action == "info" {

			ip, mask, nHosts := netcli.Info()
			log.Printf("Private IP: 		%v\nSubnetmask: 		%v\nNumber of hosts: 	%v", ip, mask, nHosts)

		}
	}

}
