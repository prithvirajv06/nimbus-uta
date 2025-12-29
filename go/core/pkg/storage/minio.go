package storage

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/prithvirajv06/nimbus-uta/go/core/config"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient(cfg config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Create bucket if not exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOClient{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (m *MinIOClient) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(
		ctx,
		m.bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	return err
}

func (m *MinIOClient) Download(ctx context.Context, objectName string) (*minio.Object, error) {
	return m.client.GetObject(ctx, m.bucket, objectName, minio.GetObjectOptions{})
}

func (m *MinIOClient) Delete(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}

func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (*url.URL, error) {
	return m.client.PresignedGetObject(ctx, m.bucket, objectName, expiry, nil)
}

func (m *MinIOClient) ListObjects(ctx context.Context, prefix string) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, m.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
}

func (m *MinIOClient) ObjectExists(ctx context.Context, objectName string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MinIOClient) CopyObject(ctx context.Context, srcObject, destObject string) error {
	src := minio.CopySrcOptions{
		Bucket: m.bucket,
		Object: srcObject,
	}
	dest := minio.CopyDestOptions{
		Bucket: m.bucket,
		Object: destObject,
	}
	_, err := m.client.CopyObject(ctx, dest, src)
	return err
}

// FileMetadata represents file information
type FileMetadata struct {
	Name         string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
}

func (m *MinIOClient) GetObjectInfo(ctx context.Context, objectName string) (*FileMetadata, error) {
	info, err := m.client.StatObject(ctx, m.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileMetadata{
		Name:         info.Key,
		Size:         info.Size,
		ContentType:  info.ContentType,
		LastModified: info.LastModified,
		ETag:         info.ETag,
	}, nil
}
