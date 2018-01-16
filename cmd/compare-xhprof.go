package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(compareCmd)
	xhprofCmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of rows to display")
}

var (
	limit int
)

var compareCmd = &cobra.Command{
	Use:   "compare-xhprof [options]... filepaths...",
	Short: "Compare two JSON serialized XHProf outputs and display them in a sorted table.",
	Long:  `Compare two JSON serialized XHProf outputs and display them in a sorted table.`,
	Args:  cobra.ExactArgs(2),
	RunE:  compare,
}

func compare(cmd *cobra.Command, args []string) error {
	profiles := make([]*xhprof.Profile, 0, len(args))
	for _, arg := range args {
		rawData, err := ioutil.ReadFile(arg)
		if err != nil {
			return err
		}

		var data map[string]*xhprof.PairCall
		err = json.Unmarshal(rawData, &data)
		if err != nil {
			return err
		}

		profile, err := xhprof.Flatten(data)
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
