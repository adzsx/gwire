package gui

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	chatW    fyne.Window
	messages *widget.List
	con      *widget.Label
	text     []string
)

func chat(version string) {
	chatW = GApp.NewWindow(version)

	chatW.Resize(fyne.NewSize(1066, 600))

	if len(input.Port) > 1 {
		con = widget.NewLabel(strings.Join(input.Port, ", "))
	} else if input.Port != nil {
		con = widget.NewLabel(input.Ip + ":" + input.Port[0])
	} else {
		con = widget.NewLabel("Not Connected")
	}

	quit := widget.NewButton("Quit", func() {
		os.Exit(0)
	})

	// Create the listbox

	messages = widget.NewList(
		func() int {
			return len(text)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(text[id])
			item.(*widget.Label).TextStyle.Monospace = false
		},
	)

	entry := widget.NewEntry()

	send := widget.NewButton("Send", func() {
		//SendMsg(entry.Text)
		AddMsg(entry.Text)
	})

	content := container.NewBorder(
		container.NewHBox(container.NewCenter(quit), container.NewCenter(con)),

		container.NewGridWithColumns(2, entry, container.NewHBox(send)),

		nil,
		nil,

		container.NewStack(messages),
	)
	chatW.SetContent(content)
}

func connected() {
	if len(input.Port) > 1 {
		con.SetText(strings.Join(input.Port, ", "))
	} else if input.Port != nil {
		con.SetText(input.Ip + ":" + input.Port[0])
	}
}

func AddMsg(msg string) {
	if msg == "" {
		return
	}
	text = append(text, msg)
	messages.Refresh()

	if messages.Length() > 0 {
		messages.ScrollToBottom()
	}
}
