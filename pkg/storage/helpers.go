// Package storage provides the public entrypoint and abstract contract
// for cloud object storage operations across AWS S3 and S3-compatible providers.
package storage

import (
	"bytes"
	"io"
)

// ReadAllBody reads all data from GetObjectOutput.Body and closes it safely.
// If output is nil or output.Body is nil, it returns nil bytes and nil error.
func ReadAllBody(output *GetObjectOutput) ([]byte, error) {
	if output == nil {
		return nil, nil
	}
	return readAndClose(output.Body)
}

func readAndClose(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	data, err := io.ReadAll(body)
	closeErr := body.Close()
	if err != nil {
		return data, err
	}
	return data, closeErr
}

// NewPutObjectBytes creates a PutObjectInput from a byte slice with given key, content type, and bucket.
func NewPutObjectBytes(bucket, key string, data []byte, contentType string) PutObjectInput {
	return PutObjectInput{
		Bucket:      bucket,
		Key:         key,
		Body:        bytes.NewReader(data),
		ContentType: contentType,
	}
}
