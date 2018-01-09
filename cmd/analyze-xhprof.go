package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(xhprofCmd)
	xhprofCmd.Flags().StringVarP(&xhprofDimension, "dimension", "d", "excl_wt", "Dimension to view/sort (wt, excl_wt, cpu, excl_cpu, memory, excl_memory).")
	xhprofCmd.Flags().IntVarP(&xhprofNumItems, "size", "s", 30, "Number of items to list in table")
}

type ByWallTime []xhprof.FlatInfo

func (a ByWallTime) Len() int           { return len(a) }
func (a ByWallTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWallTime) Less(i, j int) bool { return a[i].WallTime > a[j].WallTime }

type ByCpuTime []xhprof.FlatInfo

func (a ByCpuTime) Len() int           { return len(a) }
func (a ByCpuTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCpuTime) Less(i, j int) bool { return a[i].CpuTime > a[j].CpuTime }

type ByIoTime []xhprof.FlatInfo

func (a ByIoTime) Len() int           { return len(a) }
func (a ByIoTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIoTime) Less(i, j int) bool { return a[i].IoTime > a[j].IoTime }

type ByMemory []xhprof.FlatInfo

func (a ByMemory) Len() int           { return len(a) }
func (a ByMemory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMemory) Less(i, j int) bool { return a[i].Memory > a[j].Memory }

type ByExclusiveWallTime []xhprof.FlatInfo

func (a ByExclusiveWallTime) Len() int      { return len(a) }
func (a ByExclusiveWallTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByExclusiveWallTime) Less(i, j int) bool {
	return a[i].ExclusiveWallTime > a[j].ExclusiveWallTime
}

type ByExclusiveCpuTime []xhprof.FlatInfo

func (a ByExclusiveCpuTime) Len() int      { return len(a) }
func (a ByExclusiveCpuTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByExclusiveCpuTime) Less(i, j int) bool {
	return a[i].ExclusiveCpuTime > a[j].ExclusiveCpuTime
}

type ByExclusiveIoTime []xhprof.FlatInfo

func (a ByExclusiveIoTime) Len() int      { return len(a) }
func (a ByExclusiveIoTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByExclusiveIoTime) Less(i, j int) bool {
	return a[i].ExclusiveIoTime > a[j].ExclusiveIoTime
}

type ByExclusiveMemory []xhprof.FlatInfo

func (a ByExclusiveMemory) Len() int      { return len(a) }
func (a ByExclusiveMemory) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByExclusiveMemory) Less(i, j int) bool {
	return a[i].ExclusiveMemory > a[j].ExclusiveMemory
}

var xhprofDimension string
var xhprofNumItems int

var xhprofCmd = &cobra.Command{
	Use:   "analyze-xhprof filepath",
	Short: "Parse the output of JSON serialized XHProf output into a sorted tabular output.",
	Long:  `Parse the output of JSON serialized XHProf output into a sorted tabular output.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		var dimensionLabel string
		var header string
		switch xhprofDimension {
		case "cpu":
			sort.Sort(ByCpuTime(profile))
			dimensionLabel = "Inclusive CPU-Time"
			header = "CPU-Time"
		case "excl_cpu":
			sort.Sort(ByExclusiveCpuTime(profile))
			dimensionLabel = "Exclusive CPU-Time"
			header = "CPU-Time"
		case "io":
			sort.Sort(ByIoTime(profile))
			dimensionLabel = "Inclusive I/O-Time"
			header = "I/O-Time"
		case "excl_io":
			sort.Sort(ByExclusiveIoTime(profile))
			dimensionLabel = "Exclusive I/O-Time"
			header = "I/O-Time"
		case "memory":
			sort.Sort(ByMemory(profile))
			dimensionLabel = "Inclusive Memory"
			header = "Memory"
		case "excl_memory":
			sort.Sort(ByExclusiveMemory(profile))
			dimensionLabel = "Exclusive Memory"
			header = "Memory"
		case "wt":
			sort.Sort(ByWallTime(profile))
			dimensionLabel = "Inclusive Wall-Time"
			header = "Wall-Time"
		case "excl_wt":
			fallthrough
		default:
			sort.Sort(ByExclusiveWallTime(profile))
			dimensionLabel = "Exclusive Wall-Time"
			header = "Wall-Time"
		}

		fmt.Printf("Showing XHProf data by %s (%s)\n", dimensionLabel, xhprofDimension)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Function", "Count", header, fmt.Sprintf("Excl. %s", header)})

		for _, flatInfo := range profile[0:xhprofNumItems] {
			switch xhprofDimension {
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
		table.Render() // Send output

		return nil
	},
}
