package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(xhprofCmd)
	xhprofCmd.Flags().StringVarP(&field, "dimension", "d", "excl_wt", "Dimension to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory, io, excl_io)")
	xhprofCmd.Flags().Float32VarP(&minPercent, "min", "m", 1, "Display items having minimum percentage (default 1% for inclusive, and 10% for exclusive dimensions) of --dimension, with respect to main()")
	xhprofCmd.Flags().StringVarP(&outFile, "out-file", "o", "", "If provided, the path to store the resulting profile (e.g. after averaging)")
	xhprofCmd.Flags().StringVarP(&function, "function", "", "", "If provided, one table for parents, and one for children of this function will be displayed")
}

var (
	field      string
	minPercent float32
	outFile    string
	function   string
)

var xhprofCmd = &cobra.Command{
	Use:   "analyze-xhprof filepaths...",
	Short: "Parse the output of JSON serialized XHProf outputs into a sorted tabular output.",
	Long:  `Parse the output of JSON serialized XHProf outputs into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  analyzeXhprof,
}

func analyzeXhprof(cmd *cobra.Command, args []string) error {
	maps := make([]*xhprof.PairCallMap, 0, len(args))
	for _, arg := range args {
		f := xhprof.NewFile(arg, "xhprof")
		m, err := f.GetPairCallMap()
		if err != nil {
			return err
		}

		maps = append(maps, m)
	}

	avgMap := xhprof.AvgPairCallMaps(maps)
	if outFile != "" {
		fmt.Printf("Writing profile to %s\n", outFile)
		f := xhprof.NewFile(outFile, "xhprof")
		err := f.WritePairCallMap(avgMap)
		if err != nil {
			return err
		}
	}

	profile := avgMap.Flatten()

	// Change default to 10 for exclusive fields, only when user
	// hasn't manually provided 1%
	if strings.HasPrefix(field, "excl_") && !cmd.Flags().Changed("min") {
		minPercent = float32(10)
	}
	minPercent = minPercent / 100.0

	if function == "" {
		fieldInfo, ok := fieldsMap[field]
		if !ok {
			fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
			field = "excl_wt"
			fieldInfo = fieldsMap[field]
		}

		profile.SortBy(fieldInfo.Name)
		minValue := minPercent * profile.Calls[0].GetFloat32Field(fieldInfo.Name)
		profile = profile.SelectGreater(fieldInfo.Name, minValue)
		err := renderProfile(profile, field, fieldInfo, minValue)
		if err != nil {
			return err
		}
	} else {
		family := avgMap.ComputeNearestFamily(function)
		parentsProfile := family.Parents.Flatten()
		childrenProfile := family.Children.Flatten()

		field = "wt"
		fieldInfo := fieldsMap[field]
		minPercent = 0.1

		functionCall := profile.GetCall(function)
		if functionCall == nil {
			return errors.New("Profile doesn't contain function")
		}
		minValue := minPercent * functionCall.GetFloat32Field(fieldInfo.Name)
		profile.SortBy(fieldInfo.Name)
		profile = profile.SelectGreater(fieldInfo.Name, minValue)

		fmt.Printf("Parents of %s:\n", function)
		err := renderProfile(parentsProfile, field, fieldInfo, minValue)
		if err != nil {
			return err
		}

		fmt.Printf("Children of %s:\n", function)
		err = renderProfile(childrenProfile, field, fieldInfo, minValue)
		if err != nil {
			return err
		}
	}

	return nil
}
