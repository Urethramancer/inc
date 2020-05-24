package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"io"

	"github.com/andybalholm/brotli"
)

// Compress byte slice.
func Compress(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, errors.New("no input data")
	}

	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	size, err := w.Write(content)
	if err != nil {
		return nil, err
	}

	if size == 0 {
		return nil, errors.New("zero size compression output")
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// CompressBrotli for better compression than gzip.
func CompressBrotli(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, errors.New("no input data")
	}

	opt := brotli.WriterOptions{
		Quality: 5,
	}
	b := bytes.Buffer{}
	w := brotli.NewWriterOptions(&b, opt)
	if w == nil {
		return nil, errors.New("couldn't allocate writer")
	}

	defer w.Close()
	in := bytes.NewReader(content)
	_, err := io.Copy(w, in)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
