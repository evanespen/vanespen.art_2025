package albums

import "github.com/google/uuid"

type Album struct {
	UUID        string   `json:"uuid" parquet:"name=uuid, type=FIXED_LEN_BYTE_ARRAY, length=36, convertedtype=UTF8, encoding=DELTA_BYTE_ARRAY"`
	Title       string   `json:"title" parquet:"name=title, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Description string   `json:"description" parquet:"name=description, type=BYTE_ARRAY, encoding=DELTA_LENGTH_BYTE_ARRAY"`
	Pictures    []string `json:"pictures" parquet:"name=pictures, type=MAP convertedtype=LIST, valuetype=FIXED_LEN_BYTE_ARRAY, valueconvertedtype=UTF8, valuelength=36"`
}

func NewAlbum() *Album {
	return &Album{
		UUID:     uuid.New().String(),
		Pictures: make([]string, 0), // Initialise un slice de chaînes vide, non-nil.
	}
}
