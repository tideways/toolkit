package cmd

import (
	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(compareCallgrindCmd)
	compareCallgrindCmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of rows to display")
}

var compareCallgrindCmd = &cobra.Command{
	Use:   "compare-callgrind filepaths...",
	Short: "Compare two callgrind outputs and display them in a sorted table.",
	Long:  `Compare two callgrind outputs and display them in a sorted table.`,
	Args:  cobra.ExactArgs(2),
	RunE:  compareCallgrind,
}

func compareCallgrind(cmd *cobra.Command, args []string) error {
	profiles := make([]*xhprof.Profile, 0, len(args))
	for _, arg := range args {
		profile, err := xhprof.ParseFile(arg, true)
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
