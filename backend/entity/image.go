package entity

type UploadImageRequest struct {
	Image []byte
}

type UploadImageResponse struct {
	ImageId string
}
