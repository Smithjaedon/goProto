package scaffold

import (
	"embed"
	"log"
	"os"
	"text/template"
)

//go:embed templates
var templateFiles embed.FS

type ProjectConfig struct {
	ModuleName string
	Database   string
}

func ModifyDatabase(projectPath, moduleName string) {
	config := ProjectConfig{
		ModuleName: moduleName,
		Database:   "postgres",
	}

	generateFile("templates/sqlc/database.go.tmpl", projectPath+"/internal/database/database.go", config)
	generateFile("templates/sqlc/server.go.tmpl", projectPath+"/internal/server/server.go", config)
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

func AddDatabaseTest(projectPath string) {
	config := ProjectConfig{
		ModuleName: "",
		Database:   "",
	}
	generateFile("templates/database_test.go.tmpl", projectPath+"/internal/database/database_test.go", config)
}
