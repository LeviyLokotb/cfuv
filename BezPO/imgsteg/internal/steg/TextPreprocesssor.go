package steg

import "fmt"

type TextPreprocessor struct{}

func (t TextPreprocessor) EncodeMeta(d MetaData) []byte {
	return withNumPrefix(d.Data, uint32(len(d.Data)))
}

func (t TextPreprocessor) DecodeMeta(data []byte) (MetaData, error) {
	// Достаём префикс (длину данных) и обрезаем
	dataLen := btoi(data[:4])

	if dataLen < 0 {
		return MetaData{}, fmt.Errorf("invalid message length: %d", dataLen)
	}
	if dataLen+4 > len(data) {
		return MetaData{}, fmt.Errorf(
			"message length %d exceeds buffer %d (wrong options or wrong image?)",
			dataLen, len(data),
		)
	}

	return MetaData{
		Data: data[4 : 4+dataLen],
		Meta: "",
	}, nil
}

func (t TextPreprocessor) GetFileType() DataType {
	return DataTypeText
}
