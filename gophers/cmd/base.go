// Package cmd contains implementation of gophers tool commands.
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/go-gophers/gophers/utils/log"
)

var (
	GoBin    string
	GoFmtBin string
	WD       string

	debugF bool

	RootCmd = &cobra.Command{
		Use:   "gophers",
		Short: "Gophers tool runs functional and load tests.",
	}
)

func init() {
	var err error

	GoBin, err = exec.LookPath("go")
	if err != nil {
		log.Fatal(err)
	}

	GoFmtBin, err = exec.LookPath("gofmt")
	if err != nil {
		log.Fatal(err)
	}

	WD, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	RootCmd.PersistentFlags().BoolVar(&debugF, "debug", false, "Enable debugging")
}
