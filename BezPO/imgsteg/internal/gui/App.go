package gui

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"golang.org/x/image/bmp"

	"github.com/LeviyLokotb/imgsteg/internal/img"
	"github.com/LeviyLokotb/imgsteg/internal/steg"
)

type App struct {
	window      fyne.Window
	imageBase   img.Bitmap
	options     steg.Options
	encoder     steg.Encoder
	decoder     steg.Decoder
	contentType steg.DataType

	// Загруженный файл (для DataTypeFile)
	fileData []byte
	fileName string

	onDataTypeChanged func()
}

func NewApp(w fyne.Window) App {
	a := App{
		window:  w,
		options: steg.NewOptionsOptimal(),
		// contentType выставится в switchDataType ниже
	}
	// Применяем дефолтный тип данных, чтобы Prep был проинициализирован.
	// Изображения ещё нет — switchOptions просто сохранит опции.
	if err := a.switchDataType(steg.DataTypeText); err != nil {
		log.Printf("init switchDataType error: %v", err)
	}
	return a
}

func (a *App) switchDataType(t steg.DataType) error {
	o := a.options
	switch t {
	case steg.DataTypeText:
		o.Prep = steg.TextPreprocessor{}
	case steg.DataTypeFile:
		o.Prep = steg.FilePreprocessor{}
	}
	if err := a.switchOptions(o); err != nil {
		return err
	}
	a.contentType = t
	if a.onDataTypeChanged != nil {
		a.onDataTypeChanged()
	}
	return nil
}

// setDataFile сохраняет загруженный файл в App.
func (a *App) setDataFile(data []byte, name string) {
	a.fileData = data
	a.fileName = name
}

func (a *App) switchOptions(options steg.Options) error {
	if err := options.Validate(); err != nil {
		return err
	}
	a.options = options
	// Если изображение уже загружено, пересоздаём кодировщики
	if a.imageBase.Px != nil {
		a.switchBaseImage(a.imageBase)
	}
	return nil
}

func (a *App) switchBaseImage(imageBase img.Bitmap) error {
	encoder, err := steg.NewEncoder(imageBase, a.options)
	if err != nil {
		return err
	}
	decoder, err := steg.NewDecoder(imageBase, a.options)
	if err != nil {
		return err
	}

	a.imageBase = imageBase
	a.encoder = encoder
	a.decoder = decoder
	return nil
}

// showLoadImageDialog открывает диалог выбора .bmp-изображения
// и передаёт результат в callback в виде img.Bitmap.
func (a *App) showLoadImageDialog(onLoad func(img.Bitmap)) {
	d := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
		if err != nil {
			a.Errorf("open file error: %w", err)
			return
		}
		if rc == nil {
			// Пользователь закрыл диалог без выбора.
			return
		}
		defer rc.Close()

		bitmap, err := a.decodeBitmap(rc)
		if err != nil {
			a.Errorf("decode image error: %w", err)
			return
		}
		onLoad(bitmap)
	}, a.window)

	d.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
	d.Show()
	d.Resize(fyne.NewSize(700, 500))
}

// showLoadDataFileDialog открывает диалог выбора произвольного файла
// и передаёт его содержимое и имя в callback.
func (a *App) showLoadDataFileDialog(onLoad func(data []byte, name string)) {
	d := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
		if err != nil {
			a.Errorf("open file error: %w", err)
			return
		}
		if rc == nil {
			return
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			a.Errorf("read file error: %w", err)
			return
		}
		onLoad(data, rc.URI().Name())
	}, a.window)

	d.Show()
	d.Resize(fyne.NewSize(700, 500))
}

// showSaveImageDialog открывает диалог сохранения .bmp-изображения
// и передаёт writer в callback.
func (a *App) showSaveImageDialog(onSave func(fyne.URIWriteCloser)) {
	d := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
		if err != nil {
			a.Errorf("save file error: %w", err)
			return
		}
		if wc == nil {
			return
		}
		// Закрывать writer — ответственность вызывающего кода,
		// потому что запись может быть не мгновенной.
		onSave(wc)
	}, a.window)

	d.SetFilter(storage.NewExtensionFileFilter([]string{".bmp"}))
	d.SetFileName("encoded.bmp")
	d.Show()
	d.Resize(fyne.NewSize(700, 500))
}

