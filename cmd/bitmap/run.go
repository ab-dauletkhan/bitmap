package bitmap

import (
	"fmt"
	"os"
	"strings"

	"github.com/ab-dauletkhan/bitmap/internal/core"
)

func Run() {
	// Manually parses the program arguments
	//
	// If no argument is given, print usage (help)
	// Otherwise, checks the first argument and handle them separately
	args := os.Args[1:]
	if len(args) < 1 {
		core.PrintUsage()
		return
	}
	command := args[0]
	args = args[1:]
	switch command {
	// If the "header" command is provided, it requires a second argument,
	// which should be the file path of the bitmap image or help flag.
	// It attempts to identify it as a flag first, if not reads the file
	// and parses its header.
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
		transforms, inFile, outFile, err := core.ParseTransformations(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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

		// Apply all transformations in sequence
		if err := core.ApplyTransformations(image, transforms); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := core.SaveBMP(image, outFile); err != nil {
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
