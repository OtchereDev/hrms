package storage

import (
	"context"
	"io"
	"time"
)

// FileMetadata contains metadata about uploaded files
type FileMetadata struct {
	Filename    string
	ContentType string
	Size        int64
	UploadedAt  time.Time
	UploadedBy  string
}

// Storage defines the interface for file storage operations
// Supports both local filesystem and S3-compatible storage
type Storage interface {
	// Upload uploads a file to storage and returns the file key/path
	Upload(ctx context.Context, key string, reader io.Reader, metadata FileMetadata) (string, error)

	// Download retrieves a file from storage
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, key string) error

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, key string) (bool, error)

	// GetURL returns a public or signed URL for a file
	GetURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)

	// List lists files in a given prefix/directory
	List(ctx context.Context, prefix string) ([]string, error)

	// GetMetadata retrieves metadata for a file
	GetMetadata(ctx context.Context, key string) (*FileMetadata, error)
}

// Config contains configuration for storage backends
type Config struct {
	Type string // "local" or "s3"

	// Local storage config
	LocalBasePath string

	// S3 storage config
	S3Bucket          string
	S3Region          string
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3UseSSL          bool
}

// NewStorage creates a new storage instance based on config
func NewStorage(config Config) (Storage, error) {
	switch config.Type {
	case "s3":
		return NewS3Storage(config)
	case "local":
		return NewLocalStorage(config)
	default:
		return NewLocalStorage(config)
	}
}
