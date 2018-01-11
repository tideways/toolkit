package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(xhprofCmd)
	xhprofCmd.Flags().StringVarP(&field, "field", "f", "excl_wt", "Field to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory, io, excl_io)")
	xhprofCmd.Flags().Float32VarP(&minPercent, "min", "m", 1, "Display items having minimum percentage (default 1%) of --field, with respect to main()")
}

var (
	field      string
	minPercent float32
)

var xhprofCmd = &cobra.Command{
	Use:   "analyze-xhprof filepath",
	Short: "Parse the output of JSON serialized XHProf output into a sorted tabular output.",
	Long:  `Parse the output of JSON serialized XHProf output into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  analyzeXhprof,
}

func analyzeXhprof(cmd *cobra.Command, args []string) error {
	rawData, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	var data map[string]xhprof.PairCall
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return err
	}

	profile, err := xhprof.Flatten(data)
	if err != nil {
		return err
	}

	fieldInfo, ok := fieldsMap[field]
	if !ok {
		fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
		field = "excl_wt"
		fieldInfo = fieldsMap[field]
	}

	minPercent = minPercent / 100.0
	err = renderProfile(profile, field, fieldInfo, minPercent)
	if err != nil {
		return err
	}

	return nil
}
