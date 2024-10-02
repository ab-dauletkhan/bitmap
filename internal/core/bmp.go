package core

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/ab-dauletkhan/bitmap/internal/utils"
)

// BMPHeader defines the structure for the BMP file header.
// It contains basic information about the file type, size, and data offset.
type BMPHeader struct {
	Signature  [2]byte // Should always be "BM" for valid BMP files.
	FileSize   uint32  // Total size of the BMP file in bytes.
	Reserved   uint32  // Reserved, must be 0
	DataOffset uint32  // Offset to the start of pixel data
}

// DIBHeader defines the structure for the DIB (Device Independent Bitmap) header.
// It stores detailed information about the image, such as dimensions and color format.
type DIBHeader struct {
	Size            uint32 // Size of the DIB header
	Width           int32  // Width of the image in pixels
	Height          int32  // Height of the image in pixels
	Planes          uint16 // Number of color planes (must be 1)
	BitsPerPixel    uint16 // Bits per pixel
	Compression     uint32 // Compression method used
	ImageSize       uint32 // Size of the raw bitmap data
	XPixelsPerMeter int32  // Horizontal resolution of the image
	YPixelsPerMeter int32  // Vertical resolution of the image
	ColorsUsed      uint32 // Number of colors in the color palette
	ColorsImportant uint32 // Number of important colors used
}

// Pixel represents a single pixel in the BMP image with BGR channels.
type Pixel struct {
	Blue  byte
	Green byte
	Red   byte
}

// BMPImage encapsulates both the BMP and DIB headers, along with the actual image data.
type BMPImage struct {
	Header     BMPHeader
	InfoHeader DIBHeader
	Data       [][]Pixel
}

// ParseBMP parses a BMP file from a byte slice and returns a BMPImage struct.
// It performs various checks to ensure the validity and supported format of the BMP file.
//
// The function:
// - Verifies the file size and header integrity.
// - Ensures the file type is "BM".
// - Parses both BMP and DIB headers.
// - Validates header information including dimensions, bit depth, and compression.
// - Checks the raw image data size against the calculated expected size.
//
// Returns:
// - *BMPImage: A pointer to the parsed BMPImage struct.
// - error: An error if the BMP is invalid, unsupported, or corrupted.
func ParseBMP(b []byte) (*BMPImage, error) {
	if len(b) < 54 {
		return nil, ErrInvalidBMP
	}

	bmp := &BMPImage{}

	// Parse BMP Header
	bmp.Header.Signature = [2]byte{b[0], b[1]}
	if string(bmp.Header.Signature[:]) != "BM" {
		return nil, ErrInvalidFileType
	}
	bmp.Header.FileSize = binary.LittleEndian.Uint32(b[2:6])
	bmp.Header.Reserved = binary.LittleEndian.Uint32(b[6:10])
	bmp.Header.DataOffset = binary.LittleEndian.Uint32(b[10:14])

	// Parse DIB Header
	bmp.InfoHeader.Size = binary.LittleEndian.Uint32(b[14:18])
	if bmp.InfoHeader.Size < 40 {
		return nil, ErrInvalidHeaderSize
	}
	bmp.InfoHeader.Width = int32(binary.LittleEndian.Uint32(b[18:22]))
	bmp.InfoHeader.Height = int32(binary.LittleEndian.Uint32(b[22:26]))
	bmp.InfoHeader.Planes = binary.LittleEndian.Uint16(b[26:28])
	bmp.InfoHeader.BitsPerPixel = binary.LittleEndian.Uint16(b[28:30])
	bmp.InfoHeader.Compression = binary.LittleEndian.Uint32(b[30:34])
	bmp.InfoHeader.ImageSize = binary.LittleEndian.Uint32(b[34:38])
	bmp.InfoHeader.XPixelsPerMeter = int32(binary.LittleEndian.Uint32(b[38:42]))
	bmp.InfoHeader.YPixelsPerMeter = int32(binary.LittleEndian.Uint32(b[42:46]))
	bmp.InfoHeader.ColorsUsed = binary.LittleEndian.Uint32(b[46:50])
	bmp.InfoHeader.ColorsImportant = binary.LittleEndian.Uint32(b[50:54])

	// Validate header information
	if err := validateHeaders(bmp, len(b)); err != nil {
		return nil, err
	}

	// Set pixel data
	h := utils.Abs(int(bmp.InfoHeader.Height))
	w := int(bmp.InfoHeader.Width)
	bytesPerPixel := int(bmp.InfoHeader.BitsPerPixel) / 8
	rowSize := w * bytesPerPixel
	dataOffset := int(bmp.Header.DataOffset)
	bmp.Data = make([][]Pixel, h)

	for y := 0; y < h; y++ {
		bmp.Data[y] = make([]Pixel, w)
		for x := 0; x < w; x++ {
			pixelOffset := dataOffset + y*rowSize + x*bytesPerPixel
			bmp.Data[y][x] = Pixel{
				Blue:  b[pixelOffset],
				Green: b[pixelOffset+1],
				Red:   b[pixelOffset+2],
			}
		}
	}

	return bmp, nil
}

