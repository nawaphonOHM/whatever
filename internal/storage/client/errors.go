// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"errors"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func isClientSentinel(err error) bool {
	return errors.Is(err, contracts.ErrNilClient) || errors.Is(err, contracts.ErrNilConfig)
}

func isObjectOrBucketSentinel(err error) bool {
	return errors.Is(err, contracts.ErrObjectNotFound) || errors.Is(err, contracts.ErrBucketNotFound)
}

func isAccessOrExpirySentinel(err error) bool {
	return errors.Is(err, contracts.ErrAccessDenied) || errors.Is(err, contracts.ErrInvalidExpiry)
}

func isKnownSentinel(err error) bool {
	return isObjectOrBucketSentinel(err) || isClientSentinel(err) || isAccessOrExpirySentinel(err)
}

func isNilOrKnownSentinel(err error) bool {
	return err == nil || isKnownSentinel(err)
}

func matchAnyError(err error) error {
	if matched := matchTypedError(err); matched != nil {
		return matched
	}
	return matchAPIError(err)
}

func mapError(err error) error {
	if isNilOrKnownSentinel(err) {
		return err
	}
	if matched := matchAnyError(err); matched != nil {
		return matched
	}
	return err
}

func mapBucketError(err error) error {
	if err == nil {
		return nil
	}
	mapped := mapError(err)
	if errors.Is(mapped, contracts.ErrObjectNotFound) {
		return fmt.Errorf(errWrapFmt, contracts.ErrBucketNotFound, err)
	}
	return mapped
}
