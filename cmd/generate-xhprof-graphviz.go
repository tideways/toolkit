package cmd

import (
	"io/ioutil"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generateXhprofGraphvizCmd)
	generateXhprofGraphvizCmd.Flags().Float32VarP(&threshold, "threshold", "t", 1, "Display items having greater ratio of excl_wt (default 1%) with respect to main()")
	generateXhprofGraphvizCmd.Flags().StringVarP(&function, "function", "f", "", "If provided, the graph will be generated only for functions directly related to this one")
	generateXhprofGraphvizCmd.Flags().BoolVarP(&criticalPath, "critical-path", "", false, "If present, the critical path will be highlighted")
	generateXhprofGraphvizCmd.Flags().StringVarP(&outFile, "out-file", "o", "", "The path to store the resulting graph")
}

var (
	threshold    float32
	criticalPath bool
)

var generateXhprofGraphvizCmd = &cobra.Command{
	Use:   "generate-xhprof-graphviz filepaths...",
	Short: "Parse the output of JSON serialized XHProf outputs into a dot script for graphviz.",
	Long:  `Parse the output of JSON serialized XHProf outputs into a dot script for graphviz.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  generateXhprofGraphviz,
}

func generateXhprofGraphviz(cmd *cobra.Command, args []string) error {
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

	threshold /= 100
	dot, err := xhprof.GenerateDotScript(avgMap, threshold, function, criticalPath, nil, nil)
	if err != nil {
		return err
	}

	if len(outFile) == 0 {
		outFile = "callgraph.dot"
	}

	err = ioutil.WriteFile(outFile, []byte(dot), 0755)
	if err != nil {
		return err
	}

	return nil
}
