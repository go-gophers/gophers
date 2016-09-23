package cmd

import (
	"github.com/spf13/cobra"

	"github.com/go-gophers/gophers/utils/log"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Generate main file to run load tests.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Default.Debug = debugF
		if len(args) != 0 {
			log.Fatalf("expected 0 arguments, got %d", len(args))
		}

		output, _ := cmd.Flags().GetString("output")
		weighted, _ := cmd.Flags().GetBool("weighted")
		loadPackage(WD, output, weighted)
	},
}

func init() {
	RootCmd.AddCommand(loadCmd)

	loadCmd.Flags().StringP("output", "o", "main-load.go", "output file name")
	loadCmd.Flags().Bool("weighted", false, "generate file for weighted load tests")
}

func loadPackage(dir string, output string, weighted bool) {
	data := importPackage(dir)
	data.Load = !weighted
	data.LoadWeighted = weighted
	renderTemplate(data, output)
}
