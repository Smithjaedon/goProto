package scaffold

import (
	"embed"
	"log"
	"os"
	"os/exec"
	"text/template"
)

//go:embed templates
var templateFiles embed.FS

type ProjectConfig struct {
	ModuleName string
	Database   string
}

func GenerateDatabaseFiles(projectPath, moduleName string) {
	config := ProjectConfig{
		ModuleName: moduleName,
		Database:   "postgres",
	}

	generateFile("templates/sqlc/database.go.tmpl", projectPath+"/internal/database/database.go", config)
	generateFile("templates/sqlc/server.go.tmpl", projectPath+"/internal/server/server.go", config)
	generateFile("templates/sqlc/user_queries.go.tmpl", projectPath+"/sqlc/queries/users.sql", config)
	generateFile("templates/sqlc/user_schemas.go.tmpl", projectPath+"/sqlc/schemas/users.sql", config)
	generateFile("templates/auth/routers.go.tmpl", projectPath+"/internal/server/routes.go", config)
}

func GenerateAuthFiles(projectPath, moduleName string) {
	config := ProjectConfig{
		ModuleName: moduleName,
		Database:   "postgres",
	}

	generateFile("templates/auth/auth_handler.go.tmpl", projectPath+"/internal/server/auth_handler.go", config)
	generateFile("templates/auth/auth.go.tmpl", projectPath+"/internal/middleware/auth.go", config)
	generateFile("templates/add_users.go.tmpl", projectPath+"/migrations/00001_add_users.sql", config)
}

func generateFile(templatePath, outputPath string, config ProjectConfig) {
	tmpl, err := template.ParseFS(templateFiles, templatePath)
	if err != nil {
		log.Fatalf("failed to parse %s template: %v", outputPath, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create %s file: %v", outputPath, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, config)
	if err != nil {
		log.Fatalf("failed to execute %s template: %v", outputPath, err)
	}
}

func WriteMakefile(projectPath string) {
	f, err := os.OpenFile(projectPath+"/Makefile", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed to open Makefile: %v", err)
	}
	defer f.Close()

	content := `# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integration Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest generate goose_create

generate:
	cd sqlc && sqlc generate

goose_create:
	@read -p "Migration name: " name; goose create -s $$name sql
`

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalf("failed to write to Makefile: %v", err)
	}
}

func AddDatabaseTest(projectPath, moduleName string) {
	config := ProjectConfig{
		ModuleName: moduleName,
		Database:   "postgres",
	}
	generateFile("templates/database_test.go.tmpl", projectPath+"/internal/database/database_test.go", config)
}

func AuthDependencies(projectPath string) {
	deps := []string{
		"github.com/gin-gonic/gin",
		"golang.org/x/crypto/bcrypt",
		"github.com/google/uuid",
	}

	for _, dep := range deps {
		cmd := exec.Command("go", "get", dep)
		cmd.Dir = projectPath
		if err := cmd.Run(); err != nil {
			log.Fatalf("failed to install dependency %s: %v", dep, err)
		}
	}
}
