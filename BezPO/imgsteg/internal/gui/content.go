package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func (a *App) CreateContent() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Encoding", a.createEncodingContent()),
		container.NewTabItem("Decoding", a.createDecodingContent()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	settings := a.createSettingsPanel()

	split := container.NewHSplit(
		container.NewVScroll(tabs),
		container.NewVScroll(settings),
	)
	split.Offset = 0.75
	return split
}
