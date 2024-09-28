package repository

import (
	"context"
	"dish_as_a_service/entity"
	"fmt"
	"github.com/Falokut/go-kit/http/client"
	"github.com/pkg/errors"
)

type Image struct {
	cli          *client.Client
	baseImageUrl string
}

func NewImage(cli *client.Client, baseImageUrl string) Image {
	return Image{
		cli:          cli,
		baseImageUrl: baseImageUrl,
	}
}

const (
	uploadImageEndpoint  = "image/%s"
	deleteImageEndpoint  = "image/%s/%s"
	replaceImageEndpoint = "image/%s/%s"
)

func (r Image) UploadImage(ctx context.Context, category string, image []byte) (string, error) {
	url := fmt.Sprintf(uploadImageEndpoint, category)
	req := entity.UploadImageRequest{Image: image}
	var resp entity.UploadImageResponse
	_, err := r.cli.Post(url).
		JsonRequestBody(req).
		JsonResponseBody(&resp).
		StatusCodeToError().
		Do(ctx)
	if err != nil {
		return "", errors.WithMessage(err, "send upload image request")
	}

	return resp.ImageId, nil
}

func (r Image) DeleteImage(ctx context.Context, category, imageId string) error {
	url := fmt.Sprintf(deleteImageEndpoint, category, imageId)
	_, err := r.cli.Delete(url).
		StatusCodeToError().
		Do(ctx)
	if err != nil {
		return errors.WithMessage(err, "send delete image request")
	}

	return nil
}

func (r Image) GetImageUrl(category, imageId string) string {
	if imageId == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", r.baseImageUrl, category, imageId)
}
