package cmd

import (
	"github.com/spf13/cobra"

	"github.com/go-gophers/gophers/utils/log"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate main file to run tests.",

	Run: func(cmd *cobra.Command, args []string) {
		log.Default.Debug = debugF
		if len(args) != 0 {
			log.Fatalf("expected 0 arguments, got %d", len(args))
		}

		output, _ := cmd.Flags().GetString("output")
		testPackage(WD, output)
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().StringP("output", "o", "main-test.go", "output file name")
}

func testPackage(dir string, output string) {
	data := importPackage(dir)
	data.Test = true
	renderTemplate(data, output)
}
