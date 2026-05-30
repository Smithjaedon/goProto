package cmd

import (
	"fmt"
	"log"
	"os"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
)

func PlanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Plan the project structure and dependencies",
		Run:   runPlan,
	}
}

func runPlan(cmd *cobra.Command, args []string) {

	var planProblems string
	var inputsNeeded string
	var outputsNeeded string
	var planningDecision bool
	var mvpDecision string

	if err := huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&planningDecision).Run(); err != nil {
		return
	}

	if !planningDecision {
		fmt.Println("No plan will be created.")
		return
	}

	if err := huh.NewText().
		Title("Define the problem clearly.").
		Value(&planProblems).Run(); err != nil {
		return
	}

	if err := huh.NewText().
		Title("What are the inputs needed for this project?. ").
		Value(&inputsNeeded).Run(); err != nil {
		return
	}

	if err := huh.NewText().
		Title("What are the outputs needed for this project?. ").
		Value(&outputsNeeded).Run(); err != nil {
		return
	}

	if err := huh.NewText().
		Title("What does \"done\" mean for this project?. ").
		Value(&mvpDecision).Run(); err != nil {
		return
	}

	content := fmt.Sprintf(`
# Project Plan

## Problem Definition
%v

## Inputs Needed
%v

## Outputs Needed
%v

## MVP Decision
%v`, planProblems, inputsNeeded, outputsNeeded, mvpDecision)

	newfile, err := os.Create("plan.md")
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()
	newfile.WriteString(content)
}
