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
	xhprofCmd.Flags().StringVarP(&field, "field", "f", "excl_wt", "Field to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory, io, excl_io).")
	xhprofCmd.Flags().IntVarP(&numItems, "size", "s", 30, "Number of items to list in table")
}

var (
	field    string
	numItems int
)

var xhprofCmd = &cobra.Command{
	Use:   "analyze-xhprof filepath",
	Short: "Parse the output of JSON serialized XHProf output into a sorted tabular output.",
	Long:  `Parse the output of JSON serialized XHProf output into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  analyzeXhprof,
}

type FieldInfo struct {
	Name   string
	Label  string
	Header string
}

var fieldsMap map[string]FieldInfo = map[string]FieldInfo{
	"wt": FieldInfo{
		Name:   "WallTime",
		Label:  "Inclusive Wall-Time",
		Header: "Wall-Time",
	},
	"excl_wt": FieldInfo{
		Name:   "ExclusiveWallTime",
		Label:  "Exclusive Wall-Time",
		Header: "Wall-Time",
	},
	"cpu": FieldInfo{
		Name:   "CpuTime",
		Label:  "Inclusive CPU-Time",
		Header: "CPU-Time",
	},
	"excl_cpu": FieldInfo{
		Name:   "ExclusiveCpuTime",
		Label:  "Exclusive CPU-Time",
		Header: "CPU-Time",
	},
	"memory": FieldInfo{
		Name:   "Memory",
		Label:  "Inclusive Memory",
		Header: "Memory",
	},
	"excl_memory": FieldInfo{
		Name:   "ExclusiveMemory",
		Label:  "Exclusive Memory",
		Header: "Memory",
	},
	"io": FieldInfo{
		Name:   "IoTime",
		Label:  "Inclusive I/O-Time",
		Header: "I/O-Time",
	},
	"excl_io": FieldInfo{
		Name:   "ExclusiveIoTime",
		Label:  "Exclusive I/O-Time",
		Header: "I/O-Time",
	},
}

func analyzeXhprof(cmd *cobra.Command, args []string) error {
	rawData, err := ioutil.ReadFile(args[0])
	if err != nil {
		return err
	}

	var data map[string]xhprof.Info
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return err
	}

	profile := xhprof.Flatten(data)

	fieldInfo, ok := fieldsMap[field]
	if !ok {
		fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
		field = "excl_wt"
		fieldInfo = fieldsMap[field]
	}

	renderProfile(profile, fieldInfo)

	return nil
}
