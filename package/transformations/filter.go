package transformations

import "github.com/ab-dauletkhan/bitmap/internal/core"

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
		applyColor(image, pixelate)
	case "blur":
		applyColor(image, blur)
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
