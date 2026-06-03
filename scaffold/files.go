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

func AppendToMakefile(projectPath string) {
	f, err := os.OpenFile(projectPath+"/Makefile", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open Makefile: %v", err)
	}
	defer f.Close()

	content := `generate:
		cd sqlc && sqlc generate
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
