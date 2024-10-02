package core

import (
	"fmt"
	"strconv"
	"strings"
)

// TransformationType represents the type of transformation
type TransformationType int

const (
	MirrorTransform TransformationType = iota
	FilterTransform
	RotateTransform
	CropTransform
)

// Transform represents a single transformation operation
type Transform struct {
	Type    TransformationType
	Options interface{} // stores the specific options for each transformation
}

// MirrorOptions represents options for mirror transformation
type MirrorOptions struct {
	Direction string
}

// FilterOptions represents options for filter transformation
type FilterOptions struct {
	FilterType string
}

// RotateOptions represents options for rotate transformation
type RotateOptions struct {
	Angle int
}

func ParseTransformations(args []string) ([]Transform, string, string, error) {
	var transforms []Transform

	if len(args) < 2 {
		return nil, "", "", ErrIncorrectArgument
	}

	inFile := args[len(args)-2]
	outFile := args[len(args)-1]

	// Process each argument
	for _, arg := range args[:len(args)-2] {
		switch {
		case strings.HasPrefix(arg, "--mirror="):
			opts := strings.Split(strings.TrimPrefix(arg, "--mirror="), ",")
			for _, opt := range opts {
				direction := ""
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

		case strings.HasPrefix(arg, "--crop="):
			cropInfo, err := parseCropInfo(strings.TrimPrefix(arg, "--crop="))
			if err != nil {
				return nil, "", "", err
			}
			transforms = append(transforms, Transform{
				Type:    CropTransform,
				Options: cropInfo,
			})
		}
	}

	return transforms, inFile, outFile, nil
}

func parseCropInfo(cropStr string) (CropInfo, error) {
	info := strings.Split(cropStr, "-")
	var cropInfo CropInfo

	if len(info) != 2 && len(info) != 4 {
		return cropInfo, fmt.Errorf("crop option must have either 2 or 4 values")
	}

	var err error
	cropInfo.OffsetX, err = strconv.Atoi(info[0])
	if err != nil {
		return cropInfo, fmt.Errorf("invalid OffsetX value")
	}

	cropInfo.OffsetY, err = strconv.Atoi(info[1])
	if err != nil {
		return cropInfo, fmt.Errorf("invalid OffsetY value")
	}

	if len(info) == 4 {
		cropInfo.Width, err = strconv.Atoi(info[2])
		if err != nil {
			return cropInfo, fmt.Errorf("invalid Width value")
		}

		cropInfo.Height, err = strconv.Atoi(info[3])
		if err != nil {
			return cropInfo, fmt.Errorf("invalid Height value")
		}
	}

	return cropInfo, nil
}

// ApplyTransformations applies all transformations in sequence
func ApplyTransformations(image *BMPImage, transforms []Transform) error {
	for _, t := range transforms {
		switch t.Type {
		case MirrorTransform:
			opts := t.Options.(MirrorOptions)
			image = MirrorImage(image, opts.Direction)
		case FilterTransform:
			opts := t.Options.(FilterOptions)
			Filter(image, opts.FilterType)
		case RotateTransform:
			opts := t.Options.(RotateOptions)
			Rotate(image, opts.Angle)
		case CropTransform:
			opts := t.Options.(CropInfo)
			err := Crop(image, opts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
