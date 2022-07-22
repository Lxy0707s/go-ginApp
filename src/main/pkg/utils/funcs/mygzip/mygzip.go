package mygzip

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

// Compress compress data to gzip format
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("write err: %v", err)
	}
	err = w.Flush()
	if err != nil {
		return nil, fmt.Errorf("flush err: %v", err)
	}
	err = w.Close() //cannot be defined as defer w.Close, in case of gzip read before close, which will cause err "unexpected EOF"
	if err != nil {
		return nil, fmt.Errorf("close err: %v", err)
	}
	return b.Bytes(), nil
}

// DeCompress decompress gzip data
func DeCompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	return undatas, err
}
