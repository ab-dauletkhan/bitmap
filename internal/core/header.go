package core

import (
	"encoding/binary"
	"fmt"
)

// BMPHeader defines the structure for the BMP file header.
// It contains basic information like file type, file size in bytes, and header size.
type BMPHeader struct {
	FileType        string // Should always be "BM" for valid BMP files.
	FileSizeInBytes uint32 // Total size of the BMP file in bytes.
	HeaderSize      uint16 // Size of the BMP header, typically 54 bytes.
}

// DIBHeader defines the structure for the DIB (Device Independent Bitmap) header.
// It stores detailed information about the image, such as dimensions and pixel format.
type DIBHeader struct {
	DibHeaderSize    uint32 // Size of the DIB header.
	WidthInPixels    int32  // Width of the image in pixels.
	HeightInPixels   int32  // Height of the image in pixels. Can be negative if stored upside-down.
	PixelSizeInBits  uint16 // Number of bits per pixel. Usually 24 for BMP.
	ImageSizeInBytes uint32 // Size of the raw bitmap data in bytes.
}

// BMPImage encapsulates both the BMP and DIB headers, along with the actual image data.
type BMPImage struct {
	bmpHeader *BMPHeader // BMP file header.
	dibHeader *DIBHeader // DIB header containing image information.
	Data      []byte     // Raw image data (pixel array).
}

// ParseHeader parses the BMP and DIB headers from a byte slice.
// The function checks if the file is a valid 24-bit uncompressed BMP and handles various errors.
//
// The function:
// - Verifies the file size and header integrity.
// - Ensures the file type is "BM".
// - Extracts the header fields and validates the dimensions and pixel size.
// - Checks the raw image data size matches the expected size.
//
// Returns:
// - BMPImage struct containing the parsed headers and image data.
// - Error if the BMP is invalid or corrupted.
func ParseHeader(b []byte) (*BMPImage, error) {
	// The BMP file must be at least 54 bytes (BMP and DIB headers combined).
	if len(b) < 54 {
		return nil, ErrInvalidBMP
	}

	// Parse the BMP header.
	header := &BMPHeader{}
	header.FileType = string(b[:2]) // The first two bytes should be "BM".
	if header.FileType != "BM" {
		return nil, ErrInvalidFileType
	}
	header.FileSizeInBytes = binary.LittleEndian.Uint32(b[2:6]) // Extract file size from the header.
	if int(header.FileSizeInBytes) != len(b) {
		return nil, ErrCorruptFile
	}
	header.HeaderSize = binary.LittleEndian.Uint16(b[10:14]) // Extract header size.

	// Parse the DIB header.
	dibheader := &DIBHeader{}
	dibheader.DibHeaderSize = binary.LittleEndian.Uint32(b[14:18]) // Extract DIB header size.
	if dibheader.DibHeaderSize < 40 {                              // Check if DIB header size is valid.
		return nil, ErrInvalidHeaderSize
	}
	dibheader.WidthInPixels = int32(binary.LittleEndian.Uint32(b[18:22]))  // Extract image width.
	dibheader.HeightInPixels = int32(binary.LittleEndian.Uint32(b[22:26])) // Extract image height.

	// Ensure that both width and height are positive. Height can be negative (stored upside-down).
	if dibheader.WidthInPixels <= 0 || dibheader.HeightInPixels == 0 {
		return nil, ErrNonPositiveDimensions
	}

	dibheader.PixelSizeInBits = binary.LittleEndian.Uint16(b[28:30]) // Extract bits per pixel.
	if dibheader.PixelSizeInBits != 24 {                             // Only 24-bit BMPs are supported.
		return nil, ErrUnsupportedFormat
	}

	dibheader.ImageSizeInBytes = binary.LittleEndian.Uint32(b[34:38]) // Extract image data size.

	// Calculate the expected image size in bytes based on dimensions and pixel size.
	bytesPerPixel := int(dibheader.PixelSizeInBits / 8)
	rowSize := ((int(dibheader.WidthInPixels)*bytesPerPixel + 3) / 4) * 4 // BMP rows are padded to a multiple of 4 bytes.
	expectedImageSize := rowSize * int(abs32(dibheader.HeightInPixels))   // Adjust for negative height if upside-down.

	// Ensure the actual image size matches the expected size.
	if expectedImageSize != int(dibheader.ImageSizeInBytes) {
		return nil, ErrInvalidImageData
	}

	// Return the parsed BMPImage struct containing headers and image data.
	return &BMPImage{
		bmpHeader: header,
		dibHeader: dibheader,
		Data:      b[54:], // Image data starts after the 54-byte header.
	}, nil
}

// PrintBMPHeaderInfo prints the BMP and DIB header information in a formatted style.
// It consolidates the BMP and DIB header details in one optimized print statement.
func PrintBMPHeaderInfo(image *BMPImage) {
	// Print all header information in one go using `fmt.Printf` for efficiency.
	fmt.Printf(`BMP Header:
- FileType %s
- FileSizeInBytes %d
- HeaderSize %d
DIB Header:
- DibHeaderSize %d
- WidthInPixels %d
- HeightInPixels %d
- PixelSizeInBits %d
- ImageSizeInBytes %d
`,
		image.bmpHeader.FileType,
		image.bmpHeader.FileSizeInBytes,
		image.bmpHeader.HeaderSize,
		image.dibHeader.DibHeaderSize,
		image.dibHeader.WidthInPixels,
		image.dibHeader.HeightInPixels,
		image.dibHeader.PixelSizeInBits,
		image.dibHeader.ImageSizeInBytes,
	)
}

// abs32 returns the absolute value of a signed 32-bit integer.
// It is used to handle potential negative image height, which indicates upside-down storage.
func abs32(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}
