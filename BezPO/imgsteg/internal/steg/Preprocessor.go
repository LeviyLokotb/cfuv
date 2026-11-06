package steg

type Preprocessor interface {
	EncodeMeta(MetaData) []byte
	DecodeMeta([]byte) (MetaData, error)
	GetFileType() DataType
}
