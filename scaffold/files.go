package scaffold

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func modifyDatabase(projectPath string) {
	dirName := filepath.Base(projectPath)

	content, err := os.ReadFile(projectPath + "/.env")
	if err != nil {
		log.Fatalf("failed to read .env file at %s: %v", projectPath+"/.env", err)
	}
	data := string(content)

	data = strings.Replace(data, "import (", "import (\n\tdb \""+dirName+"/db\"", 1)

  data = strings.Replace(data, "sqlc.Queries", "db.Queries", -1)

}
