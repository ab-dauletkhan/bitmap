package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ab-dauletkhan/bitmap/internal/utils"
)

// CropInfo holds the parameters needed for cropping an image.
type CropInfo struct {
	OffsetX int // The x-coordinate of the top-left corner of the crop area.
	OffsetY int // The y-coordinate of the top-left corner of the crop area.
	Width   int // The width of the crop area.
	Height  int // The height of the crop area.
}

// parseCropInfo parses the crop string format into CropInfo.
// The crop string can contain either two values (OffsetX, OffsetY)
// or four values (OffsetX, OffsetY, Width, Height).
// It returns a CropInfo struct and an error if parsing fails.
func parseCropInfo(cropStr string) (CropInfo, error) {
	info := strings.Split(cropStr, "-")
	var cropInfo CropInfo

	if len(info) != 2 && len(info) != 4 {
		return cropInfo, fmt.Errorf("crop option must have either 2 or 4 values")
	}

	var err error
	cropInfo.OffsetX, err = strconv.Atoi(info[0])
	if err != nil || cropInfo.OffsetX <= 0 {
		return cropInfo, fmt.Errorf("invalid OffsetX value")
	}

	cropInfo.OffsetY, err = strconv.Atoi(info[1])
	if err != nil || cropInfo.OffsetY <= 0 {
		return cropInfo, fmt.Errorf("invalid OffsetY value")
	}

	if len(info) == 4 {
		cropInfo.Width, err = strconv.Atoi(info[2])
		if err != nil || cropInfo.Width <= 0 {
			return cropInfo, fmt.Errorf("invalid Width value")
		}

		cropInfo.Height, err = strconv.Atoi(info[3])
		if err != nil || cropInfo.Height <= 0 {
			return cropInfo, fmt.Errorf("invalid Height value")
		}
	}

	return cropInfo, nil
}

// Crop modifies the BMPImage to only include the specified area defined by CropInfo.
// It adjusts the image dimensions and discards pixels outside the crop area.
// The crop area is defined by OffsetX and OffsetY as the top-left corner,
// with the specified Width and Height. An error is returned if the crop
// area exceeds the image boundaries or if it results in invalid dimensions.
func Crop(image *BMPImage, opts CropInfo) error {

	if opts.OffsetX+opts.Width > len(image.Data[0]) || opts.OffsetY+opts.Height > len(image.Data) {
		return fmt.Errorf("crop area exceeds image boundaries")
	}

	if opts.Width == 0 {
		opts.Width = int(image.InfoHeader.Width) - opts.OffsetX
	}

	if opts.Height == 0 {
		opts.Height = utils.Abs(int(image.InfoHeader.Height)) - opts.OffsetY
	}

	croppedData := make([][]Pixel, opts.Height)
	for i := range croppedData {
		croppedData[i] = make([]Pixel, opts.Width)
	}

	for y := 0; y < opts.Height; y++ {
		for x := 0; x < opts.Width; x++ {
			croppedData[y][x] = image.Data[opts.OffsetY+y][opts.OffsetX+x]
		}
	}

	image.Data = croppedData
	image.InfoHeader.Width = int32(opts.Width)
	image.InfoHeader.Height = int32(opts.Height)

	return nil
}
