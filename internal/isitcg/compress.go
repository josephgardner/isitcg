package isitcg

import (
	"bytes"
	"compress/flate"

	"encoding/base64"
	"io"
)

func Compress(uncompressedString string) string {
	var compressedBytes bytes.Buffer
	compressor, _ := flate.NewWriter(&compressedBytes, flate.DefaultCompression)
	compressor.Write([]byte(uncompressedString))
	compressor.Close()
	return base64.RawURLEncoding.EncodeToString(compressedBytes.Bytes())
}

func Decompress(compressedString string) string {
	compressedBytes, err := base64.RawURLEncoding.DecodeString(compressedString)
	if err != nil {
		return ""
	}

	decompressor := flate.NewReader(bytes.NewReader(compressedBytes))
	defer decompressor.Close()

	var decompressedBytes bytes.Buffer
	io.Copy(&decompressedBytes, decompressor)

	return decompressedBytes.String()
}
