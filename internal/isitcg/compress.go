package isitcg

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

func Compress(uncompressedString string) string {
	uncompressedBytes := []byte(uncompressedString)
	var compressedBytes bytes.Buffer

	// create a new compressor that writes to the buffer
	compressor := zlib.NewWriter(&compressedBytes)
	// compress the uncompressed data
	compressor.Write(uncompressedBytes)
	// close the compressor to ensure that all data has been written
	compressor.Close()

	return base64.URLEncoding.EncodeToString(compressedBytes.Bytes())
}

func Decompress(compressedString string) string {
	compressedBytes, err := base64.URLEncoding.DecodeString(compressedString)
	if err != nil {
		// handle error
		return ""
	}

	// create a new decompressor that reads from the compressed bytes
	decompressor, err := zlib.NewReader(bytes.NewReader(compressedBytes))
	if err != nil {
		// handle error
		return ""
	}
	defer decompressor.Close()

	// create a buffer to hold the decompressed data
	var decompressedBytes bytes.Buffer
	// decompress the data
	io.Copy(&decompressedBytes, decompressor)

	return decompressedBytes.String()
}
