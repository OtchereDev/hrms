package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Storage implements Storage interface for S3-compatible storage
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(cfg Config) (*S3Storage, error) {
	// Load AWS config
	var awsCfg aws.Config
	var err error

	if cfg.S3Endpoint != "" {
		// Custom endpoint (MinIO, DigitalOcean Spaces, etc.)
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.S3Endpoint,
					SigningRegion: cfg.S3Region,
				}, nil
			})

		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.S3Region),
			config.WithEndpointResolverWithOptions(customResolver),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cfg.S3AccessKeyID,
					cfg.S3SecretAccessKey,
					"",
				),
			),
		)
	} else {
		// Standard AWS S3
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.S3Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cfg.S3AccessKeyID,
					cfg.S3SecretAccessKey,
					"",
				),
			),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // Required for MinIO and some S3-compatible services
	})

	return &S3Storage{
		client: client,
		bucket: cfg.S3Bucket,
	}, nil
}

// Upload uploads a file to S3
func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, metadata FileMetadata) (string, error) {
	// Prepare metadata
	metadataMap := map[string]string{
		"uploaded-by": metadata.UploadedBy,
		"uploaded-at": metadata.UploadedAt.Format(time.RFC3339),
	}

	// Upload to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(metadata.ContentType),
		Metadata:    metadataMap,
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

// Download retrieves a file from S3
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

// Exists checks if a file exists in S3
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if error is NotFound
		return false, nil
	}

	return true, nil
}

// GetURL returns a presigned URL for a file
func (s *S3Storage) GetURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiresIn
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}

// List lists files in S3 with a given prefix
func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var keys []string
	for _, obj := range result.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

// GetMetadata retrieves file metadata from S3
func (s *S3Storage) GetMetadata(ctx context.Context, key string) (*FileMetadata, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	var uploadedAt time.Time
	if result.LastModified != nil {
		uploadedAt = *result.LastModified
	}

	metadata := &FileMetadata{
		Filename:    key,
		ContentType: aws.ToString(result.ContentType),
		Size:        aws.ToInt64(result.ContentLength),
		UploadedAt:  uploadedAt,
	}

	if uploadedBy, ok := result.Metadata["uploaded-by"]; ok {
		metadata.UploadedBy = uploadedBy
	}

	return metadata, nil
}
