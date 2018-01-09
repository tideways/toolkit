package cmd

import (
	"fmt"
	"os"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
)

func renderProfile(profile []xhprof.FlatInfo, fieldInfo FieldInfo) {
	xhprof.SortBy(profile, fieldInfo.Name)

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
