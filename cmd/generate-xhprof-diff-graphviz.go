package cmd

import (
	"io/ioutil"

	"github.com/tideways/toolkit/xhprof"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(generateXhprofDiffGraphvizCmd)
	generateXhprofDiffGraphvizCmd.Flags().Float32VarP(&threshold, "threshold", "t", 1, "Display items having greater ratio of excl_wt (default 1%) with respect to main()")
	generateXhprofDiffGraphvizCmd.Flags().StringVarP(&outFile, "out-file", "o", "callgraph.dot", "The path to store the resulting graph")
}

var generateXhprofDiffGraphvizCmd = &cobra.Command{
	Use:   "generate-xhprof-diff-graphviz filepaths...",
	Short: "Parse the output of two JSON serialized XHProf outputs, and generate a dot script out of their diff.",
	Long:  `Parse the output of two JSON serialized XHProf outputs, and generate a dot script out of their diff.`,
	Args:  cobra.ExactArgs(2),
	RunE:  generateXhprofDiffGraphviz,
}

func generateXhprofDiffGraphviz(cmd *cobra.Command, args []string) error {
	f := xhprof.NewFile(args[0], "xhprof")
	m1, err := f.GetPairCallMap()
	if err != nil {
		return err
	}

	f = xhprof.NewFile(args[1], "xhprof")
	m2, err := f.GetPairCallMap()
	if err != nil {
		return err
	}

	threshold /= 100
	dot, err := xhprof.GenerateDiffDotScript(m1, m2, threshold)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outFile, []byte(dot), 0755)
	if err != nil {
		return err
	}

	return nil
}
