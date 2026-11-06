package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/LeviyLokotb/imgsteg/internal/gui"
)

func main() {
	fyneApp := app.NewWithID("imgsteg.LeviyLokotb")
	window := fyneApp.NewWindow("image steganographia")

	myApp := gui.NewApp(window)
	content := myApp.CreateContent()

	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
