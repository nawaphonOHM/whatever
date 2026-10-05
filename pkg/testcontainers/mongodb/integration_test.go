//go:build testcontainers

package mongodb_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testcontainers/mongodb"
)

const (
	testDBName  = "integration_test_db"
	testUser    = "testuser"
	testPass    = "testpass"
	testColName = "test_collection"
)

type testDoc struct {
	ID    string `bson:"_id,omitempty"`
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

var (
	sharedContainer *mongodb.Container
	initErr         error
)

func setupSharedContainer() error {
	ctx := context.Background()
	var err error
	dbOpt := mongodb.WithDatabase(testDBName)
	uOpt, pOpt := mongodb.WithUsername(testUser), mongodb.WithPassword(testPass)
	sharedContainer, err = mongodb.Run(ctx, dbOpt, uOpt, pOpt)
	return err
}

func teardownSharedContainer() {
	if sharedContainer == nil {
		return
	}
	if err := sharedContainer.Terminate(context.Background()); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
}

func TestMain(m *testing.M) {
	initErr = setupSharedContainer()
	defer teardownSharedContainer()
	m.Run()
}

func requireSharedReady(t *testing.T) {
	t.Helper()
	require.NoError(t, initErr, "shared container failed to initialize")
	require.NotNil(t, sharedContainer, "shared container is nil")
}
