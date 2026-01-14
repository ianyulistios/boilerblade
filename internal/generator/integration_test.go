package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_GenerateAllLayers(t *testing.T) {
	// Setup test directory
	testDir := "test_integration"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// Create subdirectories
	modelDir := filepath.Join(testDir, "model")
	repoDir := filepath.Join(testDir, "repository")
	usecaseDir := filepath.Join(testDir, "usecase")
	handlerDir := filepath.Join(testDir, "handler")
	dtoDir := filepath.Join(testDir, "dto")

	os.MkdirAll(modelDir, 0755)
	os.MkdirAll(repoDir, 0755)
	os.MkdirAll(usecaseDir, 0755)
	os.MkdirAll(handlerDir, 0755)
	os.MkdirAll(dtoDir, 0755)

	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: "required"},
		{Name: "Stock", Type: "int", Tag: "required"},
	}

	gen := NewGenerator("Product", "product", fields)

	// Test generating all layers
	// Note: We'll test with simplified paths to avoid conflicts
	t.Run("GenerateModel", func(t *testing.T) {
		testPath := filepath.Join(modelDir, "product.go")
		err := gen.generateFile(testPath, `package model
type {{.EntityName}} struct {
	ID uint
}`, gen.prepareModelData())
		if err != nil {
			t.Fatalf("Failed to generate model: %v", err)
		}
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("Model file not created")
		}
	})

	t.Run("GenerateRepository", func(t *testing.T) {
		testPath := filepath.Join(repoDir, "product.go")
		err := gen.generateFile(testPath, `package repository
type {{.EntityName}}Repository interface {}`, gen.prepareRepositoryData())
		if err != nil {
			t.Fatalf("Failed to generate repository: %v", err)
		}
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("Repository file not created")
		}
	})

	t.Run("GenerateUsecase", func(t *testing.T) {
		testPath := filepath.Join(usecaseDir, "product.go")
		err := gen.generateFile(testPath, `package usecase
type {{.EntityName}}Usecase interface {}`, gen.prepareUsecaseData())
		if err != nil {
			t.Fatalf("Failed to generate usecase: %v", err)
		}
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("Usecase file not created")
		}
	})

	t.Run("GenerateHandler", func(t *testing.T) {
		testPath := filepath.Join(handlerDir, "product.go")
		err := gen.generateFile(testPath, `package handler
type {{.EntityName}}Handler struct {}`, gen.prepareHandlerData())
		if err != nil {
			t.Fatalf("Failed to generate handler: %v", err)
		}
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("Handler file not created")
		}
	})

	t.Run("GenerateDTO", func(t *testing.T) {
		testPath := filepath.Join(dtoDir, "product.go")
		err := gen.generateFile(testPath, `package dto
type {{.EntityName}}Response struct {}`, gen.prepareDTOData())
		if err != nil {
			t.Fatalf("Failed to generate DTO: %v", err)
		}
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Errorf("DTO file not created")
		}
	})
}

func TestIntegration_FieldParsing(t *testing.T) {
	tests := []struct {
		name     string
		fields   []Field
		expected int
	}{
		{
			name: "Empty fields",
			fields: []Field{},
			expected: 0,
		},
		{
			name: "Single field",
			fields: []Field{
				{Name: "Name", Type: "string", Tag: "required"},
			},
			expected: 1,
		},
		{
			name: "Multiple fields",
			fields: []Field{
				{Name: "Name", Type: "string", Tag: "required"},
				{Name: "Price", Type: "float64", Tag: "required"},
				{Name: "Stock", Type: "int", Tag: ""},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator("Product", "product", tt.fields)
			if len(gen.Fields) != tt.expected {
				t.Errorf("Expected %d fields, got %d", tt.expected, len(gen.Fields))
			}
		})
	}
}

func TestIntegration_EntityNameVariations(t *testing.T) {
	tests := []struct {
		entityName      string
		entityNameLower string
		expectedTable   string
		expectedRoute   string
	}{
		{"Product", "product", "products", "products"},
		{"Order", "order", "orders", "orders"},
		{"OrderItem", "orderitem", "orderitems", "orderitems"},
		{"User", "user", "users", "users"},
	}

	for _, tt := range tests {
		t.Run(tt.entityName, func(t *testing.T) {
			gen := NewGenerator(tt.entityName, tt.entityNameLower, nil)
			
			tableName := gen.getTableName()
			if tableName != tt.expectedTable {
				t.Errorf("Expected table name '%s', got '%s'", tt.expectedTable, tableName)
			}

			routeName := gen.getRouteName()
			if routeName != tt.expectedRoute {
				t.Errorf("Expected route name '%s', got '%s'", tt.expectedRoute, routeName)
			}
		})
	}
}

func TestIntegration_TemplateExecution(t *testing.T) {
	testDir := "test_template"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gen := NewGenerator("Product", "product", []Field{
		{Name: "Name", Type: "string", Tag: "required"},
	})

	testPath := filepath.Join(testDir, "test.go")
	tmpl := `Entity: {{.EntityName}}
Lower: {{.EntityNameLower}}
Table: {{.TableName}}
Route: {{.RouteName}}
Fields: {{len .Fields}}
`

	// Use prepareModelData which includes all necessary fields
	data := gen.prepareModelData()
	
	err := gen.generateFile(testPath, tmpl, data)
	if err != nil {
		t.Fatalf("Failed to generate file: %v", err)
	}

	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Entity: Product") {
		t.Errorf("Template should contain 'Entity: Product', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "Lower: product") {
		t.Errorf("Template should contain 'Lower: product', got: %s", contentStr)
	}
}
