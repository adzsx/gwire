package gui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/adzsx/gwire/pkg/netcli"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	a       fyne.App
	listen  *widget.Check
	encrypt *widget.Check
)

func Setup() utils.Input {

	var input utils.Input

	a = app.New()
	setup := a.NewWindow("Gwire Setup")

	title := widget.NewLabel("GWire Setup")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP address")

	ipWarning := widget.NewLabel("")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")

	portWarning := widget.NewLabel("")

	ipWarning.Hide()
	portWarning.Hide()

	scan := widget.NewCheck("Scan Network", func(b bool) {
		if b {
			listen.SetChecked(false)
			ipEntry.Disable()
			portEntry.Enable()
			input.Ip = "scan"
		} else {
			ipEntry.Enable()
		}
	})

	listen = widget.NewCheck("Listen", func(b bool) {
		listen.Refresh()
		if b {
			scan.SetChecked(false)
			ipEntry.Disable()
			input.Action = "listen"
			portEntry.SetPlaceHolder("Port(s) seperated by space")
		} else {
			ipEntry.Enable()
			portEntry.SetPlaceHolder("Port")
		}
	})

	password := widget.NewEntry()
	password.SetPlaceHolder("Password")

	pwWarning := widget.NewLabel("")

	autoEncrypt := widget.NewCheck("RSA", func(b bool) {
		if b {
			encrypt.SetChecked(true)
			password.Disable()

		}

		if !b && password.Disabled() {
			password.Enable()
		}
	})

	encrypt = widget.NewCheck("Encrypt", func(b bool) {
		if b {
			password.Enable()
		} else {
			autoEncrypt.SetChecked(false)
			password.Disable()
			input.Enc = ""
		}
	})

	autoEncrypt.SetChecked(true)

	info := widget.NewButton("Info", func() {
		infoWindow := a.NewWindow("Gwire Info")

		ip, mask, nHosts := netcli.Info()
		ipL := widget.NewLabel("Private IP:")
		ipLabel := widget.NewLabel(ip)

		maskL := widget.NewLabel("Subnetmask: ")
		maskLabel := widget.NewLabel(mask)

		hostL := widget.NewLabel("Possible Hosts: ")
		hostLabel := widget.NewLabel(nHosts)

		content := container.NewGridWithColumns(2,
			ipL, ipLabel,
			maskL, maskLabel,
			hostL, hostLabel,
		)

		infoWindow.SetContent(content)

		infoWindow.Show()
	})

	start := widget.NewButton("Go", func() {
		clear := true
		if encrypt.Checked {
			if autoEncrypt.Checked {
				input.Enc = "auto"
			} else if len(password.Text) != 32 {
				pwWarning.SetText("Password has to be 32 characters")
				clear = false

			} else {
				input.Enc = password.Text
			}
		}

		if userEntry.Text == "" {
			input.Username = "anonymous"
		} else {
			input.Username = userEntry.Text
		}

		input.Port = strings.Split(portEntry.Text, " ")

		if !listen.Checked && !scan.Checked {
			if utils.IP(ipEntry.Text) {
				input.Ip = ipEntry.Text
				ipWarning.Hide()
			} else {
				ipWarning.Show()
				ipWarning.SetText("Invalid IP")
				clear = false
			}
		}

		if portEntry.Text == "" {
			clear = false
			portWarning.Show()
			portWarning.SetText("Missing port")
		} else {
			portWarning.Hide()
		}

		port, err := strconv.Atoi(portEntry.Text)
		if err != nil {
			portWarning.Show()
			portWarning.SetText("Invalid port")
		} else if port < 1023 && os.Geteuid() != 0 {
			portWarning.Show()
			portWarning.SetText("Use root for port below 1023")
			clear = false
		} else if port > 65000 {
			portWarning.Show()
			portWarning.SetText("Use port below 65535")
			clear = false
		} else {
			portWarning.Hide()
		}

		input.TimeOut = 100

		if clear {
			fmt.Println(input)
			setup.Close()
		}
	})

	spacer := widget.NewLabel("")

	content := container.NewVBox(
		container.NewCenter(title),
		userEntry,
		spacer,

		container.NewGridWithColumns(2,
			listen, scan,
		),
		ipEntry,
		ipWarning,
		portEntry,
		portWarning,

		spacer,

		container.NewGridWithColumns(2, encrypt, autoEncrypt),
		password,
		pwWarning,

		container.NewGridWithColumns(2, info, start),
	)

	setup.SetContent(content)

	setup.ShowAndRun()

	return input
}
