package steg

import "fmt"

type FilePreprocessor struct{}

func (t FilePreprocessor) EncodeMeta(d MetaData) []byte {
	filename := ([]byte)(d.Meta)
	lname := len(filename)
	ldata := len(d.Data)
	pref := withNumPrefix(filename, uint32(ldata))
	pref = withNumPrefix(pref, uint32(lname))
	return append(pref, d.Data...)
}

func (t FilePreprocessor) DecodeMeta(data []byte) (MetaData, error) {
	// Достаём префикс (длину данных) и обрезаем
	lname := btoi(data[:4])
	ldata := btoi(data[4:8])

	if ldata < 0 || lname < 0 {
		return MetaData{}, fmt.Errorf("invalid length: %d+%d", lname, ldata)
	}
	if ldata+lname+4 > len(data) {
		return MetaData{}, fmt.Errorf(
			"message length %d exceeds buffer %d (wrong options or wrong image?)",
			ldata+lname+4, len(data),
		)
	}

	return MetaData{
		Data: data[8+lname : 8+lname+ldata],
		Meta: string(data[8 : 8+lname]),
	}, nil
}

func (t FilePreprocessor) GetFileType() DataType {
	return DataTypeText
}
