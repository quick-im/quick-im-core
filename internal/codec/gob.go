package codec

import (
	"bytes"
	"encoding/gob"
)

type GobUtils[T any] struct{}

func (gu *GobUtils[T]) Encode(data T) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (gu *GobUtils[T]) Decode(data []byte, result *T) error {
	buf := bytes.NewReader(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}
