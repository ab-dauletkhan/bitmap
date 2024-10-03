package core

import (
	"fmt"
	"strings"
)

// TransformationType defines various types of transformations that can be applied to an image.
type TransformationType int

const (
	// MirrorTransform applies a horizontal or vertical flip to the image.
	MirrorTransform TransformationType = iota
	// FilterTransform applies a color or effect filter to the image (e.g., grayscale, blur).
	FilterTransform
	// RotateTransform rotates the image by a specified angle (90, 180 degrees).
	RotateTransform
	// CropTransform crops the image to a specified region.
	CropTransform
)

// Transform represents a single transformation operation, storing its type and any options.
type Transform struct {
	Type    TransformationType
	Options interface{} // Options vary depending on the transformation type
}

// MirrorOptions stores the direction for mirror transformations (e.g., "horizontal" or "vertical").
type MirrorOptions struct {
	Direction string
}

// FilterOptions stores the type of filter to be applied (e.g., "grayscale", "negative").
type FilterOptions struct {
	FilterType string
}

// RotateOptions stores the rotation angle (90 degrees left or right).
type RotateOptions struct {
	Angle int
}

// ParseTransformations parses command-line arguments to extract a list of image transformations,
// along with input and output file names. It handles multiple transformation flags, ensuring
// the transformations are applied in the specified order.
func ParseTransformations(args []string) ([]Transform, string, string, error) {
	var transforms []Transform

	if len(args) < 2 {
		return nil, "", "", ErrIncorrectArgument // Require at least input and output files.
	}

	inFile := args[len(args)-2]
	outFile := args[len(args)-1]

	for _, arg := range args[:len(args)-2] {
		switch {
		// Handle mirror transformations with various directional options.
		case strings.HasPrefix(arg, "--mirror="):
			opts := strings.Split(strings.TrimPrefix(arg, "--mirror="), ",")
			for _, opt := range opts {
				var direction string
				switch opt {
				case "horizontal", "h", "horizontally", "hor":
					direction = "horizontal"
				case "vertical", "v", "vertically", "ver":
					direction = "vertical"
				default:
					return nil, "", "", fmt.Errorf("invalid mirror option: %s", opt)
				}
				transforms = append(transforms, Transform{
					Type:    MirrorTransform,
					Options: MirrorOptions{Direction: direction},
				})
			}

		// Handle filter transformations for different color effects.
		case strings.HasPrefix(arg, "--filter="):
			filterType := strings.TrimPrefix(arg, "--filter=")
			switch filterType {
			case "blue", "red", "green", "grayscale", "negative", "pixelate", "blur":
				transforms = append(transforms, Transform{
					Type:    FilterTransform,
					Options: FilterOptions{FilterType: filterType},
				})
			default:
				return nil, "", "", fmt.Errorf("invalid filter option: %s", filterType)
			}

		// Handle rotate transformations with multiple angles (left, right, 180 degrees).
		case strings.HasPrefix(arg, "--rotate="):
			opts := strings.Split(strings.TrimPrefix(arg, "--rotate="), ",")
			for _, opt := range opts {
				var angle int
				switch opt {
				case "right", "90", "-270":
					angle = 1
				case "left", "-90", "270":
					angle = -1
				case "-180", "180":
					// 180-degree rotation is handled by applying two mirror operations.
					transforms = append(transforms,
						Transform{Type: MirrorTransform, Options: MirrorOptions{Direction: "horizontal"}},
						Transform{Type: MirrorTransform, Options: MirrorOptions{Direction: "vertical"}},
					)
					continue
				default:
					return nil, "", "", fmt.Errorf("invalid rotate option: %s", opt)
				}
				transforms = append(transforms, Transform{
					Type:    RotateTransform,
					Options: RotateOptions{Angle: angle},
				})
			}

		// Handle crop transformations by parsing crop-specific options.
		case strings.HasPrefix(arg, "--crop="):
			cropInfo, err := parseCropInfo(strings.TrimPrefix(arg, "--crop="))
			if err != nil {
				return nil, "", "", err
			}
			transforms = append(transforms, Transform{
				Type:    CropTransform,
				Options: cropInfo,
			})
		default:
			return nil, "", "", fmt.Errorf("incorrect argument: %s", arg)
		}
	}

	return transforms, inFile, outFile, nil
}

// ApplyTransformations applies the parsed transformations sequentially to the BMP image.
// Each transformation modifies the image based on the options provided.
func ApplyTransformations(image *BMPImage, transforms []Transform) error {
	for _, t := range transforms {
		switch t.Type {
		case MirrorTransform:
			opts := t.Options.(MirrorOptions)
			MirrorImage(image, opts.Direction)
		case FilterTransform:
			opts := t.Options.(FilterOptions)
			Filter(image, opts.FilterType)
		case RotateTransform:
			opts := t.Options.(RotateOptions)
			Rotate(image, opts.Angle)
		case CropTransform:
			opts := t.Options.(CropInfo)
			if err := Crop(image, opts); err != nil {
				return err
			}
		}
	}
	return nil
}
