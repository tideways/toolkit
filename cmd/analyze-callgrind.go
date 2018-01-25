package cmd

import (
	"fmt"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(analyzeCallgrindCmd)
	analyzeCallgrindCmd.Flags().StringVarP(&field, "dimension", "d", "excl_wt", "Dimension to view/sort (wt, excl_wt)")
	analyzeCallgrindCmd.Flags().Float32VarP(&minPercent, "min", "m", 1, "Display items having minimum percentage (default 1%) of --dimension, with respect to max value")
}

var analyzeCallgrindCmd = &cobra.Command{
	Use:   "analyze-callgrind filepaths...",
	Short: "Parse the output of callgrind outputs into a sorted tabular output.",
	Long:  `Parse the output of callgrind outputs into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  analyzeCallgrind,
}

func analyzeCallgrind(cmd *cobra.Command, args []string) error {
	maps := make([]*xhprof.PairCallMap, 0, len(args))
	for _, arg := range args {
		f := xhprof.NewFile(arg, "callgrind")
		m, err := f.GetPairCallMap()
		if err != nil {
			return err
		}

		maps = append(maps, m)
	}

	avgMap := xhprof.AvgPairCallMaps(maps)
	profile := avgMap.Flatten()

	fieldInfo, ok := fieldsMap[field]
	if !ok {
		fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
		field = "excl_wt"
		fieldInfo = fieldsMap[field]
	}

	profile.SortBy(fieldInfo.Name)

	// Change default to 10 for exclusive fields, only when user
	// hasn't manually provided 1%
	if strings.HasPrefix(field, "excl_") && !cmd.Flags().Changed("min") {
		minPercent = float32(10)
	}
	minPercent = minPercent / 100.0
	minValue := minPercent * profile.Calls[0].GetFloat32Field(fieldInfo.Name)
	profile = profile.SelectGreater(fieldInfo.Name, minValue)
	err := renderProfile(profile, field, fieldInfo, minValue)
	if err != nil {
		return err
	}

	return nil
}
