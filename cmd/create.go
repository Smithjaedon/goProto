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
	initInTmp()

	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory in createAPIStructure: %v", err)
	}

	dirName := filepath.Base(projectPath)

	scaffold.ModifyDatabase(projectPath, dirName)

	scaffold.GenerateSqlcFiles(projectPath)
	scaffold.GenerateReadmeFile(projectPath)
	scaffold.GenerateGooseFiles(projectPath)
	scaffold.AppendToMakefile(projectPath)

	// handlers
	err = os.Mkdir(projectPath+"/internal/server/handlers", 0755)
	if err != nil {
		log.Fatalf("failed to create handlers directory: %v", err)
	}
}

func initInTmp() {
	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	dirName := filepath.Base(projectPath)

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

	scaffold.AddDatabaseTest(projectPath)
}
