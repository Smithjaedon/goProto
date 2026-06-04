package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"goproto/scaffold"

	"github.com/spf13/cobra"
)

func CreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Initialize a new Go project",
		Run:   runCreate,
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	createAPIStructure()
}

func createAPIStructure() {
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory in createAPIStructure: %v", err)
	}

	dirName := filepath.Base(projectPath)

	initInTmp(projectPath, dirName)

	err = os.Mkdir(projectPath+"/internal/middleware", 0755)
	if err != nil {
		log.Fatalf("failed to create middleware directory: %v", err)
	}

	scaffold.GenerateSqlcFiles(projectPath)
	scaffold.GenerateGooseFiles(projectPath)
	scaffold.GenerateDatabaseFiles(projectPath, dirName)
	scaffold.GenerateAuthFiles(projectPath, dirName)
	scaffold.AuthDependencies(projectPath)
	scaffold.AddDatabaseTest(projectPath, dirName)
	scaffold.WriteMakefile(projectPath)
	scaffold.GenerateReadmeFile(projectPath)

	cmd := exec.Command("bash", "-c", "cd sqlc && sqlc generate")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run sqlc generate: %v", err)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run go mod tidy: %v", err)
	}

	fmt.Println("Project created successfully!")
}

func initInTmp(projectPath, dirName string) {
	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("go-blueprint", "create",
		"--name", dirName,
		"--framework", "gin",
		"--driver", "postgres",
		"--git", "skip",
	)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run go-blueprint command: %v", err)
	}

	mv := exec.Command("bash", "-c", fmt.Sprintf(
		"mv %s/%s/* %s/ && mv %s/%s/.* %s/ 2>/dev/null; rm -rf %s/%s",
		tmpDir, dirName, projectPath,
		tmpDir, dirName, projectPath,
		tmpDir, dirName,
	))
	if err := mv.Run(); err != nil {
		log.Fatalf("failed to move files from temporary blueprint: %v", err)
	}
}
