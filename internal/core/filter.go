package core

const (
	blue = iota
	green
	red
	grayscale
	negative
	pixelate
	blur
)

// Filter applies a specified filter to the given BMPImage.
// Supported filters: "blue", "green", "red", "grayscale", "negative", "pixelate", and "blur".
// The "pixelate" filter uses a default block size of 20 pixels.
// The "blur" filter applies a blur with a default radius of 50 pixels.
func Filter(image *BMPImage, filter string) {
	switch filter {
	case "blue":
		applyColor(image, blue)
	case "green":
		applyColor(image, green)
	case "red":
		applyColor(image, red)
	case "grayscale":
		applyColor(image, grayscale)
	case "negative":
		applyColor(image, negative)
	case "pixelate":
		applyPixelate(image, 50)
	case "blur":
		applyBlur(image, 20)
	}
}

// applyColor applies a color-based filter to the BMPImage data.
// The color argument determines which channel to keep: blue, green, red, grayscale, or negative.
func applyColor(image *BMPImage, color int) {
	h := len(image.Data)
	w := len(image.Data[0])

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			switch color {
			case blue:
				image.Data[y][x].Green = 0
				image.Data[y][x].Red = 0
			case green:
				image.Data[y][x].Blue = 0
				image.Data[y][x].Red = 0
			case red:
				image.Data[y][x].Blue = 0
				image.Data[y][x].Green = 0
			case grayscale:
				// Convert pixel to grayscale using standard luminance calculation
				gray := byte(float64(image.Data[y][x].Red)*0.2126 + float64(image.Data[y][x].Green)*0.7152 + float64(image.Data[y][x].Blue)*0.0722)
				image.Data[y][x].Blue = gray
				image.Data[y][x].Green = gray
				image.Data[y][x].Red = gray
			case negative:
				// Invert the color values for a negative effect
				image.Data[y][x].Blue = 255 - image.Data[y][x].Blue
				image.Data[y][x].Green = 255 - image.Data[y][x].Green
				image.Data[y][x].Red = 255 - image.Data[y][x].Red
			}
		}
	}
}

// applyPixelate applies a pixelation effect to the BMPImage data.
// The blocksize argument specifies the size of the pixelation blocks.
func applyPixelate(image *BMPImage, blocksize int) {
	h := len(image.Data)
	w := len(image.Data[0])

	for y := 0; y < h; y += blocksize {
		for x := 0; x < w; x += blocksize {
			colorPixel := avgColorBlock(image, x, y, blocksize)
			fillBlock(image, x, y, blocksize, colorPixel)
		}
	}
}

// avgColorBlock calculates the average color of a block starting at (startX, startY).
// The blocksize determines the dimensions of the block to average.
func avgColorBlock(image *BMPImage, startX, startY, blocksize int) Pixel {
	var rSum, gSum, bSum, cnt uint32
	h := len(image.Data)
	w := len(image.Data[0])

	for y := startY; y < startY+blocksize && y < h; y++ {
		for x := startX; x < startX+blocksize && x < w; x++ {
			rSum += uint32(image.Data[y][x].Red)
			gSum += uint32(image.Data[y][x].Green)
			bSum += uint32(image.Data[y][x].Blue)
			cnt++
		}
	}

	// Return the average color for the block
	return Pixel{
		Red:   byte(rSum / cnt),
		Green: byte(gSum / cnt),
		Blue:  byte(bSum / cnt),
	}
}

// fillBlock fills a block of the image with a given color.
// The block starts at (startX, startY) and has the size of blocksize.
func fillBlock(image *BMPImage, startX, startY, blocksize int, colorPixel Pixel) {
	h := len(image.Data)
	w := len(image.Data[0])

	for y := startY; y < startY+blocksize && y < h; y++ {
		for x := startX; x < startX+blocksize && x < w; x++ {
			image.Data[y][x] = colorPixel
		}
	}
}

