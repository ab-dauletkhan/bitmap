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
	// Otherwise, checks the first argument and handles them separately
	args := os.Args[1:]
	if len(args) < 1 {
		core.PrintUsage()
		return
	}
	switch args[0] {
	// The switch statement handles different commands based on the first argument.
	// If the "header" command is provided, it requires a second argument,
	// which should be the file path of the bitmap image.
	// It attempts to read the file and parses its header.
	// If any error occurs (e.g., incorrect arguments or file read error),
	// the program exits with an appropriate error message.
	// If flags --help or -h are provided, then prints help message
	case "header":
		if len(args) < 2 {
			fmt.Println(core.ErrIncorrectArgument)
			os.Exit(1)
		}

		if strings.HasPrefix(args[1], "-") {
			if args[1] != "--help" && args[1] != "-h" {
				fmt.Println(core.ErrIncorrectArgument)
			}
			core.PrintUsage("header")
			return
		}
		bytes, err := os.ReadFile(args[1])
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
		fmt.Println("handle apply command")
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
