package provider

import (
	"context"
	"io"
	"time"
)

// Metadata describes a backup artifact stored in a provider.
type Metadata struct {
	Key          string
	Size         int64
	LastModified time.Time
	Checksum     string
}

// StorageProvider is the abstraction for backup storage destinations.
// Implementations must not know about Kiroku's XLSX generation,
// scheduling, or fingerprinting logic.
type StorageProvider interface {
	// Name returns the provider identifier used in state tracking.
	Name() string

	// Upload stores an artifact under key. r must be fully consumed or closed.
	Upload(ctx context.Context, key string, r io.Reader, meta Metadata) error

	// Download retrieves an artifact by key. The caller must close the reader.
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// List returns all known artifacts with their metadata.
	List(ctx context.Context) ([]Metadata, error)

	// Delete removes an artifact by key.
	Delete(ctx context.Context, key string) error

	// HealthCheck verifies provider connectivity and credential validity.
	HealthCheck(ctx context.Context) error
}
