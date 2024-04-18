package netcli

import (
	"fmt"

	"github.com/adzsx/gwire/internal/utils"
)

func AddMsg(msg string, gui bool) {
	if !gui {
		fmt.Println()
		utils.Ansi("\033[2A\033[999D\033[K\033[L")

		color := utils.GetRandomString(colorList, utils.FilterChar(msg, ">", true))
		utils.Ansi(color)
		fmt.Print(msg)
		utils.Ansi("\033[999B\033[999D\033[1C")
	} else {
		rcv <- msg
	}

}
