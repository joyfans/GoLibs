package main

import (
	"os"

	"goqt/ui"
)

func main() {
	ui.RunEx(os.Args, func() {
		ed := NewCodeEdit()
		ed.ResizeWithWidthHeight(400, 400)
		ed.Show()
	})
}
