package bitmap

import (
	"os"

	"github.com/ab-dauletkhan/bitmap/internal/core"
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
	// If the "header" command is provided, it requires a second argument,
	// which should be the file path of the bitmap image or help flag.
	// It attempts to identify it as a flag first, if not reads the file
	// and parses its header.
	// If any error occurs (e.g., incorrect arguments or file read error),
	// the program exits with an appropriate error message.
	// If flags --help or -h are provided, then prints help message
	case "header":
		if len(args) != 1 {
			core.PrintErrorUsageExit(core.ErrIncorrectArgument, "header")
		}

		// if !strings.HasSuffix(args[0], ".bmp") {
		// 	if args[0] == "--help" || args[0] == "-h" {
		// 		core.PrintUsage("header")
		// 		return
		// 	} else {
		// 		core.PrintErrorUsageExit(core.ErrIncorrectArgument, "header")
		// 	}
		// }

		if args[0] == "--help" || args[0] == "-h" {
			core.PrintUsage("header")
			return
		}

		bytes, err := os.ReadFile(args[0])
		if err != nil {
			core.PrintErrorExit(err)
		}

		image, err := core.ParseBMP(bytes)
		if err != nil {
			core.PrintErrorExit(err)
		}
		core.PrintBMPHeaderInfo(image)

	// If the "apply" command is provided, it processes various transformation options
	// (mirror, filter, rotate, crop) and applies them to the input image in sequence.
	// The command requires an input file and output file as the last two arguments.
	// All files must be .bmp format.
	// If any error occurs during processing (invalid options, file operations, etc.),
	// the program exits with an appropriate error message.
	case "apply":
		if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
			core.PrintUsage("apply")
			return
		}
		transforms, inFile, outFile, err := core.ParseTransformations(args)
		if err != nil {
			core.PrintErrorUsageExit(err, "apply")
		}

		// if !strings.HasSuffix(inFile, ".bmp") || !strings.HasSuffix(outFile, ".bmp") {
		// 	core.PrintError(core.ErrInvalidFileType)
		// }

		bytes, err := os.ReadFile(inFile)
		if err != nil {
			core.PrintErrorExit(err)
		}

		image, err := core.ParseBMP(bytes)
		if err != nil {
			core.PrintErrorExit(err)
		}

		if err := core.ApplyTransformations(image, transforms); err != nil {
			core.PrintErrorExit(err)
		}

		if err := core.SaveBMP(image, outFile); err != nil {
			core.PrintErrorExit(err)
		}

	case "--help", "-h":
		core.PrintUsage()

	default:
		core.PrintErrorUsageExit(core.ErrUnknownCmd, "main")
	}
}
