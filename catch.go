package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Try 捕获panic
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic: %v\n", err)
			if handler != nil {
				handler(err)
			}
		}
	}()
	fun()
}

// DeepCopy 深拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
