package bitmap

import (
	"fmt"
	"os"
	"strings"

	"github.com/ab-dauletkhan/bitmap/internal/core"
	"github.com/ab-dauletkhan/bitmap/package/transformations"
)

func Run() {
	// Manually parses the program arguments
	//
	// If no argument is given, print usage (help)
	// Otherwise, checks the first argument and handles them separately
	args := os.Args[1:]
	if len(args) < 1 {
		core.PrintUsage()
		return
	}
	command := args[0]
	args = args[1:]
	switch command {
	// The switch statement handles different commands based on the first argument.
	// If the "header" command is provided, it requires a second argument,
	// which should be the file path of the bitmap image.
	// It attempts to read the file and parses its header.
	// If any error occurs (e.g., incorrect arguments or file read error),
	// the program exits with an appropriate error message.
	// If flags --help or -h are provided, then prints help message
	case "header":
		if len(args) != 1 {
			fmt.Println(core.ErrIncorrectArgument)
			os.Exit(1)
		}

		if strings.HasPrefix(args[0], "-") {
			if args[0] != "--help" && args[0] != "-h" {
				fmt.Println(core.ErrIncorrectArgument)
			}
			core.PrintUsage("header")
			return
		}
		bytes, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		image, err := core.ParseBMP(bytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		core.PrintBMPHeaderInfo(image)
	case "apply":
		if len(args) < 2 {
			fmt.Println(core.ErrIncorrectArgument)
			core.PrintUsage("apply")
			os.Exit(1)
		}

		mirrorOpts := []string{}
		filterOpts := []string{}
		for _, arg := range args {
			if strings.HasPrefix(arg, "--mirror=") {
				opts := strings.Split(strings.TrimPrefix(arg, "--mirror="), ",")
				mirrorOpts = append(mirrorOpts, opts...)
			} else if strings.HasPrefix(arg, "--filter") {
				filterOpts = append(filterOpts, strings.TrimPrefix(arg, "--filter="))
			}
		}

		inFile := args[len(args)-2]
		outFile := args[len(args)-1]

		if !strings.HasSuffix(inFile, ".bmp") || !strings.HasSuffix(outFile, ".bmp") {
			fmt.Println(core.ErrInvalidFileType)
			os.Exit(1)
		}

		bytes, err := os.ReadFile(inFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		image, err := core.ParseBMP(bytes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, opt := range mirrorOpts {
			switch opt {
			case "horizontal", "h", "horizontally", "hor":
				image = transformations.MirrorImage(image, "horizontal")
			case "vertical", "v", "vertically", "ver":
				image = transformations.MirrorImage(image, "vertical")
			default:
				fmt.Printf("Invalid mirror option: %s\n", opt)
				os.Exit(1)
			}
		}

		for _, opt := range filterOpts {
			switch opt {
			case "blue", "red", "green", "grayscale", "negative", "pixelate", "blur":
				transformations.Filter(image, opt)
			default:
				fmt.Printf("Invalid filter option: %s\n", opt)
				os.Exit(1)
			}
		}

		err = core.SaveBMP(image, outFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	case "--help", "-h":
		core.PrintUsage()
		return
	default:
		fmt.Println("unknown command")
		core.PrintUsage()
		os.Exit(1)
	}
}
