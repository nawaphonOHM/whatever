package mongodb

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDummyQuery         = "dummyQuery"
	testExpectedDefaultURI = "mongodb://localhost:27017/?uuidRepresentation=unspecified&tls=false"
	testExpectedTLSURI     = "mongodb://localhost:27017/?uuidRepresentation=unspecified&tls=true"
)

func TestBuildURI_NilOptionsProof(t *testing.T) {
	assert.Nil(t, buildUserInfo(nil))
	assert.Nil(t, buildUserInfo(&Options{}))
	assert.False(t, hasCredentials(nil))
	assert.Equal(t, "localhost:27017", buildHost(nil))
	assert.Equal(t, "/", buildPath(nil))
	assert.Equal(t, "uuidRepresentation=unspecified&tls=false", buildQuery(nil, false))
	assert.Equal(t, "uuidRepresentation=unspecified&tls=true", buildQuery(nil, true))
}

func TestBuildURI_HelperNilProof(t *testing.T) {
	assert.Equal(t, testDummyQuery, appendAuthSourceQuery(testDummyQuery, nil))
	assert.Equal(t, testDummyQuery, appendAppNameQuery(testDummyQuery, nil))
	assert.Equal(t, testDummyQuery, appendDirectConnQuery(testDummyQuery, nil))
	assert.Equal(t, DefaultPort, resolvePort(nil))
	assert.Equal(t, DefaultHost, resolveHost(nil))
	assert.False(t, isSrvOrZeroPort(nil))
	assert.Equal(t, DefaultUUIDRepresentation, resolveUUIDRep(nil))
}

func TestBuildURI_URLStructNilUserProof(t *testing.T) {
	u := &url.URL{
		Scheme:   DefaultProtocol,
		User:     buildUserInfo(nil),
		Host:     buildHost(nil),
		Path:     buildPath(nil),
		RawQuery: buildQuery(nil, false),
	}
	assert.Nil(t, u.User)
	assert.Equal(t, testExpectedDefaultURI, u.String())
	assert.Equal(t, testExpectedDefaultURI, BuildURI(nil, false))
	assert.Equal(t, testExpectedTLSURI, BuildURI(nil, true))
}
