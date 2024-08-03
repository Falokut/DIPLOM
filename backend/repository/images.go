package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"dish_as_a_service/entity"

	"github.com/pkg/errors"
)

type Image struct {
	cli          *http.Client
	baseImageUrl string
	serviceAddr  string
}

func NewImage(cli *http.Client, serviceAddr, baseImageUrl string) Image {
	return Image{
		cli:          cli,
		baseImageUrl: baseImageUrl,
		serviceAddr:  serviceAddr,
	}
}

const (
	uploadImageEndpoint = "%s/api/images-service/image/%s"
	deleteImageEndpoint = "%s/api/images-service/image/%s/%s"
)

const (
	JsonContentType = "json"
)

func (r Image) UploadImage(ctx context.Context, category string, image []byte) (string, error) {
	url := fmt.Sprintf(uploadImageEndpoint, r.serviceAddr, category)

	req, err := makeRequest(ctx, http.MethodPost, url, entity.UploadImageRequest{
		Image: image,
	})
	if err != nil {
		return "", errors.WithMessage(err, "make request")
	}

	resp, err := r.cli.Do(req)
	if err != nil {
		return "", errors.WithMessage(err, "send request")
	}
	defer resp.Body.Close()

	var uploadResp entity.UploadImageResponse
	err = json.NewDecoder(resp.Body).Decode(&uploadResp)
	if err != nil {
		return "", errors.WithMessage(err, "decode response")
	}

	return uploadResp.ImageId, nil
}

func (r Image) DeleteImage(ctx context.Context, category, imageId string) error {
	url := fmt.Sprintf(deleteImageEndpoint, r.serviceAddr, category, imageId)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return errors.WithMessage(err, "make request")
	}
	resp, err := r.cli.Do(req)
	if err != nil {
		return errors.WithMessage(err, "send request")
	}
	resp.Body.Close()

	return nil
}

func (r Image) GetImageUrl(category, imageId string) string {
	if imageId == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", r.baseImageUrl, category, imageId)
}

func makeRequest(ctx context.Context, method, url string, req any) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, errors.WithMessage(err, "marshal request")
	}

	reqReader := bytes.NewReader(body)
	request, err := http.NewRequestWithContext(ctx, method, url, reqReader)
	if err != nil {
		return nil, errors.WithMessage(err, "new request")
	}
	return request, nil
}