// applyBlur applies a basic box blur to the given BMPImage.
// The blurRadius defines the size of the neighborhood around each pixel used for averaging.
// A larger blurRadius results in a more pronounced blur effect.
func applyBlur(image *BMPImage, blurRadius int) {
	width := int(image.InfoHeader.Width)
	height := int(image.InfoHeader.Height)

	// Create a copy of the original image data to store blurred results.
	blurredData := make([][]Pixel, height)
	for i := range blurredData {
		blurredData[i] = make([]Pixel, width)
	}

	// Iterate over each pixel in the image.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Initialize accumulators for each color channel.
			var redSum, greenSum, blueSum, count int

			// Iterate over the neighborhood of the current pixel.
			for ky := -blurRadius; ky <= blurRadius; ky++ {
				for kx := -blurRadius; kx <= blurRadius; kx++ {
					// Calculate the neighboring pixel's coordinates.
					nx := x + kx
					ny := y + ky

					// Ensure the neighboring pixel is within bounds.
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						// Accumulate the color values.
						pixel := image.Data[ny][nx]
						redSum += int(pixel.Red)
						greenSum += int(pixel.Green)
						blueSum += int(pixel.Blue)
						count++
					}
				}
			}

			// Calculate the average color values for the pixel.
			blurredData[y][x] = Pixel{
				Red:   byte(redSum / count),
				Green: byte(greenSum / count),
				Blue:  byte(blueSum / count),
			}
		}
	}

	// Replace the original image data with the blurred version.
	image.Data = blurredData
}

// // applyBlur applies a blur effect to the BMPImage using a blur radius.
// // The function divides the work between multiple goroutines based on the number of CPU cores available.
// func applyBlur(image *BMPImage, blurRadius int) {
// 	startTime := time.Now()

// 	h := len(image.Data)
// 	w := len(image.Data[0])
// 	k := 2*blurRadius + 1
// 	weight := 1.0 / float64(k)

// 	// Determine the number of goroutines to use
// 	numGoroutines := runtime.NumCPU()
// 	chunkSize := h / numGoroutines
// 	if chunkSize < 1 {
// 		chunkSize = 1
// 		numGoroutines = h
// 	}

// 	var wg sync.WaitGroup

// 	// Horizontal pass of the blur
// 	for i := 0; i < numGoroutines; i++ {
// 		wg.Add(1)
// 		go func(startY, endY int) {
// 			defer wg.Done()
// 			for y := startY; y < endY; y++ {
// 				tempRow := make([]float64, w*3)
// 				for x := 0; x < w; x++ {
// 					var rSum, gSum, bSum float64
// 					for i := -blurRadius; i <= blurRadius; i++ {
// 						ix := x + i
// 						if ix >= 0 && ix < w {
// 							rSum += float64(image.Data[y][ix].Red)
// 							gSum += float64(image.Data[y][ix].Green)
// 							bSum += float64(image.Data[y][ix].Blue)
// 						}
// 					}
// 					tempRow[x*3] = rSum * weight
// 					tempRow[x*3+1] = gSum * weight
// 					tempRow[x*3+2] = bSum * weight
// 				}
// 				for x := 0; x < w; x++ {
// 					image.Data[y][x].Red = byte(clamp(tempRow[x*3]))
// 					image.Data[y][x].Green = byte(clamp(tempRow[x*3+1]))
// 					image.Data[y][x].Blue = byte(clamp(tempRow[x*3+2]))
// 				}
// 			}
// 		}(i*chunkSize, min((i+1)*chunkSize, h))
// 	}
// 	wg.Wait()

// 	// Vertical pass of the blur
// 	for i := 0; i < numGoroutines; i++ {
// 		wg.Add(1)
// 		go func(startX, endX int) {
// 			defer wg.Done()
// 			for x := startX; x < endX; x++ {
// 				tempCol := make([]float64, h*3)
// 				for y := 0; y < h; y++ {
// 					var rSum, gSum, bSum float64
// 					for i := -blurRadius; i <= blurRadius; i++ {
// 						iy := y + i
// 						if iy >= 0 && iy < h {
// 							rSum += float64(image.Data[iy][x].Red)
// 							gSum += float64(image.Data[iy][x].Green)
// 							bSum += float64(image.Data[iy][x].Blue)
// 						}
// 					}
// 					tempCol[y*3] = rSum * weight
// 					tempCol[y*3+1] = gSum * weight
// 					tempCol[y*3+2] = bSum * weight
// 				}
// 				for y := 0; y < h; y++ {
// 					image.Data[y][x].Red = byte(clamp(tempCol[y*3]))
// 					image.Data[y][x].Green = byte(clamp(tempCol[y*3+1]))
// 					image.Data[y][x].Blue = byte(clamp(tempCol[y*3+2]))
// 				}
// 			}
// 		}(i*chunkSize, min((i+1)*chunkSize, w))
// 	}
// 	wg.Wait()

// 	// Output the time taken to apply the blur
// 	elapsedTime := time.Since(startTime)
// 	fmt.Printf("Parallel blur operation took %v\n", elapsedTime)
// }

// // clamp ensures that the value is within the range 0 to 255.
// func clamp(v float64) float64 {
// 	if v > 255 {
// 		return 255
// 	}
// 	if v < 0 {
// 		return 0
// 	}
// 	return v
// }

// // min returns the smaller of two integers.
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }
