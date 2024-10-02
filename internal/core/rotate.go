package core

func Rotate(image *BMPImage, direction int) {
	h := len(image.Data)
	w := len(image.Data[0])

	rotatedData := make([][]Pixel, w)
	for i := 0; i < w; i++ {
		rotatedData[i] = make([]Pixel, h)
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
