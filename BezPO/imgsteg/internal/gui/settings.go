package gui

import (
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/LeviyLokotb/imgsteg/internal/img"
	"github.com/LeviyLokotb/imgsteg/internal/steg"
)

func (a *App) createSettingsPanel() fyne.CanvasObject {
	radioDataType := widget.NewRadioGroup(
		[]string{"Text", "File"},
		func(s string) {
			var err error
			switch s {
			case "Text":
				err = a.switchDataType(steg.DataTypeText)
			case "File":
				err = a.switchDataType(steg.DataTypeFile)
			}
			if err != nil {
				a.Errorf("switch options error: %w", err)
			}
		},
	)
	radioDataType.Horizontal = true
	radioDataType.Selected = "Text"
	radioDataType.Required = true

	entryBitsToMod := widget.NewEntry()
	entryBitsToMod.SetPlaceHolder("1-8")
	entryBitsToMod.OnSubmitted = func(s string) {
		n, err := strconv.Atoi(s)
		if err != nil {
			a.Errorf("input error: %w", err)
			return
		}
		o := a.options
		o.BitsToMod = n
		if err := a.switchOptions(o); err != nil {
			a.Errorf("switch options error: %w", err)
		}
	}

	checkChansToEncode := widget.NewCheckGroup(
		[]string{"R", "G", "B"},
		func(selected []string) {
			var chans []img.ColorChan
			if slices.Contains(selected, "R") {
				chans = append(chans, img.Red)
			}
			if slices.Contains(selected, "G") {
				chans = append(chans, img.Green)
			}
			if slices.Contains(selected, "B") {
				chans = append(chans, img.Blue)
			}

			o := a.options
			o.ChansToEncode = chans
			if err := a.switchOptions(o); err != nil {
				a.Errorf("input error: %w", err)
			}
		},
	)
	checkChansToEncode.Horizontal = true
	checkChansToEncode.Required = true

	// Синхронизирует поля ручного ввода с текущими опциями.
	syncManualFields := func(o steg.Options) {
		entryBitsToMod.SetText(strconv.Itoa(o.BitsToMod))

		var selected []string
		for _, c := range o.ChansToEncode {
			switch c {
			case img.Red:
				selected = append(selected, "R")
			case img.Green:
				selected = append(selected, "G")
			case img.Blue:
				selected = append(selected, "B")
			}
		}
		checkChansToEncode.SetSelected(selected)
	}

	// Применяет готовый конфиг.
	applyOptions := func(o steg.Options) {
		if err := a.switchOptions(o); err != nil {
			a.Errorf("switch options error: %w", err)
			return
		}
		syncManualFields(o)
	}

	const manualLabel = "Manual options"

	radioConfig := widget.NewRadioGroup(
		[]string{
			"Minimal",
			"Optimal",
			"Red chan",
			"Agressive",
			manualLabel,
		},
		func(selected string) {
			switch selected {
			case manualLabel:
				entryBitsToMod.Enable()
				checkChansToEncode.Enable()
				return
			case "Minimal":
				applyOptions(steg.NewOptionsMinimal())
			case "Optimal":
				applyOptions(steg.NewOptionsOptimal())
			case "Red chan":
				applyOptions(steg.NewOptionsRedChan())
			case "Agressive":
				applyOptions(steg.NewOptionsAgressive())
			}
			// Любой готовый конфиг блокирует ручной ввод.
			entryBitsToMod.Disable()
			checkChansToEncode.Disable()
		},
	)
	// radioConfig.Horizontal = true

	// Инициализация из текущих опций.
	syncManualFields(a.options)
	radioConfig.SetSelected(manualLabel)

	fields := container.NewVBox(
		container.NewVBox(
			widget.NewLabel("Data type:"),
			radioDataType,
		),
		container.NewVBox(
			widget.NewLabel("Config:"),
			radioConfig,
		),
		container.NewHBox(
			widget.NewLabel("Bits in the color channel for modification:"),
			entryBitsToMod,
		),
		container.NewHBox(
			widget.NewLabel("Rewritable color channels:"),
			checkChansToEncode,
		),
	)

	settingsPanel := container.NewVBox(
		fields,
	)

	return settingsPanel
}
