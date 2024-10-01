package transformations

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ab-dauletkhan/bitmap/internal/core"
)

const (
	blue = iota
	green
	red
	grayscale
	negative
	pixelate
	blur
)

func Filter(image *core.BMPImage, filter string) {
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
		applyPixelate(image, 20)
	case "blur":
		applyBlur(image, 50)
	}
}

func applyColor(image *core.BMPImage, color int) {
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
				gray := byte(float64(image.Data[y][x].Red)*0.2126 + float64(image.Data[y][x].Green)*0.7152 + float64(image.Data[y][x].Blue)*0.0722)
				image.Data[y][x].Blue = gray
				image.Data[y][x].Green = gray
				image.Data[y][x].Red = gray
			case negative:
				image.Data[y][x].Blue = 255 - image.Data[y][x].Blue
				image.Data[y][x].Green = 255 - image.Data[y][x].Green
				image.Data[y][x].Red = 255 - image.Data[y][x].Red
			}
		}
	}
}

func applyPixelate(image *core.BMPImage, blocksize int) {
	h := len(image.Data)
	w := len(image.Data[0])

	for y := 0; y < h; y += blocksize {
		for x := 0; x < w; x += blocksize {
			colorPixel := avgColorBlock(image, x, y, blocksize)
			fillBlock(image, x, y, blocksize, colorPixel)
		}
	}
}

func avgColorBlock(image *core.BMPImage, startX, startY, blocksize int) core.Pixel {
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

	return core.Pixel{
		Red:   byte(rSum / cnt),
		Green: byte(gSum / cnt),
		Blue:  byte(bSum / cnt),
	}
}

func fillBlock(image *core.BMPImage, startX, startY, blocksize int, colorPixel core.Pixel) {
	h := len(image.Data)
	w := len(image.Data[0])

	for y := startY; y < startY+blocksize && y < h; y++ {
		for x := startX; x < startX+blocksize && x < w; x++ {
			image.Data[y][x] = colorPixel
		}
	}
}

func applyBlur(image *core.BMPImage, blurRadius int) {
	startTime := time.Now()

	h := len(image.Data)
	w := len(image.Data[0])
	k := 2*blurRadius + 1
	weight := 1.0 / float64(k)

	// Determine the number of goroutines to use
	numGoroutines := runtime.NumCPU()
	chunkSize := h / numGoroutines
	if chunkSize < 1 {
		chunkSize = 1
		numGoroutines = h
	}

	var wg sync.WaitGroup

	// Horizontal pass
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			for y := startY; y < endY; y++ {
				tempRow := make([]float64, w*3)
				for x := 0; x < w; x++ {
					var rSum, gSum, bSum float64
					for i := -blurRadius; i <= blurRadius; i++ {
						ix := x + i
						if ix >= 0 && ix < w {
							rSum += float64(image.Data[y][ix].Red)
							gSum += float64(image.Data[y][ix].Green)
							bSum += float64(image.Data[y][ix].Blue)
						}
					}
					tempRow[x*3] = rSum * weight
					tempRow[x*3+1] = gSum * weight
					tempRow[x*3+2] = bSum * weight
				}
				for x := 0; x < w; x++ {
					image.Data[y][x].Red = byte(clamp(tempRow[x*3]))
					image.Data[y][x].Green = byte(clamp(tempRow[x*3+1]))
					image.Data[y][x].Blue = byte(clamp(tempRow[x*3+2]))
				}
			}
		}(i*chunkSize, min((i+1)*chunkSize, h))
	}
	wg.Wait()

	// Vertical pass
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startX, endX int) {
			defer wg.Done()
			for x := startX; x < endX; x++ {
				tempCol := make([]float64, h*3)
				for y := 0; y < h; y++ {
					var rSum, gSum, bSum float64
					for i := -blurRadius; i <= blurRadius; i++ {
						iy := y + i
						if iy >= 0 && iy < h {
							rSum += float64(image.Data[iy][x].Red)
							gSum += float64(image.Data[iy][x].Green)
							bSum += float64(image.Data[iy][x].Blue)
						}
					}
					tempCol[y*3] = rSum * weight
					tempCol[y*3+1] = gSum * weight
					tempCol[y*3+2] = bSum * weight
				}
				for y := 0; y < h; y++ {
					image.Data[y][x].Red = byte(clamp(tempCol[y*3]))
					image.Data[y][x].Green = byte(clamp(tempCol[y*3+1]))
					image.Data[y][x].Blue = byte(clamp(tempCol[y*3+2]))
				}
			}
		}(i*chunkSize, min((i+1)*chunkSize, w))
	}
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Parallel blur operation took %v\n", elapsedTime)
}

func clamp(v float64) float64 {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
