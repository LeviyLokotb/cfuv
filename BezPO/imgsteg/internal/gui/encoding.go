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

func (a *App) createEncodingContent() fyne.CanvasObject {
	imageBaseView := canvas.NewImageFromResource(nil)
	imageBaseView.FillMode = canvas.ImageFillContain
	imageBaseView.SetMinSize(fyne.NewSize(300, 300))

	imageEncodedView := canvas.NewImageFromResource(nil)
	imageEncodedView.FillMode = canvas.ImageFillContain
	imageEncodedView.SetMinSize(fyne.NewSize(300, 300))

	loadButton := widget.NewButton("Load image...", func() {
		a.showLoadImageDialog(func(b img.Bitmap) {
			if err := a.switchBaseImage(b); err != nil {
				a.Errorf("switch base image error: %w", err)
				return
			}
			imageBaseView.Resource = a.toFyneResource(b)
			imageBaseView.Refresh()
			a.Notify("Загрузка", "Изображение загружено: %dx%d", b.W, b.H)
		})
	})

	leftPane := container.NewBorder(loadButton, nil, nil, nil, imageBaseView)
	rightPane := container.NewBorder(widget.NewLabel("Encoded"), nil, nil, nil, imageEncodedView)
	imagesRow := container.NewGridWithColumns(2, leftPane, rightPane)

	// --- динамическая input-зона ---
	inputHolder := container.NewMax()

	showEncoded := func() {
		enc := a.getEncodedBitmap()
		imageEncodedView.Resource = a.toFyneResource(enc)
		imageEncodedView.Refresh()
	}

	buildTextInput := func() fyne.CanvasObject {
		messageInput := widget.NewMultiLineEntry()
		messageInput.SetPlaceHolder("Message...")

		encodeButton := widget.NewButton("Encode", func() {
			if err := a.encodeText(messageInput.Text); err != nil {
				a.Errorf("encode error: %w", err)
				return
			}
			a.Notify("Encoded!", "Data encoded")
			showEncoded()
		})

		saveButton := widget.NewButton("Save As...", func() {
			a.showSaveImageDialog(func(w fyne.URIWriteCloser) {
				if err := a.saveEncodedTo(w); err != nil {
					a.Errorf("save error: %w", err)
				} else {
					a.Notify("Сохранение", "Файл сохранён")
				}
			})
		})

		buttons := container.NewHBox(encodeButton, saveButton)
		return container.NewBorder(nil, nil, nil, buttons, messageInput)
	}

	buildFileInput := func() fyne.CanvasObject {
		fileLabel := widget.NewLabel("Файл не выбран")

		loadFileButton := widget.NewButton("Load file...", func() {
			a.showLoadDataFileDialog(func(data []byte, name string) {
				a.setDataFile(data, name)
				fileLabel.SetText(fmt.Sprintf("%s (%d байт)", name, len(data)))
				a.Notify("Файл загружен", "%s, %d байт", name, len(data))
			})
		})

		encodeButton := widget.NewButton("Encode", func() {
			if err := a.encodeFile(); err != nil {
				a.Errorf("encode error: %w", err)
				return
			}
			a.Notify("Encoded!", "File encoded")
			showEncoded()
		})

		saveButton := widget.NewButton("Save As...", func() {
			a.showSaveImageDialog(func(w fyne.URIWriteCloser) {
				if err := a.saveEncodedTo(w); err != nil {
					a.Errorf("save error: %w", err)
				} else {
					a.Notify("Сохранение", "Файл сохранён")
				}
			})
		})

		buttons := container.NewHBox(encodeButton, saveButton)
		return container.NewBorder(
			nil,
			nil,
			nil,
			buttons,
			container.NewBorder(nil, nil, nil, loadFileButton, fileLabel),
		)
	}

	refreshInput := func() {
		var obj fyne.CanvasObject
		if a.contentType == steg.DataTypeFile {
			obj = buildFileInput()
		} else {
			obj = buildTextInput()
		}
		inputHolder.Objects = []fyne.CanvasObject{obj}
		inputHolder.Refresh()
	}

	refreshInput()

	// Не перезаписываем, а дополняем существующий колбэк.
	prev := a.onDataTypeChanged
	a.onDataTypeChanged = func() {
		if prev != nil {
			prev()
		}
		refreshInput()
	}

	return container.NewVScroll(
		container.NewVBox(imagesRow, inputHolder),
	)
}
