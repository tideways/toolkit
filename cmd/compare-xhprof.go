package cmd

import (
	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(compareXhprofCmd)
	compareXhprofCmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of rows to display")
}

var (
	limit int
)

var compareXhprofCmd = &cobra.Command{
	Use:   "compare-xhprof filepaths...",
	Short: "Compare two JSON serialized XHProf outputs and display them in a sorted table.",
	Long:  `Compare two JSON serialized XHProf outputs and display them in a sorted table.`,
	Args:  cobra.ExactArgs(2),
	RunE:  compareXhprof,
}

func compareXhprof(cmd *cobra.Command, args []string) error {
	profiles := make([]*xhprof.Profile, 0, len(args))
	for _, arg := range args {
		profile, err := xhprof.ParseFile(arg, false)
		if err != nil {
			return err
		}

		profiles = append(profiles, profile)
	}

	diff := profiles[0].Subtract(profiles[1])

	err := renderProfileDiff(diff, limit)
	if err != nil {
		return err
	}

	return nil
}
