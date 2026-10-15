// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type getResult struct {
	out    *s3.GetObjectOutput
	err    error
	cancel context.CancelCauseFunc
	guard  *time.Timer
}

// cancelOnClose releases the request context when the caller closes the response body.
type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelCauseFunc
}

func (b *cancelOnClose) Close() error {
	err := b.ReadCloser.Close()
	b.cancel(nil)
	return err
}

// startHeaderTimer cancels the request with context.DeadlineExceeded if headers do not arrive in d.
func startHeaderTimer(d time.Duration, cancel context.CancelCauseFunc) *time.Timer {
	if d <= 0 {
		return nil
	}
	return time.AfterFunc(d, func() { cancel(context.DeadlineExceeded) })
}

func stopHeaderTimer(t *time.Timer) bool {
	return t == nil || t.Stop()
}

func closeBody(out *s3.GetObjectOutput) error {
	if out == nil || out.Body == nil {
		return nil
	}
	return out.Body.Close()
}

func timeoutError(err error) error {
	if err == nil {
		return fmt.Errorf("get object: %w", context.DeadlineExceeded)
	}
	return fmt.Errorf("get object: %w: %w", context.DeadlineExceeded, err)
}

// failure returns the mapped error for a failed or timed-out request, or nil on success.
func (r getResult) failure() error {
	if !stopHeaderTimer(r.guard) {
		return timeoutError(r.err)
	}
	if r.err != nil {
		return mapError(r.err)
	}
	return nil
}

// release closes a discarded response body and the request context.
func (r getResult) release() error {
	err := closeBody(r.out)
	r.cancel(nil)
	return err
}

func finishGet(r getResult) (*contracts.GetObjectOutput, error) {
	if err := r.failure(); err != nil {
		return nil, errors.Join(err, r.release())
	}
	res := mapGetObjectOutput(r.out)
	res.Body = &cancelOnClose{ReadCloser: r.out.Body, cancel: r.cancel}
	return res, nil
}
