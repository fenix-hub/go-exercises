package main

import (
	"golang.org/x/tour/reader"
)
	
type MyReader struct{}

// TODO: Add a Read([]byte) (int, error) method to MyReader.

func (r MyReader) Read(b []byte) (int, error) {
	l := len(b)
	for i := range l {
		b[i] = byte('A')
	}
	return l, nil
}

func main() {
	reader.Validate(MyReader{})
}
