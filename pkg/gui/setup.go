package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	setup.Resize(fyne.NewSize(450, 450))

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("Username")

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("IP address")

	ipWarning := widget.NewLabel("")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")

	portWarning := widget.NewLabel("")

	scan := widget.NewCheck("Scan Network", func(b bool) {
		if b {
			listen.SetChecked(false)
			ipEntry.Disable()
			portEntry.Enable()
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

		}
	})

	autoEncrypt.SetChecked(true)

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
			} else {
				ipWarning.SetText("Invalid IP")
				clear = false
			}
		}

		if clear {
			fmt.Println(input)
			setup.Close()
		}
	})

	spacer := widget.NewLabel("")

	content := container.NewVBox(
		userEntry,
		spacer,

		container.NewGridWithColumns(2,
			listen, scan,
			ipEntry, portEntry,
			ipWarning, portWarning,
		),
		spacer,

		container.NewGridWithColumns(2, encrypt, autoEncrypt),
		password,
		pwWarning,

		start,
	)

	setup.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == "Return" {
			fmt.Println("Continue")
		}
	})

	setup.SetContent(content)

	setup.ShowAndRun()

	return input
}
