package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsFirestoreURL tests detection of Google Cloud Firestore endpoints in URLs and hostnames.
func TestIsFirestoreURL(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		expected bool
	}{
		{
			name:     "bare firestore hostname",
			target:   "cluster0.firestore.goog",
			expected: true,
		},
		{
			name:     "full mongodb connection uri with firestore host",
			target:   "mongodb://user:pass@project.firestore.goog:27017",
			expected: true,
		},
		{
			name:     "mongodb+srv connection uri with uppercase firestore host",
			target:   "mongodb+srv://user:pass@my-project.FIRESTORE.GOOG/",
			expected: true,
		},
		{
			name:     "mixed case firestore target",
			target:   "Project.Firestore.Goog",
			expected: true,
		},
		{
			name:     "localhost target",
			target:   "localhost",
			expected: false,
		},
		{
			name:     "mongodb atlas cluster target",
			target:   "cluster0.mongodb.net",
			expected: false,
		},
		{
			name:     "empty string target",
			target:   "",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, IsFirestoreURL(tc.target))
		})
	}
}

// TestConfig_IsFirestore tests Firestore detection on Config instances.
func TestConfig_IsFirestore(t *testing.T) {
	t.Run("nil config returns false", func(t *testing.T) {
		var cfg *Config
		assert.False(t, cfg.IsFirestore())
	})

	t.Run("firestore host returns true", func(t *testing.T) {
		cfg := &Config{
			Host: "project.firestore.goog",
		}
		assert.True(t, cfg.IsFirestore())
	})

	t.Run("uppercase firestore host returns true", func(t *testing.T) {
		cfg := &Config{
			Host: "PROJECT.FIRESTORE.GOOG",
		}
		assert.True(t, cfg.IsFirestore())
	})

	t.Run("standard mongodb host returns false", func(t *testing.T) {
		cfg := &Config{
			Host: "localhost",
		}
		assert.False(t, cfg.IsFirestore())
	})

	t.Run("empty host returns false", func(t *testing.T) {
		cfg := &Config{
			Host: "",
		}
		assert.False(t, cfg.IsFirestore())
	})
}
