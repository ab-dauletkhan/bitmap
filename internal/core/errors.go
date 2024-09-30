package core

import "errors"

var (
	// Argument-related errors
	ErrIncorrectArgument = errors.New("incorrect arguments provided; please check your input")
	ErrMissingFilename   = errors.New("missing BMP filename; please specify a valid file for the header command")

	// Error variables for various BMP parsing and validation errors.
	ErrInvalidBMP             = errors.New("invalid BMP file")
	ErrInvalidFileType        = errors.New("invalid file type, expected 'BM'")
	ErrCorruptFile            = errors.New("corrupt BMP file")
	ErrInvalidHeaderSize      = errors.New("invalid header size")
	ErrNonPositiveDimensions  = errors.New("non-positive image dimensions")
	ErrUnsupportedFormat      = errors.New("unsupported BMP format")
	ErrInvalidImageData       = errors.New("invalid image data")
	ErrUnsupportedCompression = errors.New("unsupported compression method")
)
