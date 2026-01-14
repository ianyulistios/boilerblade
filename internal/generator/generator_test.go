package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: "required"},
	}

	gen := NewGenerator("Product", "product", fields)

	if gen.EntityName != "Product" {
		t.Errorf("Expected EntityName to be 'Product', got '%s'", gen.EntityName)
	}

	if gen.EntityNameLower != "product" {
		t.Errorf("Expected EntityNameLower to be 'product', got '%s'", gen.EntityNameLower)
	}

	if len(gen.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(gen.Fields))
	}
}

func TestGetTableName(t *testing.T) {
	gen := NewGenerator("Product", "product", nil)
	tableName := gen.getTableName()

	expected := "products"
	if tableName != expected {
		t.Errorf("Expected table name '%s', got '%s'", expected, tableName)
	}
}

func TestGetRouteName(t *testing.T) {
	gen := NewGenerator("Product", "product", nil)
	routeName := gen.getRouteName()

	expected := "products"
	if routeName != expected {
		t.Errorf("Expected route name '%s', got '%s'", expected, routeName)
	}
}

func TestGenerateModel(t *testing.T) {
	// Setup
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: "required"},
	}

	gen := NewGenerator("Product", "product", fields)

	// Use test directory for output
	testPath := filepath.Join(testDir, "product.go")

	// Test model generation
	err := gen.generateFile(testPath, `package model

type {{.EntityName}} struct {
	ID uint
{{range .Fields}}	{{.Name}} {{.Type}}
{{end}}}
`, gen.prepareModelData())

	if err != nil {
		t.Fatalf("Failed to generate model: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", testPath)
	}

	// Verify file content
	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "Product") {
		t.Errorf("Generated file should contain 'Product', got: %s", contentStr)
	}

	if !contains(contentStr, "Name string") {
		t.Errorf("Generated file should contain 'Name string', got: %s", contentStr)
	}
}

func TestGenerateRepository(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gen := NewGenerator("Product", "product", nil)
	testPath := filepath.Join(testDir, "product_repo.go")

	err := gen.generateFile(testPath, `package repository

type {{.EntityName}}Repository interface {
	Create({{.EntityNameLower}} *model.{{.EntityName}}) error
}
`, gen.prepareRepositoryData())

	if err != nil {
		t.Fatalf("Failed to generate repository: %v", err)
	}

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", testPath)
	}
}

func TestGenerateUsecase(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gen := NewGenerator("Product", "product", nil)
	testPath := filepath.Join(testDir, "product_usecase.go")

	err := gen.generateFile(testPath, `package usecase

type {{.EntityName}}Usecase interface {
	Create{{.EntityName}}(req *dto.Create{{.EntityName}}Request) error
}
`, gen.prepareUsecaseData())

	if err != nil {
		t.Fatalf("Failed to generate usecase: %v", err)
	}

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", testPath)
	}
}

func TestGenerateHandler(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gen := NewGenerator("Product", "product", nil)
	testPath := filepath.Join(testDir, "product_handler.go")

	err := gen.generateFile(testPath, `package handler

type {{.EntityName}}Handler struct {
	{{.EntityNameLower}}Usecase usecase.{{.EntityName}}Usecase
}
`, gen.prepareHandlerData())

	if err != nil {
		t.Fatalf("Failed to generate handler: %v", err)
	}

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", testPath)
	}
}

func TestGenerateDTO(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: "required"},
	}

	gen := NewGenerator("Product", "product", fields)
	testPath := filepath.Join(testDir, "product_dto.go")

	err := gen.generateFile(testPath, `package dto

type Create{{.EntityName}}Request struct {
{{range .Fields}}	{{.Name}} {{.Type}}
{{end}}}
`, gen.prepareDTOData())

	if err != nil {
		t.Fatalf("Failed to generate DTO: %v", err)
	}

	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Generated file does not exist: %s", testPath)
	}
}

func TestGenerateFile_FileExists(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	gen := NewGenerator("Product", "product", nil)
	testPath := filepath.Join(testDir, "existing.go")

	// Create file first
	os.WriteFile(testPath, []byte("existing content"), 0644)

	// Try to generate - should fail
	err := gen.generateFile(testPath, "package test", nil)

	if err == nil {
		t.Errorf("Expected error when file exists, got nil")
	}

	if !contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestPrepareModelData(t *testing.T) {
	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: ""},
	}

	gen := NewGenerator("Product", "product", fields)
	data := gen.prepareModelData()

	if data["EntityName"] != "Product" {
		t.Errorf("Expected EntityName 'Product', got '%v'", data["EntityName"])
	}

	if data["EntityNameLower"] != "product" {
		t.Errorf("Expected EntityNameLower 'product', got '%v'", data["EntityNameLower"])
	}

	fieldsData := data["Fields"].([]map[string]interface{})
	if len(fieldsData) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fieldsData))
	}

	if fieldsData[0]["Name"] != "Name" {
		t.Errorf("Expected first field name 'Name', got '%v'", fieldsData[0]["Name"])
	}
}

func TestPrepareRepositoryData(t *testing.T) {
	gen := NewGenerator("Product", "product", nil)
	data := gen.prepareRepositoryData()

	if data["EntityName"] != "Product" {
		t.Errorf("Expected EntityName 'Product', got '%v'", data["EntityName"])
	}

	if data["EntityNameLower"] != "product" {
		t.Errorf("Expected EntityNameLower 'product', got '%v'", data["EntityNameLower"])
	}
}

func TestPrepareUsecaseData(t *testing.T) {
	gen := NewGenerator("Product", "product", nil)
	data := gen.prepareUsecaseData()

	if data["EntityName"] != "Product" {
		t.Errorf("Expected EntityName 'Product', got '%v'", data["EntityName"])
	}
}

func TestPrepareHandlerData(t *testing.T) {
	gen := NewGenerator("Product", "product", nil)
	data := gen.prepareHandlerData()

	if data["EntityName"] != "Product" {
		t.Errorf("Expected EntityName 'Product', got '%v'", data["EntityName"])
	}

	if data["RouteName"] != "products" {
		t.Errorf("Expected RouteName 'products', got '%v'", data["RouteName"])
	}
}

func TestPrepareDTOData(t *testing.T) {
	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
		{Name: "Price", Type: "float64", Tag: "email"},
	}

	gen := NewGenerator("Product", "product", fields)
	data := gen.prepareDTOData()

	fieldsData := data["Fields"].([]map[string]interface{})
	if len(fieldsData) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(fieldsData))
	}

	if fieldsData[0]["ValidateTag"] != "required" {
		t.Errorf("Expected ValidateTag 'required', got '%v'", fieldsData[0]["ValidateTag"])
	}
}

func TestGenerateAll(t *testing.T) {
	testDir := "test_output"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	fields := []Field{
		{Name: "Name", Type: "string", Tag: "required"},
	}

	gen := NewGenerator("TestEntity", "testentity", fields)

	// Test that GenerateAll calls all generation methods
	// We'll test each method individually, so this is more of an integration test
	// For now, just verify the generator is set up correctly
	if gen.EntityName != "TestEntity" {
		t.Errorf("Expected EntityName 'TestEntity', got '%s'", gen.EntityName)
	}
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
