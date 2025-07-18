package compressor

import (
	_ "embed"
	"github.com/stretchr/testify/suite"
	"testing"
)

const lipsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. In vel hendrerit diam. Nam fringilla arcu quis aliquet finibus. In viverra condimentum arcu, nec dignissim augue bibendum a."

//go:embed test-data/test-data.json
var jsonData []byte

type CompressorSuite struct {
	suite.Suite
}

func (c *CompressorSuite) TestGzip() {

	// test data
	compressorInterface := New(Gzip)
	compressed, err := compressorInterface.Compress(jsonData)
	c.Nil(err)

	decompressed, err := compressorInterface.Decompress(compressed)
	c.Nil(err)
	c.Equal(jsonData, decompressed)

}

func (c *CompressorSuite) TestLZ4() {

	// test data
	compressorInterface := New(LZ4)
	compressed, err := compressorInterface.Compress(jsonData)
	c.Nil(err)

	decompressed, err := compressorInterface.Decompress(compressed)
	c.Nil(err)
	c.Equal(jsonData, decompressed)

}

func (c *CompressorSuite) TestSnappy() {

	// test data
	compressorInterface := New(Snappy)
	compressed, err := compressorInterface.Compress(jsonData)
	c.Nil(err)

	decompressed, err := compressorInterface.Decompress(compressed)
	c.Nil(err)
	c.Equal(jsonData, decompressed)

}

func (c *CompressorSuite) TestZstd() {

	// test data
	compressorInterface := New(Zstd)
	compressed, err := compressorInterface.Compress(jsonData)
	c.Nil(err)

	decompressed, err := compressorInterface.Decompress(compressed)
	c.Nil(err)
	c.Equal(jsonData, decompressed)

}

func TestNew(t *testing.T) {
	suite.Run(t, new(CompressorSuite))
}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Gzip
// BenchmarkCompressor_Compress_Gzip_lipsum-26         8630            149521 ns/op
func BenchmarkCompressor_Compress_Gzip_lipsum(b *testing.B) {

	// test data
	compressorInterface := New(Gzip)
	data := []byte(lipsum)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress(data)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_Gzip
// BenchmarkCompressor_Decompress_Gzip-26            147265              8003 ns/op
func BenchmarkCompressor_Decompress_Gzip(b *testing.B) {

	compressorInterface := New(Gzip)
	compressed, _ := compressorInterface.Compress([]byte(lipsum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Gzip_json
// BenchmarkCompressor_Compress_Gzip_json-32            907           1302005 ns/op
func BenchmarkCompressor_Compress_Gzip_json(b *testing.B) {

	// test data
	compressorInterface := New(Gzip)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress(jsonData)
	}

}

func BenchmarkCompressor_Decompress_Gzip_json(b *testing.B) {

	compressorInterface := New(Gzip)
	compressed, _ := compressorInterface.Compress(jsonData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_LZ4
// BenchmarkCompressor_Compress_LZ4-32         1566            760116 ns/op
func BenchmarkCompressor_Compress_LZ4(b *testing.B) {

	// test data
	compressorInterface := New(LZ4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress([]byte(lipsum))
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_LZ4
// BenchmarkCompressor_Decompress_LZ4-26               2576            484257 ns/op
func BenchmarkCompressor_Decompress_LZ4(b *testing.B) {

	// test data
	compressorInterface := New(LZ4)
	compressed, _ := compressorInterface.Compress([]byte(lipsum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_LZ4_json
// BenchmarkCompressor_Compress_LZ4_json-26            1616            727394 ns/op
func BenchmarkCompressor_Compress_LZ4_json(b *testing.B) {

	// test data
	compressorInterface := New(LZ4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress(jsonData)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_LZ4_json
// BenchmarkCompressor_Decompress_LZ4_json-26          1774            746229 ns/op
func BenchmarkCompressor_Decompress_LZ4_json(b *testing.B) {

	// test data
	compressorInterface := New(LZ4)
	compressed, _ := compressorInterface.Compress(jsonData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Snappy
// BenchmarkCompressor_Compress_Snappy-32           6855946               174.0 ns/
func BenchmarkCompressor_Compress_Snappy(b *testing.B) {

	// test data
	compressorInterface := New(Snappy)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress([]byte(lipsum))
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_Snappy
// BenchmarkCompressor_Decompress_Snappy-26        24138042                53.28 ns
func BenchmarkCompressor_Decompress_Snappy(b *testing.B) {

	// test data
	compressorInterface := New(Snappy)
	compressed, _ := compressorInterface.Compress([]byte(lipsum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Snappy_json
// BenchmarkCompressor_Compress_Snappy_json-26        13722             83310 ns/op
func BenchmarkCompressor_Compress_Snappy_json(b *testing.B) {

	// test data
	compressorInterface := New(Snappy)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress(jsonData)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_Snappy_json
// BenchmarkCompressor_Decompress_Snappy_json-26              40098             310
func BenchmarkCompressor_Decompress_Snappy_json(b *testing.B) {

	// test data
	compressorInterface := New(Snappy)
	compressed, _ := compressorInterface.Compress(jsonData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Zstd
// BenchmarkCompressor_Compress_Zstd-26                 468           2555395 ns/op
func BenchmarkCompressor_Compress_Zstd(b *testing.B) {

	// test data
	compressorInterface := New(Zstd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress([]byte(lipsum))
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compres
// sor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_Zstd
// BenchmarkCompressor_Decompress_Zstd-26            234162              7051 ns/op
func BenchmarkCompressor_Decompress_Zstd(b *testing.B) {

	// test data
	compressorInterface := New(Zstd)
	compressed, _ := compressorInterface.Compress([]byte(lipsum))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Compress_Zstd_json
// BenchmarkCompressor_Compress_Zstd_json-26            324           3117190 ns/op
func BenchmarkCompressor_Compress_Zstd_json(b *testing.B) {

	// test data
	compressorInterface := New(Zstd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Compress(jsonData)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/well/filesystem/compressor
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkCompressor_Decompress_Zstd_json
// BenchmarkCompressor_Decompress_Zstd_json-26         6168            223583 ns/op
func BenchmarkCompressor_Decompress_Zstd_json(b *testing.B) {

	// test data
	compressorInterface := New(Zstd)
	compressed, _ := compressorInterface.Compress(jsonData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressorInterface.Decompress(compressed)
	}

}
