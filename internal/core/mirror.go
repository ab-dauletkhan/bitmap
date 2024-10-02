package core

func MirrorImage(image *BMPImage, direction string) *BMPImage {
	h := len(image.Data)
	w := len(image.Data[0])

	switch direction {
	case "horizontal":
		for y := 0; y < h; y++ {
			for x := 0; x < w/2; x++ {
				image.Data[y][x], image.Data[y][w-x-1] = image.Data[y][w-x-1], image.Data[y][x]
			}
		}
	case "vertical":
		for y := 0; y < h/2; y++ {
			for x := 0; x < w; x++ {
				image.Data[y][x], image.Data[h-y-1][x] = image.Data[h-y-1][x], image.Data[y][x]
			}
		}
	}
	return image
}
