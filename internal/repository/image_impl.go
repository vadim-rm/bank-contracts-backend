package repository

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type ImageConfig struct {
	BucketName string
	BaseUrl    string
}

type ImageImpl struct {
	client *minio.Client
	config ImageConfig
}

func NewImageImpl(
	client *minio.Client,
	config ImageConfig,
) *ImageImpl {
	return &ImageImpl{
		client: client,
		config: config,
	}
}

func (r *ImageImpl) Upload(ctx context.Context, input UploadImageInput) (string, error) {
	_, err := r.client.PutObject(ctx, r.config.BucketName, input.Name, input.Image,
		input.Size, minio.PutObjectOptions{
			ContentType: input.ContentType,
		})

	return fmt.Sprintf("%s/%s", r.config.BaseUrl, input.Name), err
}
