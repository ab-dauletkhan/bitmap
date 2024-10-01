package transformations

import "github.com/ab-dauletkhan/bitmap/internal/core"

func Rotate(image *core.BMPImage, direction int) {
	h := len(image.Data)
	w := len(image.Data[0])

	rotatedData := make([][]core.Pixel, w)
	for i := 0; i < w; i++ {
		rotatedData[i] = make([]core.Pixel, h)
		for j := 0; j < h; j++ {
			if direction == -1 { // left
				rotatedData[i][j] = image.Data[h-1-j][i]
			} else { // right
				rotatedData[i][j] = image.Data[j][w-1-i]
			}
		}
	}
	image.InfoHeader.Height = int32(w)
	image.InfoHeader.Width = int32(h)
	image.Data = rotatedData
}
