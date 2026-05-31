package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"goProto/scaffold"

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
		log.Fatal(err)
	}

	switch structure {
	case "api":
		log.Println("API structure selected")
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
		log.Fatal(err)
	}
	scaffold.GenerateSqlcFiles(projectPath, database)
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
		).
		Value(&framework).Run(); err != nil {
		log.Fatal(err)
	}

	if err := huh.NewSelect[string]().
		Title("What Type of Database Do You Want?").
		Options(
			huh.NewOption("PostgreSQL", "pgx"),
			huh.NewOption("SQLite", "sqlite"),
		).
		Value(&database).Run(); err != nil {
		log.Fatal(err)
	}

	return framework, database
}

func initInTmp(framework, database string) {
	tmpDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go-blueprint", "create",
		"--name", "myproject",
		"--framework", framework,
		"--driver", database,
		"--git", "skip",
	)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	mv := exec.Command("bash", "-c", fmt.Sprintf("mv %s/myproject/* %s/ && mv %s/myproject/.* %s/ 2>/dev/null && rm -rf %s/myproject", tmpDir, projectPath, tmpDir, projectPath, tmpDir))
	if err := mv.Run(); err != nil {
		log.Fatal(err)
	}
}
