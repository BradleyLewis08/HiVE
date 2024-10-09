package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command {
	Use: "hive",
	Short: "Hive CLI",
	Long: "CLI for managing HiVE",
}

func Execute() error {
	return rootCmd.Execute()
}