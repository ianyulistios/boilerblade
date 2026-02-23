package cli

import (
	"boilerblade/internal/generator"
	"flag"
	"fmt"
	"os"
	"strings"
)

// HandleMakeCommand processes the make command
func HandleMakeCommand(args []string) error {
	// Ensure project has .env.example when generating (so env is part of "generate project")
	if wd, err := os.Getwd(); err == nil {
		_ = EnsureEnvExample(wd)
	}

	if len(args) < 1 {
		return fmt.Errorf("resource type is required")
	}

	resourceType := args[0]
	remainingArgs := args[1:]

	// Parse flags
	fs := flag.NewFlagSet("make", flag.ContinueOnError)
	name := fs.String("name", "", "Name of the entity or consumer (e.g., Product, OrderEvents)")
	title := fs.String("title", "", "Title for consumer only (e.g., \"Order Events\"); optional")
	fields := fs.String("fields", "", "Fields for model (format: Name:string:required,Price:float64:required)")

	if err := fs.Parse(remainingArgs); err != nil {
		return err
	}

	resourceLower := strings.ToLower(resourceType)

	// Consumer: general-purpose, only -name and -title
	if resourceLower == "consumer" {
		if *name == "" {
			return fmt.Errorf("consumer name is required (use -name flag, e.g. -name=OrderEvents or -name=order_events)")
		}
		consumerGen := generator.NewConsumerGen(*name, *title)
		if err := consumerGen.GenerateAMQPConstants(); err != nil {
			return fmt.Errorf("generating AMQP constants: %w", err)
		}
		fmt.Printf("✓ AMQP constants generated (constants/amqp_%s.go)\n", consumerGen.Identifier)
		if err := consumerGen.GenerateConsumer(); err != nil {
			return fmt.Errorf("generating consumer: %w", err)
		}
		fmt.Printf("✓ RabbitMQ consumer \"%s\" generated (src/consumer/%s.go)\n", consumerGen.Title, consumerGen.Identifier)
		fmt.Println("  Register the consumer in server/amqp.go and add your logic in handleCreatedMessage/handleUpdatedMessage.")
		return nil
	}

	// Migration: Goose SQL migration (postgres + mysql placeholder files)
	if resourceLower == "migration" {
		if *name == "" {
			return fmt.Errorf("migration name is required (use -name flag, e.g. -name=add_orders_table)")
		}
		migrationGen := generator.NewMigrationGen(*name)
		if err := migrationGen.Generate(); err != nil {
			return fmt.Errorf("generating migration: %w", err)
		}
		pg, my := migrationGen.GeneratedFiles()
		fmt.Printf("✓ Goose migration created: %s\n", pg)
		fmt.Printf("✓ Goose migration created: %s\n", my)
		fmt.Println("  Edit the files to add your Up/Down SQL, then run the app or use goose up.")
		return nil
	}

	if *name == "" {
		return fmt.Errorf("entity name is required (use -name flag)")
	}

	// Convert name to proper format for entity layers
	entityName := strings.Title(strings.ToLower(*name))
	entityNameLower := strings.ToLower(*name)

	// Parse fields if provided
	var modelFields []generator.Field
	if *fields != "" {
		modelFields = parseFields(*fields)
	}

	gen := generator.NewGenerator(entityName, entityNameLower, modelFields)

	switch resourceLower {
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
		return fmt.Errorf("unknown resource type: %s. Available: model, repository, usecase, handler, dto, consumer, migration, all", resourceType)
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
