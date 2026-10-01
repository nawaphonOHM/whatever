// Package mongodb provides internal MongoDB Testcontainers implementation,
// container lifecycle management, and option builders.
package mongodb

const (
	// DefaultImage is the default MongoDB container image.
	DefaultImage = "mongo:6"
	// DefaultPort is the default MongoDB container port.
	DefaultPort = "27017/tcp"
)
