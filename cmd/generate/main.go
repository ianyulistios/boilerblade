package main

import (
	"boilerblade/internal/generator"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	version = "1.0.0"
)

func main() {
	var (
		layer      = flag.String("layer", "", "Layer to generate: model, repository, usecase, handler, dto, all")
		name       = flag.String("name", "", "Name of the entity (e.g., Product, Order)")
		fields     = flag.String("fields", "", "Fields for model (format: name:type:tag, e.g., Name:string:required,Price:float64:required)")
		help       = flag.Bool("help", false, "Show help message")
		showVersion = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Boilerblade CLI Generator v%s\n", version)
		os.Exit(0)
	}

	if *help || *layer == "" || *name == "" {
		printUsage()
		os.Exit(0)
	}

	// Convert name to proper format
	entityName := strings.Title(strings.ToLower(*name))
	entityNameLower := strings.ToLower(*name)

	// Parse fields if provided
	var modelFields []generator.Field
	if *fields != "" {
		modelFields = parseFields(*fields)
	}

	gen := generator.NewGenerator(entityName, entityNameLower, modelFields)

	switch strings.ToLower(*layer) {
	case "model":
		if err := gen.GenerateModel(); err != nil {
			fmt.Printf("Error generating model: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Model %s generated successfully\n", entityName)

	case "repository":
		if err := gen.GenerateRepository(); err != nil {
			fmt.Printf("Error generating repository: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Repository %s generated successfully\n", entityName)

	case "usecase":
		if err := gen.GenerateUsecase(); err != nil {
			fmt.Printf("Error generating usecase: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Usecase %s generated successfully\n", entityName)

	case "handler":
		if err := gen.GenerateHandler(); err != nil {
			fmt.Printf("Error generating handler: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Handler %s generated successfully\n", entityName)

	case "dto":
		if err := gen.GenerateDTO(); err != nil {
			fmt.Printf("Error generating DTO: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ DTO %s generated successfully\n", entityName)

	case "all":
		if err := gen.GenerateAll(); err != nil {
			fmt.Printf("Error generating all: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ All layers for %s generated successfully\n", entityName)

	default:
		fmt.Printf("Unknown layer: %s\n", *layer)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Boilerblade Code Generator")
	fmt.Println("Usage: boilerblade -layer=<layer> -name=<name> [options]")
	fmt.Println()
	fmt.Println("Layers:")
	fmt.Println("  model       - Generate model file")
	fmt.Println("  repository  - Generate repository file")
	fmt.Println("  usecase     - Generate usecase file")
	fmt.Println("  handler     - Generate handler file")
	fmt.Println("  dto         - Generate DTO file")
	fmt.Println("  all         - Generate all layers")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -name       Entity name (required, e.g., Product, Order)")
	fmt.Println("  -fields     Model fields (format: Name:string:required,Price:float64:required)")
	fmt.Println("  -help       Show this help message")
	fmt.Println("  -version    Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  boilerblade -layer=model -name=Product -fields=\"Name:string:required,Price:float64:required\"")
	fmt.Println("  boilerblade -layer=all -name=Order")
	fmt.Println("  boilerblade -layer=repository -name=Product")
}

func parseFields(fieldsStr string) []generator.Field {
	var fields []generator.Field
	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fieldParts := strings.Split(part, ":")
		if len(fieldParts) < 2 {
			continue
		}
		fieldName := strings.TrimSpace(fieldParts[0])
		fieldType := strings.TrimSpace(fieldParts[1])
		fieldTag := ""
		if len(fieldParts) > 2 {
			fieldTag = strings.TrimSpace(fieldParts[2])
		}
		fields = append(fields, generator.Field{
			Name: fieldName,
			Type: fieldType,
			Tag:  fieldTag,
		})
	}
	return fields
}
