package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of HiVE",
	Long:  `All software has versions. This is HiVE's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HiVE v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
