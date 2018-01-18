package cmd

import (
	"fmt"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(analyzeCallgrindCmd)
	analyzeCallgrindCmd.Flags().StringVarP(&field, "field", "f", "excl_wt", "Field to view/sort (wt, excl_wt)")
	analyzeCallgrindCmd.Flags().Float32VarP(&minPercent, "min", "m", 1, "Display items having minimum percentage (default 1%) of --field, with respect to main()")
}

var analyzeCallgrindCmd = &cobra.Command{
	Use:   "analyze-callgrind filepaths...",
	Short: "Parse the output of callgrind outputs into a sorted tabular output.",
	Long:  `Parse the output of callgrind outputs into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  analyzeCallgrind,
}

func analyzeCallgrind(cmd *cobra.Command, args []string) error {
	profiles := make([]*xhprof.Profile, 0, len(args))
	for _, arg := range args {
		f := xhprof.NewFile(arg, "callgrind")
		profile, err := f.GetProfile()
		if err != nil {
			return err
		}

		profiles = append(profiles, profile)
	}

	avgProfile := xhprof.AvgProfiles(profiles)

	fieldInfo, ok := fieldsMap[field]
	if !ok {
		fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
		field = "excl_wt"
		fieldInfo = fieldsMap[field]
	}

	minPercent = minPercent / 100.0
	err := renderProfile(avgProfile, field, fieldInfo, minPercent)
	if err != nil {
		return err
	}

	return nil
}
