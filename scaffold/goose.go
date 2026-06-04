package scaffold

import (
	"log"
	"os"
	"os/exec"
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

	cmd := exec.Command("openssl", "rand", "-hex", "32")
	cmd.Dir = projectPath
	secretKey, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to generate secret key: %v", err)
	}

	_, err = f.WriteString("SECRET_KEY=" + string(secretKey))
	if err != nil {
		log.Fatalf("failed to write secret key to .env file: %v", err)
	}
}
