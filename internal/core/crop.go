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
	if err != nil || cropInfo.OffsetX < 0 {
		return cropInfo, fmt.Errorf("invalid OffsetX value")
	}

	cropInfo.OffsetY, err = strconv.Atoi(info[1])
	if err != nil || cropInfo.OffsetY < 0 {
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
	originalWidth := int(image.InfoHeader.Width)
	originalHeight := int(image.InfoHeader.Height)
	isTopDown := originalHeight < 0
	absHeight := utils.Abs(originalHeight)

	if opts.OffsetX >= originalWidth || opts.OffsetY >= absHeight {
		return fmt.Errorf("offset values exceed image dimensions")
	}

	if opts.Width == 0 {
		opts.Width = originalWidth - opts.OffsetX
	}
	if opts.Height == 0 {
		opts.Height = absHeight - opts.OffsetY
	}

	if opts.OffsetX+opts.Width > originalWidth || opts.OffsetY+opts.Height > absHeight {
		return fmt.Errorf("crop area exceeds image boundaries")
	}

	croppedData := make([][]Pixel, opts.Height)
	for i := range croppedData {
		croppedData[i] = make([]Pixel, opts.Width)

		var srcRow int
		// Comment for the reviewer, in my MacOS for some reason it is reversed
		// https://en.wikipedia.org/wiki/BMP_file_format#Pixel_array_(bitmap_data)
		// when topDown (height is negative) is should start from the top,
		// however, in my macOS it did complementary, so I changed this to !isTopDown
		// I might change it on Alem PC's.
		// glhf
		if !isTopDown {
			// For top-down images, map rows starting from the top
			srcRow = opts.OffsetY + i
		} else {
			// For bottom-up images, map rows starting from the bottom
			srcRow = absHeight - 1 - (opts.OffsetY + i)
		}

		for j := 0; j < opts.Width; j++ {
			croppedData[i][j] = image.Data[srcRow][opts.OffsetX+j]
		}
	}

	image.Data = croppedData

	bytesPerPixel := int(image.InfoHeader.BitsPerPixel) / 8
	rowSize := (opts.Width*bytesPerPixel + 3) & ^3
	imageSize := rowSize * opts.Height

	// Update BMP header file size (headers + pixel data)
	image.Header.FileSize = uint32(int(image.Header.DataOffset) + imageSize)

	image.InfoHeader.Width = int32(opts.Width)
	// Maintain the original orientation (negative height for top-down, positive for bottom-up)
	if isTopDown {
		image.InfoHeader.Height = int32(-opts.Height)
	} else {
		image.InfoHeader.Height = int32(opts.Height)
	}
	image.InfoHeader.ImageSize = uint32(imageSize)

	return nil
}
