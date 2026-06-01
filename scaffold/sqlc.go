package scaffold

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func GenerateSqlcFiles(projectPath, database string) {
	os.Mkdir(projectPath+"/sqlc", 0755)
	cmd := exec.Command("touch",
		projectPath+"/sqlc/schema.sql",
		projectPath+"/sqlc/query.sql")
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to create sqlc schema and queries files: %v", err)
	}
	var content string
	if database != "postgres" {
		content = fmt.Sprintf(`
version: "2"
sql:
  - engine: "%v"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        out: "../db"

	`, database)
	} else {
		content = `
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        out: "../db"
`
	}

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
