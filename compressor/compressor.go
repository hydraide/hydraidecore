// Package compressor provides a simple interface to compress and decompress data using different algorithms
// A Hydra uses a compressor to store the Gems compressed in files and, after decompression,
// it is able to transform the Gems back into a readable format.
// The control of the compressor is done through the Traits, compressing and decompressing the data according to the
// algorithm specified in the setting. Therefore, the compressor is not accessible from the outside for the head of Hydra!
//
// The New function returns a new compressor. The compressorType parameter specifies the algorithm to be used.
package compressor

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"io"
)

type Type int

// available compression methods
const (
	Gzip Type = iota + 1
	LZ4
	Snappy
	Zstd
)

type Compressor interface {
	Compress(uncompressed []byte) (compressed []byte, err error)
	Decompress(compressed []byte) (decompressed []byte, err error)
}

type compressor struct {
	compressorType Type
}

// New returns a new compressor
func New(compressorType Type) Compressor {

	return &compressor{
		compressorType: compressorType,
	}
}

// Compress compresses the given data using the given algorithm
func (c *compressor) Compress(uncompressed []byte) (compressed []byte, err error) {
	switch c.compressorType {
	case Gzip:
		return c.compressGzip(uncompressed)
	case LZ4:
		return c.compressLZ4(uncompressed)
	case Snappy:
		return c.compressSnappy(uncompressed)
	case Zstd:
		return c.compressZstd(uncompressed)
	}
	return nil, errors.New("unknown compressor type")
}

// Decompress decompresses the given data using the given algorithm
func (c *compressor) Decompress(compressed []byte) (decompressed []byte, err error) {
	switch c.compressorType {
	case Gzip:
		return c.decompressGzip(compressed)
	case LZ4:
		return c.decompressLZ4(compressed)
	case Snappy:
		return c.decompressSnappy(compressed)
	case Zstd:
		return c.decompressZstd(compressed)
	}
	return nil, errors.New("unknown compressor type")
}

// compressGzip compress the given data using the gzip algorithm
func (c *compressor) compressGzip(uncompressed []byte) (compressed []byte, err error) {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err = gz.Write(uncompressed); err != nil {
		return nil, err
	}
	if err = gz.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

// decompressGzip decompress the given data using the gzip algorithm
func (c *compressor) decompressGzip(compressed []byte) (decompressed []byte, err error) {

	reader := bytes.NewReader(compressed)
	gzReader, e1 := gzip.NewReader(reader)

	if e1 != nil {
		return nil, err
	}

	output, e2 := io.ReadAll(gzReader)
	if e2 != nil {
		return nil, err
	}

	return output, nil

}

// compressLZ4 compress the given data using the lz4 algorithm
func (c *compressor) compressLZ4(uncompressed []byte) (compressed []byte, err error) {

	// create a new buffer to write the compressed data to
	w := &bytes.Buffer{}
	zw := lz4.NewWriter(w)
	_, err = zw.Write(uncompressed)
	if err != nil {
		return nil, err
	}
	// Closing is *very* important
	if err = zw.Close(); err != nil {
		return nil, err
	}
	return w.Bytes(), nil

}

// decompressLZ4 decompress the given data using the lz4 algorithm
func (c *compressor) decompressLZ4(compressed []byte) (decompressed []byte, err error) {

	r := bytes.NewReader(compressed)
	zr := lz4.NewReader(r)
	output, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return output, nil

}

// compressSnappy compress the given data using the snappy algorithm
func (c *compressor) compressSnappy(uncompressed []byte) (compressed []byte, err error) {
	compressed = snappy.Encode(nil, uncompressed)
	return compressed, nil
}

// decompressSnappy decompress the given data using the snappy algorithm
func (c *compressor) decompressSnappy(compressed []byte) (decompressed []byte, err error) {

	decompressed, err = snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}
	return decompressed, nil

}

// compressZstd compress the given data using the zstd algorithm
func (c *compressor) compressZstd(uncompressed []byte) (compressed []byte, err error) {

	w := &bytes.Buffer{}
	zstdEncoder, err := zstd.NewWriter(w)
	if err != nil {
		return nil, err
	}
	return zstdEncoder.EncodeAll(uncompressed, nil), nil

}

// decompressZstd decompress the given data using the zstd algorithm
func (c *compressor) decompressZstd(compressed []byte) (decompressed []byte, err error) {

	r := bytes.NewReader(compressed)
	zstdDecoder, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}

	decompressed, err = zstdDecoder.DecodeAll(compressed, nil)
	return decompressed, err

}
