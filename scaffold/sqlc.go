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
		projectPath+"/sqlc/queries.sql")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf(`
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

	file, err := os.Create(projectPath + "/sqlc/sqlc.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}
