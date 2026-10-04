package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testCustomProbeDur = 5 * time.Second
	testNegativeDur    = -1 * time.Second
	testProbeHelperURI = "mongodb://127.0.0.1:59999"
)

func TestResolveProbeTimeout(t *testing.T) {
	assert.Equal(t, DefaultSocketTimeout, resolveProbeTimeout(nil))
	oZero := &Options{}
	assert.Equal(t, DefaultSocketTimeout, resolveProbeTimeout(oZero))

	oCustom := &Options{SocketTimeout: testCustomProbeDur}
	assert.Equal(t, testCustomProbeDur, resolveProbeTimeout(oCustom))
}

func TestProbeTimeoutContextHelpers(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, DefaultSocketTimeout, probeTimeoutFromContext(ctx))

	ctxCustom := withProbeTimeout(ctx, testCustomProbeDur)
	assert.Equal(t, testCustomProbeDur, probeTimeoutFromContext(ctxCustom))

	ctxZero := withProbeTimeout(ctx, 0)
	assert.Equal(t, DefaultSocketTimeout, probeTimeoutFromContext(ctxZero))

	ctxNegative := withProbeTimeout(ctx, testNegativeDur)
	assert.Equal(t, DefaultSocketTimeout, probeTimeoutFromContext(ctxNegative))
}

func TestResolveDatabaseHelpers(t *testing.T) {
	assert.Equal(t, "admin", resolveDBName(""))
	assert.Equal(t, "mydb", resolveDBName("mydb"))
	assert.Empty(t, resolveDatabaseName(nil))

	o := &Options{Database: "customdb"}
	assert.Equal(t, "customdb", resolveDatabaseName(o))

	assert.Equal(t, "__probe__", chooseTargetColl(nil))
	assert.Equal(t, "__probe__", chooseTargetColl(make([]string, 0)))

	colls := []string{"coll1", "coll2", "coll3"}
	assert.Contains(t, colls, chooseTargetColl(colls))
}

func getDisconnectedDatabase(t *testing.T) (*mongo.Client, *mongo.Database) {
	t.Helper()
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testProbeHelperURI))
	require.NoError(t, err)
	require.NoError(t, rawClient.Disconnect(context.Background()))
	return rawClient, rawClient.Database("testdb")
}

func TestProbeListCollections_DisconnectedClient(t *testing.T) {
	_, db := getDisconnectedDatabase(t)
	_, errNoMax := probeListCollectionsWithoutMaxTime(context.Background(), db, testCustomProbeDur)
	assert.ErrorContains(t, errNoMax, errStep2Substr)

	_, errNoMaxZero := probeListCollectionsWithoutMaxTime(context.Background(), db, 0)
	assert.ErrorContains(t, errNoMaxZero, errStep2Substr)

	errWithMax := probeListCollectionsWithMaxTime(context.Background(), db, testCustomProbeDur)
	assert.ErrorContains(t, errWithMax, errStep3Substr)

	errWithMaxZero := probeListCollectionsWithMaxTime(context.Background(), db, 0)
	assert.ErrorContains(t, errWithMaxZero, errStep3Substr)
}

func TestProbeDocumentAndClient_DisconnectedClient(t *testing.T) {
	rawClient, db := getDisconnectedDatabase(t)
	assert.ErrorContains(t, probeDocument(context.Background(), db, "__probe__"), errStep4Substr)
	assert.Error(t, defaultProbeClient(context.Background(), rawClient, "testdb"))
}
