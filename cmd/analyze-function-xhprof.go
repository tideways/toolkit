package cmd

import (
	"fmt"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(analyzeFunctionXhprofCmd)
	analyzeFunctionXhprofCmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of rows to display")
}

var analyzeFunctionXhprofCmd = &cobra.Command{
	Use:   "analyze-function-xhprof function filepaths...",
	Short: "Report parents and children of a function based on a JSON serialized XHProf file.",
	Long:  "Report parents and children of a function based on a JSON serialized XHProf file.",
	Args:  cobra.MinimumNArgs(2),
	RunE:  analyzeFunctionXhprof,
}

func analyzeFunctionXhprof(cmd *cobra.Command, args []string) error {
	function := args[0]
	maps := make([]*xhprof.PairCallMap, 0, len(args[1:]))
	for _, arg := range args[1:] {
		f := xhprof.NewFile(arg, "xhprof")
		m, err := f.GetPairCallMap()
		if err != nil {
			return err
		}

		maps = append(maps, m)
	}

	avgMap := xhprof.AvgPairCallMaps(maps)
	family := avgMap.ComputeNearestFamily(function)
	parentsProfile := family.Parents.Flatten()
	childrenProfile := family.Children.Flatten()

	field = "wt"
	fieldInfo := fieldsMap[field]

	fmt.Printf("Parents of %s:\n", function)
	err := renderProfile(parentsProfile, field, fieldInfo, limit, 0)
	if err != nil {
		return err
	}

	fmt.Printf("Children of %s:\n", function)
	err = renderProfile(childrenProfile, field, fieldInfo, limit, 0)
	if err != nil {
		return err
	}

	return nil
}