// validateHeaders performs various checks on the BMP and DIB headers to ensure
// the BMP file is valid and supported. It checks for correct file size, positive
// dimensions, supported bit depth, and uncompressed format. It also validates
// the image size against the calculated expected size.
//
// Parameters:
// - bmp: A pointer to the BMPImage struct containing the headers to validate.
// - fileSize: The actual size of the BMP file in bytes.
//
// Returns:
// - error: An error if any validation check fails, or nil if all checks pass.
func validateHeaders(bmp *BMPImage, fileSize int) error {
	if bmp.Header.FileSize != uint32(fileSize) {
		return ErrCorruptFile
	}
	if bmp.InfoHeader.Width <= 0 || bmp.InfoHeader.Height == 0 {
		return ErrNonPositiveDimensions
	}
	if bmp.InfoHeader.Planes != 1 {
		return ErrUnsupportedFormat
	}
	if bmp.InfoHeader.BitsPerPixel != 24 {
		return ErrUnsupportedFormat
	}
	if bmp.InfoHeader.Compression != 0 {
		return ErrUnsupportedCompression
	}

	// Validate image size
	widthInBytes := uint32(math.Abs(float64(bmp.InfoHeader.Width)) * float64(bmp.InfoHeader.BitsPerPixel) / 8)
	paddedWidth := (widthInBytes + 3) & ^uint32(3) // Round up to nearest multiple of 4
	expectedImageSize := (paddedWidth*uint32(math.Abs(float64(bmp.InfoHeader.Height))) + 3) & ^uint32(3)

	if bmp.InfoHeader.ImageSize != expectedImageSize {
		return ErrInvalidImageData
	}

	return nil
}

// SerializeBMP converts a BMPImage struct
// into a byte slice representing the complete BMP file.
// It handles the BMP and DIB headers, accounts for row padding,
// and properly organizes the pixel data.
func SerializeBMP(image *BMPImage) []byte {
	// Calculate sizes and offsets
	headerSize := int(image.Header.DataOffset)
	width := int(image.InfoHeader.Width)
	height := utils.Abs(int(image.InfoHeader.Height)) // Handle top-down BMPs
	bytesPerPixel := int(image.InfoHeader.BitsPerPixel) / 8
	rowSize := (width*bytesPerPixel + 3) & ^3 // 4-byte alignment
	dataSize := rowSize * height
	totalSize := headerSize + dataSize

	// Pre-allocate a byte slice for the entire BMP file
	data := make([]byte, totalSize)

	// Serialize BMP Header
	binary.LittleEndian.PutUint16(data[0:2], uint16(image.Header.Signature[0])|uint16(image.Header.Signature[1])<<8)
	binary.LittleEndian.PutUint32(data[2:6], image.Header.FileSize)
	binary.LittleEndian.PutUint32(data[6:10], image.Header.Reserved)
	binary.LittleEndian.PutUint32(data[10:14], image.Header.DataOffset)

	// Serialize DIB Header
	binary.LittleEndian.PutUint32(data[14:18], image.InfoHeader.Size)
	binary.LittleEndian.PutUint32(data[18:22], uint32(image.InfoHeader.Width))
	binary.LittleEndian.PutUint32(data[22:26], uint32(image.InfoHeader.Height))
	binary.LittleEndian.PutUint16(data[26:28], image.InfoHeader.Planes)
	binary.LittleEndian.PutUint16(data[28:30], image.InfoHeader.BitsPerPixel)
	binary.LittleEndian.PutUint32(data[30:34], image.InfoHeader.Compression)
	binary.LittleEndian.PutUint32(data[34:38], image.InfoHeader.ImageSize)
	binary.LittleEndian.PutUint32(data[38:42], uint32(image.InfoHeader.XPixelsPerMeter))
	binary.LittleEndian.PutUint32(data[42:46], uint32(image.InfoHeader.YPixelsPerMeter))
	binary.LittleEndian.PutUint32(data[46:50], image.InfoHeader.ColorsUsed)
	binary.LittleEndian.PutUint32(data[50:54], image.InfoHeader.ColorsImportant)

	// Serialize pixel data
	offset := headerSize
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := image.Data[y][x]
			data[offset] = pixel.Blue
			data[offset+1] = pixel.Green
			data[offset+2] = pixel.Red
			offset += bytesPerPixel
		}
	}

	return data
}

func SaveBMP(image *BMPImage, filename string) error {
	data := SerializeBMP(image)
	return os.WriteFile(filename, data, 0o644)
}

// PrintBMPHeaderInfo prints the BMP and DIB header information in a formatted style.
// It displays all relevant fields from both headers, providing a comprehensive
// overview of the BMP file structure and image properties.
//
// Parameters:
// - image: A pointer to the BMPImage struct containing the headers to print.
func PrintBMPHeaderInfo(image *BMPImage) {
	fmt.Printf(`BMP Header:
- Signature: %s
- FileSize: %d bytes
- DataOffset: %d bytes
DIB Header:
- Size: %d bytes
- Width: %d pixels
- Height: %d pixels
- Planes: %d
- BitsPerPixel: %d
- Compression: %d
- ImageSize: %d bytes
- XPixelsPerMeter: %d
- YPixelsPerMeter: %d
- ColorsUsed: %d
- ColorsImportant: %d
`,
		image.Header.Signature,
		image.Header.FileSize,
		image.Header.DataOffset,
		image.InfoHeader.Size,
		image.InfoHeader.Width,
		image.InfoHeader.Height,
		image.InfoHeader.Planes,
		image.InfoHeader.BitsPerPixel,
		image.InfoHeader.Compression,
		image.InfoHeader.ImageSize,
		image.InfoHeader.XPixelsPerMeter,
		image.InfoHeader.YPixelsPerMeter,
		image.InfoHeader.ColorsUsed,
		image.InfoHeader.ColorsImportant,
	)
}
