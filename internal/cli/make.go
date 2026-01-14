package cli

import (
	"boilerblade/internal/generator"
	"flag"
	"fmt"
	"strings"
)

// HandleMakeCommand processes the make command
func HandleMakeCommand(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("resource type is required")
	}

	resourceType := args[0]
	remainingArgs := args[1:]

	// Parse flags
	fs := flag.NewFlagSet("make", flag.ContinueOnError)
	name := fs.String("name", "", "Name of the entity (e.g., Product, Order)")
	fields := fs.String("fields", "", "Fields for model (format: Name:string:required,Price:float64:required)")

	if err := fs.Parse(remainingArgs); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("entity name is required (use -name flag)")
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

	switch strings.ToLower(resourceType) {
	case "model":
		if err := gen.GenerateModel(); err != nil {
			return fmt.Errorf("generating model: %w", err)
		}
		fmt.Printf("✓ Model %s generated successfully\n", entityName)

	case "repository":
		if err := gen.GenerateRepository(); err != nil {
			return fmt.Errorf("generating repository: %w", err)
		}
		fmt.Printf("✓ Repository %s generated successfully\n", entityName)

	case "usecase":
		if err := gen.GenerateUsecase(); err != nil {
			return fmt.Errorf("generating usecase: %w", err)
		}
		fmt.Printf("✓ Usecase %s generated successfully\n", entityName)

	case "handler":
		if err := gen.GenerateHandler(); err != nil {
			return fmt.Errorf("generating handler: %w", err)
		}
		fmt.Printf("✓ Handler %s generated successfully\n", entityName)

	case "dto":
		if err := gen.GenerateDTO(); err != nil {
			return fmt.Errorf("generating DTO: %w", err)
		}
		fmt.Printf("✓ DTO %s generated successfully\n", entityName)

	case "all":
		if err := gen.GenerateAll(); err != nil {
			return fmt.Errorf("generating all layers: %w", err)
		}
		fmt.Printf("✓ All layers for %s generated successfully\n", entityName)

	default:
		return fmt.Errorf("unknown resource type: %s. Available: model, repository, usecase, handler, dto, all", resourceType)
	}

	return nil
}

// HandleLegacyCommand handles the old flag-based command format for backward compatibility
func HandleLegacyCommand(args []string) error {
	fs := flag.NewFlagSet("legacy", flag.ContinueOnError)
	layer := fs.String("layer", "", "Layer to generate")
	name := fs.String("name", "", "Name of the entity")
	fields := fs.String("fields", "", "Fields for model")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *layer == "" || *name == "" {
		return fmt.Errorf("both -layer and -name are required")
	}

	// Convert to new format
	newArgs := []string{*layer}
	if *name != "" {
		newArgs = append(newArgs, "-name", *name)
	}
	if *fields != "" {
		newArgs = append(newArgs, "-fields", *fields)
	}

	return HandleMakeCommand(newArgs)
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
