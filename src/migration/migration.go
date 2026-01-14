package migration

import (
	"boilerblade/helper"
	"boilerblade/src/model"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// ModelRegistry holds all models that need to be migrated
var ModelRegistry []interface{}

// RegisterModel registers a model for migration
func RegisterModel(model interface{}) {
	ModelRegistry = append(ModelRegistry, model)
}

// init registers all models for migration
func init() {
	// Register all models here
	RegisterModel(&model.User{})
	RegisterModel(&model.Product{})
	// Add more models here as needed:
	// RegisterModel(&model.Order{})
	// RegisterModel(&model.Category{})
}

// RunMigrations runs all registered model migrations
func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	if len(ModelRegistry) == 0 {
		helper.LogInfo("No models registered for migration", map[string]interface{}{
			"source": "RunMigrations",
		})
		return nil
	}

	helper.LogInfo("Starting database migrations", map[string]interface{}{
		"source":      "RunMigrations",
		"model_count": len(ModelRegistry),
	})

	// Run AutoMigrate for all registered models
	for _, model := range ModelRegistry {
		modelName := getModelName(model)
		helper.LogInfo("Migrating model", map[string]interface{}{
			"source":    "RunMigrations",
			"model":     modelName,
		})

		if err := db.AutoMigrate(model); err != nil {
			helper.LogError("Failed to migrate model", err, modelName, map[string]interface{}{
				"source": "RunMigrations",
				"model":  modelName,
			})
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}

		helper.LogInfo("Model migrated successfully", map[string]interface{}{
			"source": "RunMigrations",
			"model":  modelName,
		})
	}

	helper.LogInfo("All migrations completed successfully", map[string]interface{}{
		"source":      "RunMigrations",
		"model_count": len(ModelRegistry),
	})

	return nil
}

// getModelName extracts the model name from the model interface
func getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// GetRegisteredModels returns list of registered model names
func GetRegisteredModels() []string {
	models := make([]string, len(ModelRegistry))
	for i, model := range ModelRegistry {
		models[i] = getModelName(model)
	}
	return models
}
