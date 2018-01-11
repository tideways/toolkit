package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
)

func renderProfile(profile []xhprof.FlatInfo, field string, fieldInfo FieldInfo, minPercent float32) error {
	xhprof.SortBy(profile, fieldInfo.Name)

	mainInfo, err := getMainInfo(profile)
	if err != nil {
		return err
	}

	minValue := minPercent * mainInfo.GetFloat32Field(fieldInfo.Name)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function", "Count", fieldInfo.Header, fmt.Sprintf("Excl. %s", fieldInfo.Header)})

	for _, flatInfo := range profile {
		if flatInfo.GetFloat32Field(fieldInfo.Name) < minValue {
			break
		}

		table.Append(getRowByField(flatInfo, field))
	}

	fmt.Printf("Showing XHProf data by %s\n", fieldInfo.Label)
	table.Render()

	return nil
}

func getMainInfo(profile []xhprof.FlatInfo) (*xhprof.FlatInfo, error) {
	for _, info := range profile {
		if info.Name == "main()" {
			return &info, nil
		}
	}

	return nil, errors.New("Profile has no main()")
}

func getRowByField(flatInfo xhprof.FlatInfo, field string) []string {
	var res []string

	switch field {
	case "wt":
		fallthrough
	case "excl_wt":
		res = []string{
			fmt.Sprintf("%.90s", flatInfo.Name),
			fmt.Sprintf("%d", flatInfo.Calls),
			fmt.Sprintf("%2.2f ms", flatInfo.WallTime/1000),
			fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveWallTime/1000),
		}
	case "cpu":
		fallthrough
	case "excl_cpu":
		res = []string{
			fmt.Sprintf("%.90s", flatInfo.Name),
			fmt.Sprintf("%d", flatInfo.Calls),
			fmt.Sprintf("%2.2f ms", flatInfo.CpuTime/1000),
			fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveCpuTime/1000),
		}
	case "io":
		fallthrough
	case "excl_io":
		res = []string{
			fmt.Sprintf("%.90s", flatInfo.Name),
			fmt.Sprintf("%d", flatInfo.Calls),
			fmt.Sprintf("%2.2f ms", flatInfo.IoTime/1000),
			fmt.Sprintf("%2.2f ms", flatInfo.ExclusiveIoTime/1000),
		}
	case "memory":
		fallthrough
	case "excl_memory":
		res = []string{
			fmt.Sprintf("%.90s", flatInfo.Name),
			fmt.Sprintf("%d", flatInfo.Calls),
			fmt.Sprintf("%2.2f KB", flatInfo.Memory/1024),
			fmt.Sprintf("%2.2f KB", flatInfo.ExclusiveMemory/1024),
		}
	}

	return res
}
