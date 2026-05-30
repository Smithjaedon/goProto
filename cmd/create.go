package cmd

import (
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
	cmd := exec.Command("go-blueprint", "create")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func createNormalStructure() {
	log.Println("Creating Normal structure...")
	// Implement the logic to create a normal project structure
}
