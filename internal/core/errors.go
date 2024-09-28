package core

import "errors"

var (
	// Argument-related errors
	ErrIncorrectArgument = errors.New("incorrect arguments provided; please check your input")
	ErrMissingFilename   = errors.New("missing BMP filename; please specify a valid file for the header command")

	// BMP file structure-related errors
	ErrInvalidBMP            = errors.New("invalid BMP file format; ensure the file is uncompressed 24-bit BMP")
	ErrCorruptFile           = errors.New("corrupt BMP file; file size does not match the expected size in headers")
	ErrUnsupportedFormat     = errors.New("unsupported BMP format; only 24-bit uncompressed BMP files are supported")
	ErrInvalidImageData      = errors.New("invalid image data; data size does not match the dimensions or file structure")
	ErrInvalidHeaderSize     = errors.New("invalid BMP header size; header does not conform to BMP specification")
	ErrInvalidFileType       = errors.New("file is not recognized as a valid BMP file; ensure it starts with 'BM'")
	ErrNonPositiveDimensions = errors.New("invalid image dimensions; width and height must be positive values")
)
