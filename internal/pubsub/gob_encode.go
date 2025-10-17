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

func decode[T any](data []byte) (T, error) {
	r := bytes.NewReader(data)
	decoder := gob.NewDecoder(r)
	var target T
	err := decoder.Decode(&target)
	if err != nil {
		var result T
		return result, err
	}
	return target, nil
}
