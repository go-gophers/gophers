package cmd

import (
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Generate main file to run load tests.",

	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			race, _ := cmd.Flags().GetBool("race")
			weighted, _ := cmd.Flags().GetBool("weighted")
			loadPackage(arg, race, weighted)
		}
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)

	loadCmd.Flags().Bool("race", false, "enable data race detection")
	loadCmd.Flags().Bool("weighted", false, "generate file for weighted load tests")
}

func loadPackage(path string, race bool, weighted bool) {
	data := importPackage(path, race)
	data.Load = !weighted
	data.LoadWeighted = weighted
	renderTemplate(data, "main-load.go")
}
