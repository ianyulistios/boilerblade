package main

import (
	"boilerblade/internal/cli"
	"fmt"
	"os"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		cli.ShowHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Error: Project name is required")
			fmt.Println("Usage: boilerblade new <project-name>")
			os.Exit(1)
		}
		projectName := os.Args[2]
		if err := cli.CreateNewProject(projectName); err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}

	case "make":
		if len(os.Args) < 3 {
			fmt.Println("Error: Resource type is required")
			fmt.Println("Usage: boilerblade make <resource> [options]")
			fmt.Println("\nAvailable resources:")
			fmt.Println("  model       - Generate model file")
			fmt.Println("  repository  - Generate repository file")
			fmt.Println("  usecase     - Generate usecase file")
			fmt.Println("  handler     - Generate handler file")
			fmt.Println("  dto         - Generate DTO file")
			fmt.Println("  consumer    - Generate RabbitMQ consumer (-name and optional -title; general-purpose)")
			fmt.Println("  migration   - Create Goose SQL migration (-name, e.g. add_orders_table)")
			fmt.Println("  all         - Generate all layers")
			os.Exit(1)
		}
		if err := cli.HandleMakeCommand(os.Args[2:]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "version", "-v", "--version":
		fmt.Printf("Boilerblade CLI v%s\n", version)
		os.Exit(0)

	case "help", "-h", "--help":
		cli.ShowHelp()
		os.Exit(0)

	default:
		// Backward compatibility: if it looks like old format, try to parse it
		if len(os.Args) > 1 && (os.Args[1] == "-layer" || os.Args[1] == "-name") {
			// Old format detected, redirect to make command
			if err := cli.HandleLegacyCommand(os.Args[1:]); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("Unknown command: %s\n\n", command)
			cli.ShowHelp()
			os.Exit(1)
		}
	}
}
