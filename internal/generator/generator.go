package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Generator struct {
	EntityName      string
	EntityNameLower string
	Fields          []Field
}

func NewGenerator(entityName, entityNameLower string, fields []Field) *Generator {
	return &Generator{
		EntityName:      entityName,
		EntityNameLower: entityNameLower,
		Fields:          fields,
	}
}

func (g *Generator) GenerateModel() error {
	tmpl := `package model

import (
	"time"

	"gorm.io/gorm"
)

// {{.EntityName}} represents the {{.EntityNameLower}} entity in the database
type {{.EntityName}} struct {
	ID        uint           ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.NameLower}}\" gorm:\"{{.GormTag}}\"`" + `
{{end}}	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`json:\"-\" gorm:\"index\"`" + `
}

// TableName specifies the table name for {{.EntityName}} model
func ({{.EntityName}}) TableName() string {
	return "{{.TableName}}"
}
`

	return g.generateFile("src/model/"+g.EntityNameLower+".go", tmpl, g.prepareModelData())
}

func (g *Generator) GenerateRepository() error {
	tmpl := `package repository

import (
	"boilerblade/src/model"

	"gorm.io/gorm"
)

// {{.EntityName}}Repository defines the interface for {{.EntityNameLower}} data operations
type {{.EntityName}}Repository interface {
	Create({{.EntityNameLower}} *model.{{.EntityName}}) error
	GetByID(id uint) (*model.{{.EntityName}}, error)
	GetAll(limit, offset int) ([]model.{{.EntityName}}, error)
	Update({{.EntityNameLower}} *model.{{.EntityName}}) error
	Delete(id uint) error
	Count() (int64, error)
}

// {{.EntityNameLower}}Repository implements {{.EntityName}}Repository interface
type {{.EntityNameLower}}Repository struct {
	db *gorm.DB
}

// New{{.EntityName}}Repository creates a new {{.EntityNameLower}} repository instance
func New{{.EntityName}}Repository(db *gorm.DB) {{.EntityName}}Repository {
	return &{{.EntityNameLower}}Repository{
		db: db,
	}
}

// Create creates a new {{.EntityNameLower}}
func (r *{{.EntityNameLower}}Repository) Create({{.EntityNameLower}} *model.{{.EntityName}}) error {
	return r.db.Create({{.EntityNameLower}}).Error
}

// GetByID retrieves a {{.EntityNameLower}} by ID
func (r *{{.EntityNameLower}}Repository) GetByID(id uint) (*model.{{.EntityName}}, error) {
	var {{.EntityNameLower}} model.{{.EntityName}}
	err := r.db.First(&{{.EntityNameLower}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.EntityNameLower}}, nil
}

// GetAll retrieves all {{.EntityNameLower}}s with pagination
func (r *{{.EntityNameLower}}Repository) GetAll(limit, offset int) ([]model.{{.EntityName}}, error) {
	var {{.EntityNameLower}}s []model.{{.EntityName}}
	err := r.db.Limit(limit).Offset(offset).Find(&{{.EntityNameLower}}s).Error
	return {{.EntityNameLower}}s, err
}

// Update updates an existing {{.EntityNameLower}}
func (r *{{.EntityNameLower}}Repository) Update({{.EntityNameLower}} *model.{{.EntityName}}) error {
	return r.db.Save({{.EntityNameLower}}).Error
}

// Delete soft deletes a {{.EntityNameLower}}
func (r *{{.EntityNameLower}}Repository) Delete(id uint) error {
	return r.db.Delete(&model.{{.EntityName}}{}, id).Error
}

// Count returns the total number of {{.EntityNameLower}}s
func (r *{{.EntityNameLower}}Repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.{{.EntityName}}{}).Count(&count).Error
	return count, err
}
`

	return g.generateFile("src/repository/"+g.EntityNameLower+".go", tmpl, g.prepareRepositoryData())
}

func (g *Generator) GenerateUsecase() error {
	tmpl := `package usecase

import (
	"boilerblade/src/dto"
	"boilerblade/src/model"
	"boilerblade/src/repository"
	"errors"
	"math"
)

// {{.EntityName}}Usecase defines the interface for {{.EntityNameLower}} business logic
type {{.EntityName}}Usecase interface {
	Create{{.EntityName}}(req *dto.Create{{.EntityName}}Request) (*dto.{{.EntityName}}Response, error)
	Get{{.EntityName}}ByID(id uint) (*dto.{{.EntityName}}Response, error)
	GetAll{{.EntityName}}s(limit, offset int) (*dto.{{.EntityName}}ListResponse, error)
	Update{{.EntityName}}(id uint, req *dto.Update{{.EntityName}}Request) (*dto.{{.EntityName}}Response, error)
	Delete{{.EntityName}}(id uint) error
}

