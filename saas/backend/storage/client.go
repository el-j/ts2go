package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps MinIO/S3 client
type Client struct {
	client     *minio.Client
	bucketName string
}

// Config holds storage configuration
type Config struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
	UseSSL     bool
}

// NewClient creates a new storage client
func NewClient(cfg Config) (*Client, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	client := &Client{
		client:     minioClient,
		bucketName: cfg.BucketName,
	}

	// Ensure bucket exists
	if err := client.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	return client, nil
}

// ensureBucket creates bucket if it doesn't exist
func (c *Client) ensureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// UploadFile uploads a file to storage
func (c *Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, size, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &UploadResult{
		Bucket:      info.Bucket,
		Key:         info.Key,
		ETag:        info.ETag,
		Size:        info.Size,
		ContentType: contentType,
	}, nil
}

// DownloadFile downloads a file from storage
func (c *Client) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}

// GetFileInfo gets file metadata
func (c *Client) GetFileInfo(ctx context.Context, objectName string) (*FileInfo, error) {
	stat, err := c.client.StatObject(ctx, c.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}

	return &FileInfo{
		Key:          stat.Key,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		LastModified: stat.LastModified,
		ETag:         stat.ETag,
	}, nil
}

// DeleteFile deletes a file from storage
func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// ListFiles lists files with a prefix
func (c *Client) ListFiles(ctx context.Context, prefix string, maxKeys int) ([]FileInfo, error) {
	var files []FileInfo

	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
		MaxKeys:   maxKeys,
	}

	for object := range c.client.ListObjects(ctx, c.bucketName, opts) {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}

		files = append(files, FileInfo{
			Key:          object.Key,
			Size:         object.Size,
			ContentType:  object.ContentType,
			LastModified: object.LastModified,
			ETag:         object.ETag,
		})
	}

	return files, nil
}

// CopyFile copies a file within storage
func (c *Client) CopyFile(ctx context.Context, srcObject, destObject string) error {
	src := minio.CopySrcOptions{
		Bucket: c.bucketName,
		Object: srcObject,
	}

	dst := minio.CopyDestOptions{
		Bucket: c.bucketName,
		Object: destObject,
	}

	_, err := c.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// UploadResult represents upload result
type UploadResult struct {
	Bucket      string
	Key         string
	ETag        string
	Size        int64
	ContentType string
}

// FileInfo represents file metadata
type FileInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}
