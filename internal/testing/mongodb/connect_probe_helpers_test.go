package mongodb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
)

const (
	testCustomProbeDur = 5 * time.Second
	testNegativeDur    = -1 * time.Second
	testProbeHelperURI = "mongodb://127.0.0.1:59999"
	testDBName         = "testdb"
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

	assert.Equal(t, "probe", chooseTargetColl(nil))
	assert.Equal(t, "probe", chooseTargetColl(make([]string, 0)))

	colls := []string{"coll1", "coll2", "coll3"}
	assert.Contains(t, colls, chooseTargetColl(colls))
}

func getDisconnectedDatabase(t *testing.T) (*mongo.Client, *mongo.Database) {
	t.Helper()
	rawClient, err := mongo.Connect(options.Client().ApplyURI(testProbeHelperURI))
	require.NoError(t, err)
	require.NoError(t, rawClient.Disconnect(context.Background()))
	return rawClient, rawClient.Database(testDBName)
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
	assert.ErrorContains(t, probeDocument(context.Background(), db, "probe"), errStep4Substr)
	c := client.NewClient(rawClient, testDBName)
	assert.Error(t, defaultProbeClient(context.Background(), c, testDBName))
	assert.Equal(t, ErrNilClient, defaultProbeClient(context.Background(), nil, testDBName))
}

func TestSetMockProbe_OverridesAndRestores(t *testing.T) {
	customErr := errors.New("custom mock probe error")
	restore := SetMockProbe(func(context.Context, *client.Client, string) error {
		return customErr
	})
	defer restore()

	assert.Equal(t, customErr, probeClient(context.Background(), nil, testDBName))

	restore()
	assert.Equal(t, ErrNilClient, probeClient(context.Background(), nil, testDBName))
}

func TestSetMockProbe_NilRestoresDefault(t *testing.T) {
	restore := SetMockProbe(func(context.Context, *client.Client, string) error {
		return errors.New("temporary error")
	})
	defer restore()

	SetMockProbe(nil)
	assert.Equal(t, ErrNilClient, probeClient(context.Background(), nil, testDBName))
}
