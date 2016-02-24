package main

import "goqt/ui"

func main() {
	ui.Run(func() {
		widget := ui.NewWidget()
		widget.Show()
	})
}
