package scaffold

import (
	"log"
	"os"
)

func GenerateSqlcFiles(projectPath string) {
	dirs := []string{
		projectPath + "/sqlc",
		projectPath + "/sqlc/schemas",
		projectPath + "/sqlc/queries",
	}

	for _, dir := range dirs {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatalf("failed to create directory '%s': %v", dir, err)
		}
	}

	// Create empty placeholder files — GenerateDatabaseFiles will write the real content
	placeholders := []string{
		projectPath + "/sqlc/schemas/users.sql",
		projectPath + "/sqlc/queries/users.sql",
	}

	for _, path := range placeholders {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("failed to create file '%s': %v", path, err)
		}
		f.Close()
	}

	content := `version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "schemas/"
    gen:
      go:
        package: "db"
        out: "../db"
`

	file, err := os.Create(projectPath + "/sqlc/sqlc.yaml")
	if err != nil {
		log.Fatalf("failed to create sqlc.yaml at %s: %v", projectPath+"/sqlc/sqlc.yaml", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("failed to write sqlc.yaml: %v", err)
	}
}