// showSaveDecodedFileDialog открывает диалог сохранения файла
// с предложенным именем и передаёт writer в callback.
func (a *App) showSaveDecodedFileDialog(defaultName string, onSave func(fyne.URIWriteCloser)) {
	d := dialog.NewFileSave(func(wc fyne.URIWriteCloser, err error) {
		if err != nil {
			a.Errorf("save file error: %w", err)
			return
		}
		if wc == nil {
			return
		}
		onSave(wc)
	}, a.window)

	d.SetFileName(defaultName)
	d.Show()
	d.Resize(fyne.NewSize(700, 500))
}

// encodeText вызывает кодирование текста в текущем изображении.
func (a *App) encodeText(text string) error {
	if a.imageBase.Px == nil {
		return fmt.Errorf("base image is not loaded")
	}
	err := a.encoder.EncodeBytes(steg.MetaData{
		Data: []byte(text),
	})
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

// encodeFile кодирует загруженный файл в текущем изображении.
// Имя файла сохраняется в MetaData.Meta.
func (a *App) encodeFile() error {
	if a.imageBase.Px == nil {
		return fmt.Errorf("base image is not loaded")
	}
	if a.fileData == nil {
		return fmt.Errorf("file is not loaded")
	}
	err := a.encoder.EncodeBytes(steg.MetaData{
		Data: a.fileData,
		Meta: a.fileName,
	})
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}

// decodeData создаёт декодер по базовому изображению и извлекает
// сообщение из закодированного изображения.
func (a *App) decodeData(baseBmp, encBmp img.Bitmap) (steg.MetaData, error) {
	if baseBmp.Px == nil {
		return steg.MetaData{}, fmt.Errorf("base image is not loaded")
	}
	if encBmp.Px == nil {
		return steg.MetaData{}, fmt.Errorf("encoded image is not loaded")
	}

	dec, err := steg.NewDecoder(baseBmp, a.options)
	if err != nil {
		return steg.MetaData{}, fmt.Errorf("new decoder: %w", err)
	}

	raw, err := dec.DecodeBytes(encBmp)
	if err != nil {
		return steg.MetaData{}, fmt.Errorf("decode: %w", err)
	}
	return raw, nil
}

// getEncodedBitmap возвращает текущее закодированное изображение.
func (a *App) getEncodedBitmap() img.Bitmap {
	return a.encoder.GetEncoded()
}

// saveEncodedTo сохраняет закодированное изображение в writer.
func (a *App) saveEncodedTo(w fyne.URIWriteCloser) error {
	bitmap := a.getEncodedBitmap()
	if bitmap.Px == nil {
		return fmt.Errorf("nothing to save: no encoded image")
	}
	if err := a.encodeBitmap(w, bitmap); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}
	return nil
}

// saveDecodedFile сохраняет расшифрованные байты в writer.
func (a *App) saveDecodedFile(w fyne.URIWriteCloser, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close writer: %w", err)
	}
	return nil
}

// --- Конвертация ---

// decodeBitmap читает BMP из rc и возвращает img.Bitmap.
func (a *App) decodeBitmap(rc io.Reader) (img.Bitmap, error) {
	src, err := bmp.Decode(rc)
	if err != nil {
		return img.Bitmap{}, fmt.Errorf("bmp decode: %w", err)
	}
	return img.FromImage(src), nil
}

// encodeBitmap записывает img.Bitmap в w в формате BMP.
func (a *App) encodeBitmap(w io.Writer, b img.Bitmap) error {
	if err := bmp.Encode(w, b.ToImage()); err != nil {
		return fmt.Errorf("bmp encode: %w", err)
	}
	return nil
}

// toFyneResource конвертирует img.Bitmap в fyne.Resource.
// Реализовано через PNG-байты, чтобы результат подходил всюду,
// где Fyne ожидает именно fyne.Resource (иконки, canvas.Image и т.д.).
func (a *App) toFyneResource(b img.Bitmap) fyne.Resource {
	if b.Px == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, b.ToImage()); err != nil {
		a.Errorf("png encode: %w", err)
		return nil
	}
	return fyne.NewStaticResource("bitmap.png", buf.Bytes())
}

func (app App) Errorf(format string, a ...any) {
	msg := fmt.Errorf(format, a...).Error()
	log.Println(msg)
	fyne.CurrentApp().SendNotification(fyne.NewNotification("Ошибка", msg))
}

func (app App) Notify(head string, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fyne.CurrentApp().SendNotification(fyne.NewNotification(head, msg))
}
