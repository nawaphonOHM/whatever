// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"errors"
	"fmt"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

const errWrapFmt = "%w: %w"

func matchNotFoundType(err error) error {
	if nsk, ok := errors.AsType[*types.NoSuchKey](err); ok {
		return fmt.Errorf(errWrapFmt, contracts.ErrObjectNotFound, nsk)
	}
	if nf, ok := errors.AsType[*types.NotFound](err); ok {
		return fmt.Errorf(errWrapFmt, contracts.ErrObjectNotFound, nf)
	}
	return nil
}

func matchNoSuchBucketType(err error) error {
	if nsb, ok := errors.AsType[*types.NoSuchBucket](err); ok {
		return fmt.Errorf(errWrapFmt, contracts.ErrBucketNotFound, nsb)
	}
	return nil
}

func matchTypedError(err error) error {
	if mapped := matchNotFoundType(err); mapped != nil {
		return mapped
	}
	return matchNoSuchBucketType(err)
}

func isNotFoundOr404(code string) bool {
	return code == "NotFound" || code == "404"
}

func isObjectNotFoundCode(code string) bool {
	return code == "NoSuchKey" || isNotFoundOr404(code)
}

func isAccessDeniedCode(code string) bool {
	return code == "AccessDenied" || code == "Forbidden" || code == "403"
}

func matchAPIErrorCode(code string, err error) error {
	if isObjectNotFoundCode(code) {
		return fmt.Errorf(errWrapFmt, contracts.ErrObjectNotFound, err)
	}
	return matchBucketOrAccessCode(code, err)
}

func matchBucketOrAccessCode(code string, err error) error {
	if code == "NoSuchBucket" {
		return fmt.Errorf(errWrapFmt, contracts.ErrBucketNotFound, err)
	}
	if isAccessDeniedCode(code) {
		return fmt.Errorf(errWrapFmt, contracts.ErrAccessDenied, err)
	}
	return nil
}

func mapHTTPStatus(status int, err error) error {
	switch status {
	case http.StatusNotFound:
		return fmt.Errorf(errWrapFmt, contracts.ErrObjectNotFound, err)
	case http.StatusForbidden:
		return fmt.Errorf(errWrapFmt, contracts.ErrAccessDenied, err)
	}
	return nil
}

func matchHTTPStatus(err error) error {
	respErr, ok := errors.AsType[*awshttp.ResponseError](err)
	if !ok {
		return nil
	}
	return mapHTTPStatus(respErr.HTTPStatusCode(), err)
}

func matchSmithyAPIError(err error) error {
	apiErr, ok := errors.AsType[smithy.APIError](err)
	if !ok {
		return nil
	}
	return matchAPIErrorCode(apiErr.ErrorCode(), err)
}

func matchAPIError(err error) error {
	if mapped := matchSmithyAPIError(err); mapped != nil {
		return mapped
	}
	return matchHTTPStatus(err)
}
