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

func ModifyDatabase(projectPath, moduleName, database string) {
	config := ProjectConfig{
		ModuleName: moduleName,
		Database:   database,
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
