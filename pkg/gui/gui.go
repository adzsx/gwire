package gui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	chatW    fyne.Window
	messages *widget.List
	text     []string
)

func chat(version string) {
	chatW = GApp.NewWindow(version)

	chatW.Resize(fyne.NewSize(1066, 600))

	host := widget.NewLabel(version)
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
		AddMsg(entry.Text)
	})

	content := container.NewBorder(
		container.NewHBox(container.NewCenter(quit), container.NewCenter(host)),

		container.NewGridWithColumns(2, entry, container.NewHBox(send)),

		nil,
		nil,

		container.NewStack(messages),
	)
	chatW.SetContent(content)
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

func AddLog(msg string) {

}
