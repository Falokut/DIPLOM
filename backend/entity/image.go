package entity

type UploadImageRequest struct {
	Image []byte
}

type ReplaceImageRequest struct {
	ImageData []byte
}

type UploadImageResponse struct {
	ImageId string
}
