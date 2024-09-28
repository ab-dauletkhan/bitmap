package bitmap

import (
	"fmt"
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
	}

	switch args[0] {
	case "header":
		fmt.Println("handle header command")
	case "apply":
		fmt.Println("handle apply command")
	case "--help", "-h":
		core.PrintUsage()
	default:
		fmt.Println("unknown command")
		core.PrintUsage()
		os.Exit(1)
	}

}
