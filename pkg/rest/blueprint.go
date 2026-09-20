// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// CorsSetting configures Cross-Origin Resource Sharing policy options.
type CorsSetting = contracts.CorsSetting

// NewCorsSetting initializes a new CorsSetting builder.
func NewCorsSetting() *CorsSetting {
	return contracts.NewCorsSetting()
}

// Meta holds server configuration metadata such as CORS settings.
type Meta = contracts.Meta

// NewMeta initializes a new Meta builder.
func NewMeta() *Meta {
	return contracts.NewMeta()
}

// BluePrint defines the application blueprint including metadata and APIs.
type BluePrint = contracts.BluePrint

// NewBluePrint initializes a new BluePrint builder.
func NewBluePrint() *BluePrint {
	return contracts.NewBluePrint()
}
