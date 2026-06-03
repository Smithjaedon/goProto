package main

import (
	"fmt"
	"os"

	"goproto/cmd"

	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "goProto",
		Short: "A Go Project Scaffolding Tool",
	}

	rootCmd.AddCommand(cmd.CreateCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
