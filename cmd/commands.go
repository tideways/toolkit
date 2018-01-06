package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tk",
	Short: "Tideways Toolkit is a collection of tools to interact with PHP",
	Long:  `The Tideways Toolkit (tk) is a collection of commandline tools to interact with
	PHP and perform various debugging, profiling and introspection jobs by
	interacting with PHP or with debugging extensions for PHP.`,
}

var version string

func Execute(v string) {
	version = v

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		if cerr, ok := err.(*CommandError); ok {
			os.Exit(cerr.ExitStatus())
		}

		os.Exit(1)
	}
}

type CommandError struct {
	s          string
	exitStatus int
}

func NewCommandError(exitStatus int, s string) *CommandError {
	err := new(CommandError)
	err.s = s
	err.exitStatus = exitStatus

	return err
}

func (err *CommandError) Error() string {
	return err.s
}

func (err *CommandError) ExitStatus() int {
	return err.exitStatus
}
