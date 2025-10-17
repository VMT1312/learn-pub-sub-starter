package pubsub

import (
	"bytes"
	"encoding/gob"
)

func encode(T any) ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(T)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
