// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func applyListFilters(params *s3.ListObjectsV2Input, input contracts.ListObjectsInput) {
	if input.Delimiter != "" {
		params.Delimiter = aws.String(input.Delimiter)
	}
	if input.MaxKeys > 0 {
		params.MaxKeys = aws.Int32(input.MaxKeys)
	}
}

func buildListParams(bucket string, input contracts.ListObjectsInput) *s3.ListObjectsV2Input {
	params := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if input.Prefix != "" {
		params.Prefix = aws.String(input.Prefix)
	}
	applyListFilters(params, input)
	return params
}

func stopOnDuplicateToken(o *s3.ListObjectsV2PaginatorOptions) {
	o.StopOnDuplicateToken = true
}

func (c *Client) newListRun(input contracts.ListObjectsInput) *listRun {
	params := buildListParams(c.resolveBucket(input.Bucket), input)
	return &listRun{
		client:    c,
		pager:     s3.NewListObjectsV2Paginator(c.rawClient, params, stopOnDuplicateToken),
		collector: listCollector{limit: input.MaxKeys},
	}
}

// ListObjects enumerates objects and common prefixes matching criteria, following continuation
// tokens across pages until the listing ends or MaxKeys entries (objects plus common prefixes, as
// S3 counts them) were collected; MaxKeys of 0 lists everything. RequestTimeout bounds each page
// request. If any page fails, the error is returned and the partial listing is discarded.
func (c *Client) ListObjects(
	ctx context.Context,
	input contracts.ListObjectsInput,
) (*contracts.ListObjectsOutput, error) {
	if c.isNil() {
		return nil, contracts.ErrNilClient
	}
	run := c.newListRun(input)
	if err := run.drain(ctx); err != nil {
		return nil, mapError(err)
	}
	return run.collector.output(), nil
}
