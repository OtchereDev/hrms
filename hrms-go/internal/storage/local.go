package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage implements Storage interface for local filesystem
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local filesystem storage
func NewLocalStorage(config Config) (*LocalStorage, error) {
	if config.LocalBasePath == "" {
		config.LocalBasePath = "./storage"
	}

	// Create base directory if it doesn't exist
	if err := os.MkdirAll(config.LocalBasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: config.LocalBasePath,
		baseURL:  "/api/files", // Base URL for serving files
	}, nil
}

// Upload uploads a file to local filesystem
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, metadata FileMetadata) (string, error) {
	fullPath := filepath.Join(s.basePath, key)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return key, nil
}

// Download retrieves a file from local filesystem
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, key)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from local filesystem
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)

	err := os.Remove(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.basePath, key)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetURL returns a URL for accessing the file
func (s *LocalStorage) GetURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	// For local storage, return a simple URL path
	// In production, you might want to implement signed URLs
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

// List lists files in a directory
func (s *LocalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	fullPath := filepath.Join(s.basePath, prefix)

	var files []string

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(s.basePath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// GetMetadata retrieves file metadata
func (s *LocalStorage) GetMetadata(ctx context.Context, key string) (*FileMetadata, error) {
	fullPath := filepath.Join(s.basePath, key)

	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileMetadata{
		Filename:   filepath.Base(key),
		Size:       info.Size(),
		UploadedAt: info.ModTime(),
	}, nil
}
