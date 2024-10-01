package repository

import (
	"context"
	"io"
)

type Image interface {
	Upload(ctx context.Context, input UploadImageInput) (string, error)
}

type UploadImageInput struct {
	Image       io.Reader
	Size        int64
	Name        string
	ContentType string
}
