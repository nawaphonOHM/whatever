package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_IsFirestore(t *testing.T) {
	var nilOpts *Options
	assert.False(t, nilOpts.IsFirestore())

	assert.False(t, (&Options{Host: "localhost"}).IsFirestore())
	assert.True(t, (&Options{Host: "test.firestore.goog"}).IsFirestore())
	assert.True(t, (&Options{URI: "mongodb://user:pass@test.firestore.goog:443/db"}).IsFirestore())
}

func TestOptions_BuildURI_Firestore(t *testing.T) {
	opts := &Options{
		Host:     "abc.asia-southeast1.firestore.goog",
		Database: "testdb",
		Username: "user",
		Password: "password",
	}
	uri := BuildURI(opts, false)
	assert.Contains(t, uri, "abc.asia-southeast1.firestore.goog:443")
	assert.Contains(t, uri, "tls=true")
	assert.Contains(t, uri, "loadBalanced=true")
	assert.Contains(t, uri, "retryWrites=false")
	assert.Contains(t, uri, "authMechanism=SCRAM-SHA-256")
}

func TestOptions_BuildClientOptions_Firestore(t *testing.T) {
	opts := &Options{
		Host: "abc.asia-southeast1.firestore.goog",
	}
	clientOpts := BuildClientOptions(opts)
	assert.NotNil(t, clientOpts)
	assert.NotNil(t, clientOpts.LoadBalanced)
	assert.True(t, *clientOpts.LoadBalanced)
	assert.NotNil(t, clientOpts.RetryWrites)
	assert.False(t, *clientOpts.RetryWrites)
}
