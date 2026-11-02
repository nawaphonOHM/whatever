// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type gcsAcceptEncodingKey struct{}

func handleDropGCSAcceptEncoding(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return next.HandleFinalize(ctx, in)
	}
	ae := req.Header.Get("Accept-Encoding")
	req.Header.Del("Accept-Encoding")
	nextCtx := middleware.WithStackValue(ctx, gcsAcceptEncodingKey{}, ae)
	return next.HandleFinalize(nextCtx, in)
}

func restoreAcceptEncodingHeader(ctx context.Context, req *smithyhttp.Request) {
	if ae, ok := middleware.GetStackValue(ctx, gcsAcceptEncodingKey{}).(string); ok && ae != "" {
		req.Header.Set("Accept-Encoding", ae)
	}
}

func handleRestoreGCSAcceptEncoding(
	ctx context.Context, in middleware.FinalizeInput, next middleware.FinalizeHandler,
) (middleware.FinalizeOutput, middleware.Metadata, error) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		restoreAcceptEncodingHeader(ctx, req)
	}
	return next.HandleFinalize(ctx, in)
}

var dropGCSAcceptEncoding = middleware.FinalizeMiddlewareFunc(
	"DropGCSAcceptEncoding",
	handleDropGCSAcceptEncoding,
)

var restoreGCSAcceptEncoding = middleware.FinalizeMiddlewareFunc(
	"RestoreGCSAcceptEncoding",
	handleRestoreGCSAcceptEncoding,
)

func applyGCSSigningMiddleware(stack *middleware.Stack) error {
	if err := stack.Finalize.Insert(dropGCSAcceptEncoding, "Signing", middleware.Before); err != nil {
		return err
	}
	return stack.Finalize.Insert(restoreGCSAcceptEncoding, "Signing", middleware.After)
}

func applyGCSSigningOptions(o *s3.Options, provider contracts.Provider) {
	if provider != contracts.ProviderGCS {
		return
	}
	o.APIOptions = append(o.APIOptions, applyGCSSigningMiddleware)
}
