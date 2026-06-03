package scaffold

import (
	"log"
	"os"
)

func GenerateAuthFiles(projectPath string) {
	err := os.Mkdir(projectPath+"/internal/middleware", 0755)
	if err != nil {
		log.Fatalf("failed to create middleware directory: %v", err)
	}
}
