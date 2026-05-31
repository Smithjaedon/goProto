package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
	cmd := exec.Command("go-blueprint", "create", "--name", "myproject", "--framework", framework, "--database", database, "--git", "commit")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	projectPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	os.Mkdir(projectPath+"/sqlc", 0755)
	cmd = exec.Command("touch",
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

	file, err := os.Create(projectPath + "/sqlc.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
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
