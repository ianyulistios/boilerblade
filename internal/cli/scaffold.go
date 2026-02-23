package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateNewProject scaffolds a new Boilerblade project
func CreateNewProject(projectName string) error {
	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check if directory already exists
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	fmt.Printf("Creating new Boilerblade project: %s\n", projectName)

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find the boilerblade template directory
	// Strategy: Look for go.mod with "module boilerblade" in current or parent directories
	templateDir := findBoilerbladeRoot(currentDir)
	if templateDir == "" {
		return fmt.Errorf("could not find boilerblade template directory. Please run this command from within the boilerblade project directory")
	}

	// Create project directory
	projectPath := filepath.Join(currentDir, projectName)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Files and directories to copy (excluding certain files)
	excludeDirs := map[string]bool{
		".git":        true,
		"bin":         true,
		"node_modules": true,
		"vendor":      true,
	}

	excludeFiles := map[string]bool{
		"boilerblade.exe": true,
		"generate.exe":    true,
		".env":            true, // Don't copy .env; create from .env.example
	}

	// Copy files and directories
	if err := copyProjectFiles(templateDir, projectPath, excludeDirs, excludeFiles); err != nil {
		// Cleanup on error
		os.RemoveAll(projectPath)
		return fmt.Errorf("failed to copy project files: %w", err)
	}

	// Update go.mod with new module name
	if err := updateGoMod(projectPath, projectName); err != nil {
		return fmt.Errorf("failed to update go.mod: %w", err)
	}

	// Create .env from .env.example
	if err := createEnvFile(projectPath, projectName); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	// Update config files with project name
	if err := updateConfigFiles(projectPath, projectName); err != nil {
		return fmt.Errorf("failed to update config files: %w", err)
	}

	fmt.Printf("✓ Project '%s' created successfully!\n", projectName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod download")
	fmt.Println("  # .env was created from .env.example — edit with your credentials if needed")
	fmt.Println("  go run main.go")

	return nil
}

func copyProjectFiles(src, dst string, excludeDirs, excludeFiles map[string]bool) error {
	// Get absolute paths to avoid issues
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	return filepath.Walk(srcAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source
		relPath, err := filepath.Rel(srcAbs, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Check if directory should be excluded
		dirName := filepath.Base(path)
		if info.IsDir() && excludeDirs[dirName] {
			return filepath.SkipDir
		}

		// Check if file should be excluded
		if !info.IsDir() && excludeFiles[dirName] {
			return nil
		}

		// Skip if path contains excluded directories
		pathParts := strings.Split(relPath, string(filepath.Separator))
		for _, part := range pathParts {
			if excludeDirs[part] {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Create destination path
		dstPath := filepath.Join(dstAbs, relPath)

		// Prevent copying into itself
		if strings.HasPrefix(dstPath, srcAbs) && dstPath != srcAbs {
			return nil
		}

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if needed
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}

func updateGoMod(projectPath, projectName string) error {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}

	// Replace module name
	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "module boilerblade", fmt.Sprintf("module %s", projectName))

	return os.WriteFile(goModPath, []byte(contentStr), 0644)
}

func createEnvFile(projectPath, projectName string) error {
	envExamplePath := filepath.Join(projectPath, ".env.example")
	envPath := filepath.Join(projectPath, ".env")

	var content []byte
	var err error
	if content, err = os.ReadFile(envExamplePath); err != nil {
		// Fallback to embedded template if .env.example was not copied
		content = []byte(envExampleContent)
	}

	contentStr := string(content)
	contentStr = strings.ReplaceAll(contentStr, "FIBER_APP_NAME=boilerblade", fmt.Sprintf("FIBER_APP_NAME=%s", projectName))
	contentStr = strings.ReplaceAll(contentStr, "DB_NAME=boilerblade", fmt.Sprintf("DB_NAME=%s", projectName))

	return os.WriteFile(envPath, []byte(contentStr), 0644)
}

func updateConfigFiles(projectPath, projectName string) error {
	// Update main.go if it references the app name
	mainGoPath := filepath.Join(projectPath, "main.go")
	if _, err := os.Stat(mainGoPath); err == nil {
		content, err := os.ReadFile(mainGoPath)
		if err != nil {
			return err
		}

		contentStr := string(content)
		// Replace any references to boilerblade module
		contentStr = strings.ReplaceAll(contentStr, "boilerblade/", fmt.Sprintf("%s/", projectName))

		if err := os.WriteFile(mainGoPath, []byte(contentStr), 0644); err != nil {
			return err
		}
	}

	// Update all Go files to replace module imports
	return updateModuleImports(projectPath, projectName)
}

func findBoilerbladeRoot(startDir string) string {
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			content, err := os.ReadFile(goModPath)
			if err == nil {
				contentStr := string(content)
				if strings.Contains(contentStr, "module boilerblade") {
					return dir
				}
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}
	return ""
}

func updateModuleImports(projectPath, projectName string) error {
	return filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)
		// Replace boilerblade imports with new project name
		contentStr = strings.ReplaceAll(contentStr, `"boilerblade/`, fmt.Sprintf(`"%s/`, projectName))
		contentStr = strings.ReplaceAll(contentStr, `"boilerblade"`, fmt.Sprintf(`"%s"`, projectName))

		return os.WriteFile(path, []byte(contentStr), 0644)
	})
}
