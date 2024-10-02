package core

type CropInfo struct {
	OffsetX int
	OffsetY int
	Width   int
	Height  int
}

func Crop(image *BMPImage, opts CropInfo) error {
	return nil
}
