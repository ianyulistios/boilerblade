# Generator Test Suite

Test suite untuk generator boilerplate code.

## Running Tests

### Run All Tests
```bash
go test ./internal/generator -v
```

### Run Specific Test
```bash
go test ./internal/generator -v -run TestNewGenerator
```

### Run Integration Tests
```bash
go test ./internal/generator -v -run TestIntegration
```

### Run with Coverage
```bash
go test ./internal/generator -v -cover
```

### Run with Coverage Report
```bash
go test ./internal/generator -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

Current coverage: **62.0%**

### Test Files

1. **`generator_test.go`** - Unit tests untuk generator functions
   - `TestNewGenerator` - Test generator initialization
   - `TestGetTableName` - Test table name generation
   - `TestGetRouteName` - Test route name generation
   - `TestGenerateModel` - Test model file generation
   - `TestGenerateRepository` - Test repository file generation
   - `TestGenerateUsecase` - Test usecase file generation
   - `TestGenerateHandler` - Test handler file generation
   - `TestGenerateDTO` - Test DTO file generation
   - `TestGenerateFile_FileExists` - Test file existence check
   - `TestPrepareModelData` - Test model data preparation
   - `TestPrepareRepositoryData` - Test repository data preparation
   - `TestPrepareUsecaseData` - Test usecase data preparation
   - `TestPrepareHandlerData` - Test handler data preparation
   - `TestPrepareDTOData` - Test DTO data preparation
   - `TestGenerateAll` - Test generating all layers

2. **`integration_test.go`** - Integration tests
   - `TestIntegration_GenerateAllLayers` - Test generating all layers together
   - `TestIntegration_FieldParsing` - Test field parsing with various inputs
   - `TestIntegration_EntityNameVariations` - Test different entity name formats
   - `TestIntegration_TemplateExecution` - Test template execution

## Test Structure

### Unit Tests
- Test individual functions and methods
- Use test directories to avoid conflicts
- Clean up test files after tests

### Integration Tests
- Test complete workflows
- Test multiple components together
- Verify end-to-end functionality

## Adding New Tests

### Example Test Structure

```go
func TestNewFeature(t *testing.T) {
    // Setup
    testDir := "test_output"
    os.MkdirAll(testDir, 0755)
    defer os.RemoveAll(testDir)

    // Test
    gen := NewGenerator("Product", "product", nil)
    // ... test logic ...

    // Verify
    if condition {
        t.Errorf("Expected X, got Y")
    }
}
```

### Best Practices

1. **Use test directories** - Create temporary directories for test files
2. **Clean up** - Always use `defer os.RemoveAll()` to clean up test files
3. **Table-driven tests** - Use table-driven tests for multiple scenarios
4. **Descriptive names** - Use descriptive test names
5. **Sub-tests** - Use `t.Run()` for sub-tests

## Test Data

Test files are created in temporary directories:
- `test_output/` - For unit tests
- `test_integration/` - For integration tests
- `test_template/` - For template tests

All test directories are automatically cleaned up after tests complete.
