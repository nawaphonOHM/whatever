// Package firestore provides internal Firestore Testcontainers implementation,
// container lifecycle management, and option builders.
package firestore

import (
	tcfirestore "github.com/testcontainers/testcontainers-go/modules/gcloud/firestore"
)

const (
	// DefaultImage is the default Google Cloud SDK image used for Firestore emulator.
	DefaultImage = "gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators"
	// DefaultPort is the default Firestore emulator container port.
	DefaultPort = "8080/tcp"
	// DefaultProjectID is the default Google Cloud Project ID for the Firestore container.
	DefaultProjectID = tcfirestore.DefaultProjectID
)
