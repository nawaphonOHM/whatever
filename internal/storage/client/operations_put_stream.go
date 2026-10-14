// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"bytes"
	"context"
	"io"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// streamPartSize is the in-memory buffer size and multipart part size for non-seekable bodies.
const streamPartSize = 8 * 1024 * 1024

// maxSinglePutSize is the maximum size (5 GiB) allowed by S3 for a single PutObject call.
const maxSinglePutSize = 5 * 1024 * 1024 * 1024

// maxMultipartSize is the maximum upload capacity (~78.1 GiB) with 10000 parts of 8 MiB.
const maxMultipartSize = int64(maxMultipartParts) * streamPartSize

// isSeekable reports whether body supports a working Seek (pipes and files that cannot seek are not).
func isSeekable(body io.Reader) bool {
	seeker, ok := body.(io.Seeker)
	if !ok {
		return false
	}
	_, err := seeker.Seek(0, io.SeekCurrent)
	return err == nil
}

func restoreSeek(seeker io.Seeker, pos int64) bool {
	_, err := seeker.Seek(pos, io.SeekStart)
	return err == nil
}

func measureSeek(seeker io.Seeker, curr int64) (int64, bool) {
	end, err := seeker.Seek(0, io.SeekEnd)
	ok := restoreSeek(seeker, curr)
	if err != nil || !ok {
		return 0, false
	}
	return end - curr, true
}

func seekableSize(seeker io.Seeker) (int64, bool) {
	curr, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, false
	}
	return measureSeek(seeker, curr)
}

func isLargeSeekable(body io.Reader) bool {
	seeker, ok := body.(io.Seeker)
	if !ok {
		return false
	}
	size, ok := seekableSize(seeker)
	return ok && size > maxSinglePutSize
}

func isOversizedSeekable(body io.Reader) bool {
	seeker, ok := body.(io.Seeker)
	if !ok {
		return false
	}
	size, ok := seekableSize(seeker)
	return ok && size > maxMultipartSize
}

func readChunk(r io.Reader, size int64) ([]byte, error) {
	return io.ReadAll(io.LimitReader(r, size))
}

// putStream uploads a non-seekable body: one request when it fits in a single part, multipart otherwise.
func (c *Client) putStream(ctx context.Context, input contracts.PutObjectInput) (*contracts.PutObjectOutput, error) {
	head, err := readChunk(input.Body, streamPartSize+1)
	if err != nil {
		return nil, err
	}
	if len(head) <= streamPartSize {
		return c.putSingle(ctx, input, bytes.NewReader(head))
	}
	rest := io.MultiReader(bytes.NewReader(head[streamPartSize:]), input.Body)
	return c.putMultipart(ctx, input, head[:streamPartSize], rest)
}
