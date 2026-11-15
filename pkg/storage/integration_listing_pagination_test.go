//go:build testcontainers

package storage_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

const (
	pagedListingObjectCount   = 1001
	localStackListingWorkers  = 16
	localStackDeleteBatchSize = 1000
)

func TestLocalStackListObjectsPaginates(t *testing.T) {
	requireSharedLocalStack(t)
	prefix := localStackObjectKey(t, "pages/")
	cleanupLocalStackPagedObjects(t, prefix, pagedListingObjectCount)
	putLocalStackPagedObjects(t, prefix, pagedListingObjectCount)
	listed, err := sharedStorage.ListObjects(t.Context(), storage.ListObjectsInput{Prefix: prefix})
	require.NoError(t, err)
	require.Len(t, listed.Objects, pagedListingObjectCount)
}

func putLocalStackPagedObjects(t *testing.T, prefix string, count int) {
	t.Helper()
	jobs := make(chan int)
	errs := make(chan error, count)
	var workers sync.WaitGroup
	startLocalStackListingWorkers(&workers, prefix, jobs, errs)
	queueLocalStackListingJobs(jobs, count)
	waitLocalStackListingWorkers(&workers, errs, t)
}

func startLocalStackListingWorkers(workers *sync.WaitGroup, prefix string, jobs <-chan int, errs chan<- error) {
	for range localStackListingWorkers {
		workers.Go(func() { putLocalStackPagedObjectsWorker(prefix, jobs, errs) })
	}
}

func queueLocalStackListingJobs(jobs chan<- int, count int) {
	for index := range count {
		jobs <- index
	}
	close(jobs)
}

func waitLocalStackListingWorkers(workers *sync.WaitGroup, errs chan error, t *testing.T) {
	workers.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}
}

func putLocalStackPagedObjectsWorker(prefix string, jobs <-chan int, errs chan<- error) {
	for index := range jobs {
		key := fmt.Sprintf("%s%04d", prefix, index)
		_, err := sharedStorage.PutObject(context.Background(), storage.PutObjectInput{
			Key: key, Body: strings.NewReader("pagination payload"),
		})
		if err != nil {
			errs <- fmt.Errorf("put paginated listing object %q: %w", key, err)
		}
	}
}

func cleanupLocalStackPagedObjects(t *testing.T, prefix string, count int) {
	t.Helper()
	t.Cleanup(func() {
		for start := 0; start < count; start += localStackDeleteBatchSize {
			end := min(start+localStackDeleteBatchSize, count)
			deleteLocalStackPagedObjectBatch(t, prefix, start, end)
		}
	})
}

func deleteLocalStackPagedObjectBatch(t *testing.T, prefix string, start, end int) {
	objects := make([]types.ObjectIdentifier, 0, end-start)
	for index := start; index < end; index++ {
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(fmt.Sprintf("%s%04d", prefix, index))})
	}
	_, err := sharedStorage.RawClient().DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: aws.String(sharedStorage.DefaultBucket()),
		Delete: &types.Delete{Objects: objects, Quiet: aws.Bool(true)},
	})
	if err != nil {
		t.Errorf("delete paginated listing objects %q: %v", prefix, err)
	}
}
