package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate main file to run tests.",

	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			race, _ := cmd.Flags().GetBool("race")
			testPackage(arg, race)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().Bool("race", false, "enable data race detection")
}

func testPackage(path string, race bool) {
	data := importPackage(path, race)
	data.Test = true
	renderTemplate(data, "main-test.go")
}
