package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/LeviyLokotb/imgsteg/internal/img"
	"github.com/LeviyLokotb/imgsteg/internal/steg"
)

func (a *App) createDecodingContent() fyne.CanvasObject {
	imageBaseView := canvas.NewImageFromResource(nil)
	imageBaseView.FillMode = canvas.ImageFillContain
	imageBaseView.SetMinSize(fyne.NewSize(300, 300))

	imageEncodedView := canvas.NewImageFromResource(nil)
	imageEncodedView.FillMode = canvas.ImageFillContain
	imageEncodedView.SetMinSize(fyne.NewSize(300, 300))

	// Encoded-изображение локально: пользователь может подгрузить
	// другой файл, не трогая общий imageBase.
	var encBmp img.Bitmap

	loadBaseButton := widget.NewButton("Load base...", func() {
		a.showLoadImageDialog(func(b img.Bitmap) {
			if err := a.switchBaseImage(b); err != nil {
				a.Errorf("switch base image error: %w", err)
				return
			}
			imageBaseView.Resource = a.toFyneResource(b)
			imageBaseView.Refresh()
			a.Notify("Base image", "Загружено: %dx%d", b.W, b.H)
		})
	})

	loadEncButton := widget.NewButton("Load encoded...", func() {
		a.showLoadImageDialog(func(b img.Bitmap) {
			encBmp = b
			imageEncodedView.Resource = a.toFyneResource(b)
			imageEncodedView.Refresh()
			a.Notify("Encoded image", "Загружено: %dx%d", b.W, b.H)
		})
	})

	leftPane := container.NewBorder(loadBaseButton, nil, nil, nil, imageBaseView)
	rightPane := container.NewBorder(loadEncButton, nil, nil, nil, imageEncodedView)
	imagesRow := container.NewGridWithColumns(2, leftPane, rightPane)

	// --- динамическая output-зона ---
	outputHolder := container.NewMax()

	// Результат последнего успешного decode.
	var lastDecoded []byte
	var lastMeta string
	hasDecoded := false

	resetResult := func() {
		lastDecoded = nil
		lastMeta = ""
		hasDecoded = false
	}

	// buildOutput собирает output-зону под текущий contentType.
	buildOutput := func() fyne.CanvasObject {
		if a.contentType == steg.DataTypeFile {
			// File: label с именем файла + кнопка сохранения,
			// которая активируется только после успешного decode.
			nameLabel := widget.NewLabel("Файл не расшифрован")
			nameLabel.Wrapping = fyne.TextWrapWord

			saveButton := widget.NewButton("Save decoded...", nil)
			saveButton.Disable()
			saveButton.OnTapped = func() {
				if !hasDecoded {
					a.Errorf("nothing to save: decode first")
					return
				}
				name := lastMeta
				if name == "" {
					name = "decoded.bin"
				}
				a.showSaveDecodedFileDialog(name, func(w fyne.URIWriteCloser) {
					if err := a.saveDecodedFile(w, lastDecoded); err != nil {
						a.Errorf("save decoded error: %w", err)
					} else {
						a.Notify("Сохранение", "Файл сохранён")
					}
				})
			}

			decodeButton := widget.NewButton("Decode", func() {
				resetResult()
				nameLabel.SetText("Файл не расшифрован")
				saveButton.Disable()

				meta, err := a.decodeData(a.imageBase, encBmp)
				if err != nil {
					a.Errorf("decode error: %w", err)
					return
				}
				lastDecoded = meta.Data
				lastMeta = meta.Meta
				hasDecoded = true

				a.Notify("Decoded!", "Data decoded")
				nameLabel.SetText(fmt.Sprintf(
					"Файл: %s\nРазмер: %d байт",
					meta.Meta, len(meta.Data),
				))
				saveButton.Enable()
			})

			buttons := container.NewHBox(decodeButton, saveButton)
			return container.NewBorder(nil, nil, nil, buttons, nameLabel)
		}

		// Text: entry с текстом + кнопки decode/save.
		messageOutput := widget.NewMultiLineEntry()
		messageOutput.SetPlaceHolder("Decoded message...")

		saveButton := widget.NewButton("Save decoded...", nil)
		saveButton.Disable()
		saveButton.OnTapped = func() {
			if !hasDecoded {
				a.Errorf("nothing to save: decode first")
				return
			}
			name := lastMeta
			if name == "" {
				name = "decoded.txt"
			}
			a.showSaveDecodedFileDialog(name, func(w fyne.URIWriteCloser) {
				if err := a.saveDecodedFile(w, lastDecoded); err != nil {
					a.Errorf("save decoded error: %w", err)
				} else {
					a.Notify("Сохранение", "Файл сохранён")
				}
			})
		}

		decodeButton := widget.NewButton("Decode", func() {
			resetResult()
			messageOutput.SetText("")
			saveButton.Disable()

			meta, err := a.decodeData(a.imageBase, encBmp)
			if err != nil {
				a.Errorf("decode error: %w", err)
				return
			}
			lastDecoded = meta.Data
			lastMeta = meta.Meta
			hasDecoded = true

			a.Notify("Decoded!", "Data decoded")
			messageOutput.SetText(string(meta.Data))
			saveButton.Enable()
		})

		buttons := container.NewHBox(decodeButton, saveButton)
		return container.NewBorder(nil, nil, nil, buttons, messageOutput)
	}

	refreshOutput := func() {
		resetResult()
		obj := buildOutput()
		outputHolder.Objects = []fyne.CanvasObject{obj}
		outputHolder.Refresh()
	}

	// Первичная отрисовка + дополнение общего колбэка смены типа.
	refreshOutput()
	prev := a.onDataTypeChanged
	a.onDataTypeChanged = func() {
		if prev != nil {
			prev()
		}
		refreshOutput()
	}

	return container.NewVScroll(
		container.NewVBox(imagesRow, outputHolder),
	)
}
