package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"goproto/scaffold"

	"charm.land/huh/v2"
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
	var structure string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What Project Structure Do You Want?").
				Options(
					huh.NewOption("API", "api"),
					huh.NewOption("Normal", "normal"),
				).
				Value(&structure)))
	err := form.Run()
	if err != nil {
		log.Fatalf("failed to run prompt form: %v", err)
	}

	switch structure {
	case "api":
		log.Println("Generating...")
		createAPIStructure()
	case "normal":
		log.Println("Normal structure selected")
		createNormalStructure()
	default:
		log.Println("Invalid structure selected")
	}
}

func createAPIStructure() {
	framework, database := blueprint()

	initInTmp(framework, database)

	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory in createAPIStructure: %v", err)
	}

	dirName := filepath.Base(projectPath)

	scaffold.ModifyDatabase(projectPath, dirName, database)

	scaffold.GenerateSqlcFiles(projectPath, database)
	scaffold.GenerateReadmeFile(projectPath)
	if database == "postgres" {
		database = "pgx"
	}
	scaffold.GenerateGooseFiles(projectPath, database)
	scaffold.AppendToMakefile(projectPath)

	// handlers
	err = os.Mkdir(projectPath+"/internal/server/handlers", 0755)
	if err != nil {
		log.Fatalf("failed to create handlers directory: %v", err)
	}
}

func createNormalStructure() {
	log.Println("Creating Normal structure...")
	// Implement the logic to create a normal project structure
}

func blueprint() (frmwrk, dbase string) {
	var framework string
	var database string
	if err := huh.NewSelect[string]().
		Title("What Framework Do You Want?").
		Options(
			huh.NewOption("Chi", "chi"),
			huh.NewOption("Gin", "gin"),
			huh.NewOption("Net/http", "net/http"),
		).
		Value(&framework).Run(); err != nil {
		log.Fatalf("failed to select framework: %v", err)
	}

	if err := huh.NewSelect[string]().
		Title("What Type of Database Do You Want?").
		Options(
			huh.NewOption("PostgreSQL", "postgres"),
			huh.NewOption("SQLite", "sqlite"),
		).
		Value(&database).Run(); err != nil {
		log.Fatalf("failed to select database: %v", err)
	}

	return framework, database
}

func initInTmp(framework, database string) {
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
		"--framework", framework,
		"--driver", database,
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

	if database == "postgres" {
		scaffold.AddDatabaseTest(projectPath)
	}
}
