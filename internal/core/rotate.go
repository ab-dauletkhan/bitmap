package core

// Rotate rotates the BMPImage 90 degrees to the left or right based on the direction.
// The direction value determines the rotation:
// - A value of -1 rotates the image 90 degrees to the left (counterclockwise).
// - Any other value rotates the image 90 degrees to the right (clockwise).
// The function updates the image's width and height in the DIB header after rotation.
func Rotate(image *BMPImage, direction int) {
	h := len(image.Data)
	w := len(image.Data[0])

	// Create a new 2D slice for the rotated image data with swapped width and height
	rotatedData := make([][]Pixel, w)
	for i := 0; i < w; i++ {
		rotatedData[i] = make([]Pixel, h)
		for j := 0; j < h; j++ {
			if direction == -1 { // to the left (counterclockwise)
				rotatedData[i][j] = image.Data[h-1-j][i]
			} else { // to the right (clockwise)
				rotatedData[i][j] = image.Data[j][w-1-i]
			}
		}
	}

	// Update the BMP image headers to reflect the new dimensions after rotation
	image.InfoHeader.Height = int32(w)
	image.InfoHeader.Width = int32(h)

	// Replace the original pixel data with the rotated data
	image.Data = rotatedData
}