// {{.EntityNameLower}}Usecase implements {{.EntityName}}Usecase interface
type {{.EntityNameLower}}Usecase struct {
	{{.EntityNameLower}}Repo repository.{{.EntityName}}Repository
}

// New{{.EntityName}}Usecase creates a new {{.EntityNameLower}} usecase instance
func New{{.EntityName}}Usecase({{.EntityNameLower}}Repo repository.{{.EntityName}}Repository) {{.EntityName}}Usecase {
	return &{{.EntityNameLower}}Usecase{
		{{.EntityNameLower}}Repo: {{.EntityNameLower}}Repo,
	}
}

// Create{{.EntityName}} creates a new {{.EntityNameLower}}
func (uc *{{.EntityNameLower}}Usecase) Create{{.EntityName}}(req *dto.Create{{.EntityName}}Request) (*dto.{{.EntityName}}Response, error) {
	// Create {{.EntityNameLower}} model
	{{.EntityNameLower}} := &model.{{.EntityName}}{
		// TODO: Map fields from request to model
	}

	// Save to database
	if err := uc.{{.EntityNameLower}}Repo.Create({{.EntityNameLower}}); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.{{.EntityName}}Response{
		ID: {{.EntityNameLower}}.ID,
		// TODO: Map other fields
		CreatedAt: {{.EntityNameLower}}.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: {{.EntityNameLower}}.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Get{{.EntityName}}ByID retrieves a {{.EntityNameLower}} by ID
func (uc *{{.EntityNameLower}}Usecase) Get{{.EntityName}}ByID(id uint) (*dto.{{.EntityName}}Response, error) {
	{{.EntityNameLower}}, err := uc.{{.EntityNameLower}}Repo.GetByID(id)
	if err != nil {
		return nil, errors.New("{{.EntityNameLower}} not found")
	}

	return &dto.{{.EntityName}}Response{
		ID: {{.EntityNameLower}}.ID,
		// TODO: Map other fields
		CreatedAt: {{.EntityNameLower}}.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: {{.EntityNameLower}}.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetAll{{.EntityName}}s retrieves all {{.EntityNameLower}}s with pagination
func (uc *{{.EntityNameLower}}Usecase) GetAll{{.EntityName}}s(limit, offset int) (*dto.{{.EntityName}}ListResponse, error) {
	// Validate pagination
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Get {{.EntityNameLower}}s and total count
	{{.EntityNameLower}}s, err := uc.{{.EntityNameLower}}Repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.{{.EntityNameLower}}Repo.Count()
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	{{.EntityNameLower}}Responses := make([]dto.{{.EntityName}}Response, len({{.EntityNameLower}}s))
	for i, {{.EntityNameLower}} := range {{.EntityNameLower}}s {
		{{.EntityNameLower}}Responses[i] = dto.{{.EntityName}}Response{
			ID: {{.EntityNameLower}}.ID,
			// TODO: Map other fields
			CreatedAt: {{.EntityNameLower}}.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: {{.EntityNameLower}}.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.{{.EntityName}}ListResponse{
		{{.EntityName}}s:      {{.EntityNameLower}}Responses,
		Total:      total,
		Limit:      limit,
		Offset:     offset,
		TotalPages: totalPages,
	}, nil
}

// Update{{.EntityName}} updates an existing {{.EntityNameLower}}
func (uc *{{.EntityNameLower}}Usecase) Update{{.EntityName}}(id uint, req *dto.Update{{.EntityName}}Request) (*dto.{{.EntityName}}Response, error) {
	// Get existing {{.EntityNameLower}}
	{{.EntityNameLower}}, err := uc.{{.EntityNameLower}}Repo.GetByID(id)
	if err != nil {
		return nil, errors.New("{{.EntityNameLower}} not found")
	}

	// TODO: Update fields if provided

	// Save updates
	if err := uc.{{.EntityNameLower}}Repo.Update({{.EntityNameLower}}); err != nil {
		return nil, err
	}

	// Return response DTO
	return &dto.{{.EntityName}}Response{
		ID: {{.EntityNameLower}}.ID,
		// TODO: Map other fields
		CreatedAt: {{.EntityNameLower}}.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: {{.EntityNameLower}}.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Delete{{.EntityName}} deletes a {{.EntityNameLower}}
func (uc *{{.EntityNameLower}}Usecase) Delete{{.EntityName}}(id uint) error {
	// Check if {{.EntityNameLower}} exists
	_, err := uc.{{.EntityNameLower}}Repo.GetByID(id)
	if err != nil {
		return errors.New("{{.EntityNameLower}} not found")
	}

	// Delete {{.EntityNameLower}}
	return uc.{{.EntityNameLower}}Repo.Delete(id)
}
`

	return g.generateFile("src/usecase/"+g.EntityNameLower+".go", tmpl, g.prepareUsecaseData())
}

func (g *Generator) GenerateHandler() error {
	tmpl := `package handler

import (
	"boilerblade/helper"
	"boilerblade/src/dto"
	"boilerblade/src/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// {{.EntityName}}Handler handles HTTP requests for {{.EntityNameLower}} operations
type {{.EntityName}}Handler struct {
	{{.EntityNameLower}}Usecase usecase.{{.EntityName}}Usecase
	validator                   *validator.Validate
}

// New{{.EntityName}}Handler creates a new {{.EntityName}}Handler instance
func New{{.EntityName}}Handler({{.EntityNameLower}}Usecase usecase.{{.EntityName}}Usecase) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{
		{{.EntityNameLower}}Usecase: {{.EntityNameLower}}Usecase,
		validator:                   validator.New(),
	}
}

// RegisterRoutes registers all {{.EntityNameLower}} routes to the router group
func (h *{{.EntityName}}Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/{{.RouteName}}", h.GetAll{{.EntityName}}s)
	router.Get("/{{.RouteName}}/:id", h.Get{{.EntityName}})
	router.Post("/{{.RouteName}}", h.Create{{.EntityName}})
	router.Put("/{{.RouteName}}/:id", h.Update{{.EntityName}})
	router.Delete("/{{.RouteName}}/:id", h.Delete{{.EntityName}})
}

// Create{{.EntityName}} handles the creation of a new {{.EntityNameLower}}
func (h *{{.EntityName}}Handler) Create{{.EntityName}}(c *fiber.Ctx) error {
	var req dto.Create{{.EntityName}}Request
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrBadRequest("Invalid request body"))
	}

	{{.EntityNameLower}}Response, err := h.{{.EntityNameLower}}Usecase.Create{{.EntityName}}(&req)
	if err != nil {
		return helper.HandleUsecaseError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON({{.EntityNameLower}}Response)
}

// Get{{.EntityName}} handles retrieving a {{.EntityNameLower}} by ID
func (h *{{.EntityName}}Handler) Get{{.EntityName}}(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid {{.EntityNameLower}} ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrBadRequest("Invalid {{.EntityNameLower}} ID"))
	}

	{{.EntityNameLower}}Response, err := h.{{.EntityNameLower}}Usecase.Get{{.EntityName}}ByID(uint(id))
	if err != nil {
		return helper.HandleUsecaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON({{.EntityNameLower}}Response)
}

// GetAll{{.EntityName}}s handles retrieving all {{.EntityNameLower}}s with pagination
func (h *{{.EntityName}}Handler) GetAll{{.EntityName}}s(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		offset = 0
	}

	{{.EntityNameLower}}ListResponse, err := h.{{.EntityNameLower}}Usecase.GetAll{{.EntityName}}s(limit, offset)
	if err != nil {
		return helper.HandleUsecaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON({{.EntityNameLower}}ListResponse)
}

// Update{{.EntityName}} handles updating an existing {{.EntityNameLower}}
func (h *{{.EntityName}}Handler) Update{{.EntityName}}(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid {{.EntityNameLower}} ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrBadRequest("Invalid {{.EntityNameLower}} ID"))
	}

	var req dto.Update{{.EntityName}}Request
	if err := c.BodyParser(&req); err != nil {
		helper.LogError("Failed to parse request body", err, c.Path(), nil)
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrBadRequest("Invalid request body"))
	}

	{{.EntityNameLower}}Response, err := h.{{.EntityNameLower}}Usecase.Update{{.EntityName}}(uint(id), &req)
	if err != nil {
		return helper.HandleUsecaseError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON({{.EntityNameLower}}Response)
}

// Delete{{.EntityName}} handles soft deleting a {{.EntityNameLower}}
func (h *{{.EntityName}}Handler) Delete{{.EntityName}}(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		helper.LogError("Invalid {{.EntityNameLower}} ID parameter", err, c.Path(), map[string]interface{}{"id_param": c.Params("id")})
		return c.Status(fiber.StatusBadRequest).JSON(helper.ErrBadRequest("Invalid {{.EntityNameLower}} ID"))
	}

	if err := h.{{.EntityNameLower}}Usecase.Delete{{.EntityName}}(uint(id)); err != nil {
		return helper.HandleUsecaseError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
`

	return g.generateFile("src/handler/"+g.EntityNameLower+".go", tmpl, g.prepareHandlerData())
}

func (g *Generator) GenerateDTO() error {
	tmpl := `package dto

import "boilerblade/src/model"

// Create{{.EntityName}}Request represents the request payload for creating a {{.EntityNameLower}}
type Create{{.EntityName}}Request struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.NameLower}}\" validate:\"{{.ValidateTag}}\"`" + `
{{end}}}

// Update{{.EntityName}}Request represents the request payload for updating a {{.EntityNameLower}}
type Update{{.EntityName}}Request struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.NameLower}}\" validate:\"omitempty,{{.ValidateTag}}\"`" + `
{{end}}}

// {{.EntityName}}Response represents the response payload for a single {{.EntityNameLower}}
type {{.EntityName}}Response struct {
	ID        uint   ` + "`json:\"id\"`" + `
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.NameLower}}\"`" + `
{{end}}	CreatedAt string ` + "`json:\"created_at\"`" + `
	UpdatedAt string ` + "`json:\"updated_at\"`" + `
}

// {{.EntityName}}ListResponse represents the response payload for a list of {{.EntityNameLower}}s with pagination
type {{.EntityName}}ListResponse struct {
	{{.EntityName}}s      []{{.EntityName}}Response ` + "`json:\"{{.EntityNameLower}}s\"`" + `
	Total      int64          ` + "`json:\"total\"`" + `
	Limit      int            ` + "`json:\"limit\"`" + `
	Offset     int            ` + "`json:\"offset\"`" + `
	TotalPages int            ` + "`json:\"total_pages\"`" + `
}

// To{{.EntityName}}Response converts a model.{{.EntityName}} to {{.EntityName}}Response DTO
func To{{.EntityName}}Response({{.EntityNameLower}} *model.{{.EntityName}}) {{.EntityName}}Response {
	return {{.EntityName}}Response{
		ID: {{.EntityNameLower}}.ID,
		// TODO: Map other fields
		CreatedAt: {{.EntityNameLower}}.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: {{.EntityNameLower}}.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
`

	return g.generateFile("src/dto/"+g.EntityNameLower+".go", tmpl, g.prepareDTOData())
}

func (g *Generator) GenerateAll() error {
	layers := []func() error{
		g.GenerateModel,
		g.GenerateDTO,
		g.GenerateRepository,
		g.GenerateUsecase,
		g.GenerateHandler,
	}

	for _, layer := range layers {
		if err := layer(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generateFile(filePath string, tmplStr string, data interface{}) error {
	// Create directory if not exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	}

	// Parse template
	tmpl, err := template.New("code").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Helper methods to prepare data for templates
func (g *Generator) prepareModelData() map[string]interface{} {
	fields := make([]map[string]interface{}, len(g.Fields))
	for i, f := range g.Fields {
		gormTag := "not null"
		if f.Tag != "" {
			gormTag = f.Tag
		}
		fields[i] = map[string]interface{}{
			"Name":      f.Name,
			"Type":      f.Type,
			"NameLower": strings.ToLower(f.Name[:1]) + f.Name[1:],
			"GormTag":   gormTag,
		}
	}
	return map[string]interface{}{
		"EntityName":      g.EntityName,
		"EntityNameLower": g.EntityNameLower,
		"TableName":       g.getTableName(),
		"Fields":          fields,
	}
}

func (g *Generator) prepareRepositoryData() map[string]interface{} {
	return map[string]interface{}{
		"EntityName":      g.EntityName,
		"EntityNameLower": g.EntityNameLower,
	}
}

func (g *Generator) prepareUsecaseData() map[string]interface{} {
	return map[string]interface{}{
		"EntityName":      g.EntityName,
		"EntityNameLower": g.EntityNameLower,
	}
}

func (g *Generator) prepareHandlerData() map[string]interface{} {
	return map[string]interface{}{
		"EntityName":      g.EntityName,
		"EntityNameLower": g.EntityNameLower,
		"RouteName":       g.getRouteName(),
	}
}

func (g *Generator) prepareDTOData() map[string]interface{} {
	fields := make([]map[string]interface{}, len(g.Fields))
	for i, f := range g.Fields {
		validateTag := "required"
		if f.Tag != "" {
			validateTag = f.Tag
		}
		fields[i] = map[string]interface{}{
			"Name":        f.Name,
			"Type":        f.Type,
			"NameLower":   strings.ToLower(f.Name[:1]) + f.Name[1:],
			"ValidateTag": validateTag,
		}
	}
	return map[string]interface{}{
		"EntityName":      g.EntityName,
		"EntityNameLower": g.EntityNameLower,
		"Fields":          fields,
	}
}

func (g *Generator) getTableName() string {
	// Convert EntityName to plural snake_case
	// e.g., Product -> products, OrderItem -> order_items
	return strings.ToLower(g.EntityNameLower) + "s"
}

func (g *Generator) getRouteName() string {
	// Convert EntityName to plural lowercase
	// e.g., Product -> products
	return strings.ToLower(g.EntityNameLower) + "s"
}
