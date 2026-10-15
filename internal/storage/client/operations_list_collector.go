// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

type listCollector struct {
	objects  []contracts.ObjectSummary
	prefixes []string
	limit    int32
}

func (c *listCollector) full() bool {
	return c.limit > 0 && len(c.objects)+len(c.prefixes) >= int(c.limit)
}

func (c *listCollector) addObjects(items []types.Object) {
	for _, item := range items {
		if c.full() {
			return
		}
		c.objects = append(c.objects, contracts.ObjectSummary{
			Key:          aws.ToString(item.Key),
			Size:         aws.ToInt64(item.Size),
			ETag:         aws.ToString(item.ETag),
			LastModified: aws.ToTime(item.LastModified),
		})
	}
}

func (c *listCollector) addPrefixes(items []types.CommonPrefix) {
	for _, item := range items {
		if c.full() {
			return
		}
		c.prefixes = append(c.prefixes, aws.ToString(item.Prefix))
	}
}

// collectPage records a page and reports whether the requested limit has been reached.
func (c *listCollector) collectPage(page *s3.ListObjectsV2Output) bool {
	c.addObjects(page.Contents)
	c.addPrefixes(page.CommonPrefixes)
	return c.full()
}

func (c *listCollector) output() *contracts.ListObjectsOutput {
	return &contracts.ListObjectsOutput{Objects: c.objects, Prefixes: c.prefixes}
}

type listRun struct {
	client    *Client
	pager     *s3.ListObjectsV2Paginator
	collector listCollector
}

func (r *listRun) nextPage(ctx context.Context) (bool, error) {
	pageCtx, cancel := r.client.requestContext(ctx)
	defer cancel()
	page, err := r.pager.NextPage(pageCtx)
	if err != nil {
		return false, err
	}
	return r.collector.collectPage(page), nil
}

func shouldStop(stop bool, err error) bool {
	return stop || err != nil
}

func (r *listRun) drain(ctx context.Context) error {
	if !r.pager.HasMorePages() {
		return nil
	}
	stop, err := r.nextPage(ctx)
	if shouldStop(stop, err) {
		return err
	}
	return r.drain(ctx)
}
