package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(xhprofCmd)
	xhprofCmd.Flags().StringVarP(&field, "field", "f", "excl_wt", "Field to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory, io, excl_io).")
	xhprofCmd.Flags().IntVarP(&numItems, "size", "s", 30, "Number of items to list in table")
}

var field string
var numItems int

var xhprofCmd = &cobra.Command{
	Use:   "analyze-xhprof filepath",
	Short: "Parse the output of JSON serialized XHProf output into a sorted tabular output.",
	Long:  `Parse the output of JSON serialized XHProf output into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  runXhprof,
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

func runXhprof(cmd *cobra.Command, args []string) error {
	var xhprofData map[string]xhprof.Info
	var symbols map[string]xhprof.FlatInfo
	var child string
	var parent string
	data, err := ioutil.ReadFile(args[0])

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &xhprofData)

	if err != nil {
		return err
	}

	symbols = make(map[string]xhprof.FlatInfo)

	for call, info := range xhprofData {
		var flatInfo xhprof.FlatInfo
		var ok bool
		fns := strings.Split(call, "==>")

		if len(fns) == 2 {
			parent = fns[0]
			child = fns[1]
		} else {
			parent = ""
			child = fns[0]
		}

		if flatInfo, ok = symbols[child]; !ok {
			flatInfo = xhprof.FlatInfo{Name: child}
		}

		flatInfo.Calls += info.Calls

		flatInfo.WallTime += info.WallTime
		flatInfo.ExclusiveWallTime += info.WallTime

		flatInfo.CpuTime += info.CpuTime
		flatInfo.ExclusiveCpuTime += info.CpuTime

		flatInfo.IoTime += (info.WallTime - info.CpuTime)
		flatInfo.ExclusiveIoTime += (info.WallTime - info.CpuTime)

		flatInfo.Memory += info.Memory
		flatInfo.PeakMemory += info.PeakMemory
		flatInfo.ExclusiveMemory += info.Memory

		symbols[child] = flatInfo

		if len(parent) == 0 {
			continue
		}

		if flatInfo, ok = symbols[parent]; !ok {
			flatInfo = xhprof.FlatInfo{Name: parent}
		}

		flatInfo.ExclusiveWallTime -= info.WallTime
		flatInfo.ExclusiveCpuTime -= info.CpuTime
		flatInfo.ExclusiveMemory -= info.Memory
		flatInfo.ExclusiveIoTime -= (info.WallTime - info.CpuTime)

		symbols[parent] = flatInfo
	}

	profile := make([]xhprof.FlatInfo, len(symbols))

	for _, flatInfo := range symbols {
		profile = append(profile, flatInfo)
	}

	fieldInfo, ok := fieldsMap[field]
	if !ok {
		fmt.Printf("Provided field (%s) is not valid, defaulting to excl_wt\n", field)
		field = "excl_wt"
		fieldInfo = fieldsMap[field]
	}

	xhprof.SortBy(profile, fieldInfo.Name)

	renderProfile(profile, fieldInfo)

	return nil
}

func renderProfile(profile []xhprof.FlatInfo, fieldInfo FieldInfo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function", "Count", fieldInfo.Header, fmt.Sprintf("Excl. %s", fieldInfo.Header)})

	for _, flatInfo := range profile[0:numItems] {
		switch field {
		case "wt":
			fallthrough
		case "excl_wt":
			table.Append([]string{
				fmt.Sprintf("%.90s", flatInfo.Name),
				fmt.Sprintf("%d", flatInfo.Calls),
				fmt.Sprintf("%2.2f ms", flatInfo.WallTime/1000),
				fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveWallTime/1000),
			})
		case "cpu":
			fallthrough
		case "excl_cpu":
			table.Append([]string{
				fmt.Sprintf("%.90s", flatInfo.Name),
				fmt.Sprintf("%d", flatInfo.Calls),
				fmt.Sprintf("%2.2f ms", flatInfo.CpuTime/1000),
				fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveCpuTime/1000),
			})
		case "io":
			fallthrough
		case "excl_io":
			table.Append([]string{
				fmt.Sprintf("%.90s", flatInfo.Name),
				fmt.Sprintf("%d", flatInfo.Calls),
				fmt.Sprintf("%2.2f ms", flatInfo.IoTime/1000),
				fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveIoTime/1000),
			})
		case "memory":
			fallthrough
		case "excl_memory":
			table.Append([]string{
				fmt.Sprintf("%.90s", flatInfo.Name),
				fmt.Sprintf("%d", flatInfo.Calls),
				fmt.Sprintf("%2.2f KB", flatInfo.Memory/1024),
				fmt.Sprintf("%2.2f KB", flatInfo.ExclusiveMemory/1024),
			})
		}
	}

	fmt.Printf("Showing XHProf data by %s\n", fieldInfo.Label)
	table.Render()
}
