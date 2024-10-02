package core

// MirrorImage mirrors the BMPImage either horizontally or vertically based on the given direction.
// The "horizontal" direction swaps pixels from left to right, while the "vertical" direction flips
// the image by adjusting the pixel positions from top to bottom. In the vertical case, the height
// of the DIB header is inverted to reflect the change.
func MirrorImage(image *BMPImage, direction string) {
	h := len(image.Data)
	w := len(image.Data[0])

	switch direction {
	case "horizontal":
		// Mirror the image horizontally by swapping pixels along each row
		for y := 0; y < h; y++ {
			for x := 0; x < w/2; x++ {
				image.Data[y][x], image.Data[y][w-x-1] = image.Data[y][w-x-1], image.Data[y][x]
			}
		}
	case "vertical":
		// Mirror the image vertically by flipping the height in the DIB header
		// This changes the image orientation without manually swapping pixel data
		// image.InfoHeader.Height *= -1 is used to indicate that the image is mirrored vertically
		image.InfoHeader.Height *= -1

		// or we could swap them as well with O(n) time complexity (n is the number of the bytes in bmp)
		// for y := 0; y < h/2; y++ {
		// 	for x := 0; x < w; x++ {
		// 		image.Data[y][x], image.Data[h-y-1][x] = image.Data[h-y-1][x], image.Data[y][x]
		// 	}
		// }
	}
}
