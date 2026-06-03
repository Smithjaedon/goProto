package scaffold

import (
	"log"
	"os"
)

func GenerateGooseFiles(projectPath string) {
	err := os.Mkdir(projectPath+"/migrations", 0755)
	if err != nil {
		log.Fatalf("failed to create migrations directory '%s': %v", projectPath+"/migrations", err)
	}

	f, err := os.OpenFile(projectPath+"/.env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open or create .env file at %s: %v", projectPath+"/.env", err)
	}
	defer f.Close()

	content := `GOOSE_DRIVER=pgx
GOOSE_DBSTRING=postgresql://username:password@localhost:port/dbname?sslmode=disable&search_path=public
GOOSE_MIGRATION_DIR=./migrations/
`

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalf("failed to write goose env content to %s: %v", projectPath+"/.env", err)
	}
}
