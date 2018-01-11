package cmd

import (
	"fmt"
	"os"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
)

func renderProfile(profile *xhprof.Profile, field string, fieldInfo FieldInfo, minPercent float32) error {
	profile.SortBy(fieldInfo.Name)
	main, err := profile.GetMain()
	if err != nil {
		return err
	}

	minValue := minPercent * main.GetFloat32Field(fieldInfo.Name)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function", "Count", fieldInfo.Header, fmt.Sprintf("Excl. %s", fieldInfo.Header)})
	for _, call := range profile.Calls {
		if call.GetFloat32Field(fieldInfo.Name) < minValue {
			break
		}

		table.Append(getRowByField(call, field))
	}

	fmt.Printf("Showing XHProf data by %s\n", fieldInfo.Label)
	table.Render()

	return nil
}

func getRowByField(call *xhprof.Call, field string) []string {
	var res []string

	switch field {
	case "wt":
		fallthrough
	case "excl_wt":
		res = []string{
			fmt.Sprintf("%.90s", call.Name),
			fmt.Sprintf("%d", call.Count),
			fmt.Sprintf("%2.2f ms", call.WallTime/1000),
			fmt.Sprintf("%2.2f ms", call.ExclusiveWallTime/1000),
		}
	case "cpu":
		fallthrough
	case "excl_cpu":
		res = []string{
			fmt.Sprintf("%.90s", call.Name),
			fmt.Sprintf("%d", call.Count),
			fmt.Sprintf("%2.2f ms", call.CpuTime/1000),
			fmt.Sprintf("%2.2f ms", call.ExclusiveCpuTime/1000),
		}
	case "io":
		fallthrough
	case "excl_io":
		res = []string{
			fmt.Sprintf("%.90s", call.Name),
			fmt.Sprintf("%d", call.Count),
			fmt.Sprintf("%2.2f ms", call.IoTime/1000),
			fmt.Sprintf("%2.2f ms", call.ExclusiveIoTime/1000),
		}
	case "memory":
		fallthrough
	case "excl_memory":
		res = []string{
			fmt.Sprintf("%.90s", call.Name),
			fmt.Sprintf("%d", call.Count),
			fmt.Sprintf("%2.2f KB", call.Memory/1024),
			fmt.Sprintf("%2.2f KB", call.ExclusiveMemory/1024),
		}
	}

	return res
}
