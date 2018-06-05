package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/tideways/toolkit/xhprof"

	"github.com/olekukonko/tablewriter"
)

type Unit struct {
	Name    string
	Divisor float32
}

var (
	ms    Unit = Unit{Name: "ms", Divisor: 1000.0}
	kb    Unit = Unit{Name: "KB", Divisor: 1024.0}
	plain Unit = Unit{Name: "", Divisor: 1.0}
)

type FieldInfo struct {
	Name   string
	Label  string
	Header string
	Unit   Unit
}

var fieldsMap map[string]FieldInfo = map[string]FieldInfo{
	"wt": FieldInfo{
		Name:   "WallTime",
		Label:  "Inclusive Wall-Time",
		Header: "Wall-Time",
		Unit:   ms,
	},
	"excl_wt": FieldInfo{
		Name:   "ExclusiveWallTime",
		Label:  "Exclusive Wall-Time",
		Header: "Wall-Time",
		Unit:   ms,
	},
	"cpu": FieldInfo{
		Name:   "CpuTime",
		Label:  "Inclusive CPU-Time",
		Header: "CPU-Time",
		Unit:   ms,
	},
	"excl_cpu": FieldInfo{
		Name:   "ExclusiveCpuTime",
		Label:  "Exclusive CPU-Time",
		Header: "CPU-Time",
		Unit:   ms,
	},
	"memory": FieldInfo{
		Name:   "Memory",
		Label:  "Inclusive Memory",
		Header: "Memory",
		Unit:   kb,
	},
	"excl_memory": FieldInfo{
		Name:   "ExclusiveMemory",
		Label:  "Exclusive Memory",
		Header: "Memory",
		Unit:   kb,
	},
	"io": FieldInfo{
		Name:   "IoTime",
		Label:  "Inclusive I/O-Time",
		Header: "I/O-Time",
		Unit:   ms,
	},
	"excl_io": FieldInfo{
		Name:   "ExclusiveIoTime",
		Label:  "Exclusive I/O-Time",
		Header: "I/O-Time",
		Unit:   ms,
	},
	"num_alloc": FieldInfo{
		Name:   "NumAlloc",
		Label:  "Number of Allocations",
		Header: "Num. Alloc.",
		Unit:   plain,
	},
	"alloc_amt": FieldInfo{
		Name:   "AllocAmount",
		Label:  "Amount of allocated Memory",
		Header: "Alloc. Amount",
		Unit:   kb,
	},
	"num_free": FieldInfo{
		Name:   "NumFree",
		Label:  "Number of Frees",
		Header: "Num. Frees",
		Unit:   plain,
	},
}

func renderProfile(profile *xhprof.Profile, field string, fieldInfo FieldInfo, minValue float32) error {
	header := fieldInfo.Header
	exclHeader := "Excl. " + fieldInfo.Header
	var fields []FieldInfo
	var headers []string
	if strings.HasPrefix(field, "excl_") {
		fields = []FieldInfo{fieldsMap[strings.TrimPrefix(field, "excl_")], fieldInfo}
		exclHeader = fmt.Sprintf("%s (>= %2.2f %s)", exclHeader, minValue/fieldInfo.Unit.Divisor, fieldInfo.Unit.Name)
		headers = []string{"Function", "Count", header, exclHeader}
	} else if field == "num_alloc" {
		fields = []FieldInfo{fieldsMap["num_alloc"], fieldsMap["alloc_amt"], fieldsMap["num_free"]}
		header = fmt.Sprintf("%s (>= %2.2f %s)", header, minValue/fieldInfo.Unit.Divisor, fieldInfo.Unit.Name)
		headers = []string{"Function", "Count", fieldsMap["num_alloc"].Header, fieldsMap["alloc_amt"].Header, fieldsMap["num_free"].Header}
	} else {
		fields = []FieldInfo{fieldInfo}
		header = fmt.Sprintf("%s (>= %2.2f %s)", header, minValue/fieldInfo.Unit.Divisor, fieldInfo.Unit.Name)
		headers = []string{"Function", "Count", header}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, call := range profile.Calls {
		table.Append(getRow(call, fields))
	}

	fmt.Printf("Showing XHProf data by %s\n", fieldInfo.Label)
	table.Render()

	return nil
}

func getRow(call *xhprof.Call, fields []FieldInfo) []string {
	res := []string{
		fmt.Sprintf("%.90s", call.Name),
		fmt.Sprintf("%d", call.Count),
	}

	for _, field := range fields {
		col := fmt.Sprintf("%2.2f %s", call.GetFloat32Field(field.Name)/field.Unit.Divisor, field.Unit.Name)
		res = append(res, col)
	}

	return res
}

func renderProfileDiff(diff *xhprof.ProfileDiff, limit int) error {
	diff.Sort()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Function", "Count", "Wall-Time", "Fraction Wall-Time From", "Fraction Wall-Time To"})
	for i, call := range diff.Calls {
		if i >= limit {
			break
		}

		row := []string{
			fmt.Sprintf("%.90s", call.Name),
			fmt.Sprintf("%d", call.Count),
			fmt.Sprintf("%2.2f ms", call.WallTime/1000),
			fmt.Sprintf("%2.2f", call.FractionWtFrom),
			fmt.Sprintf("%2.2f", call.FractionWtTo),
		}

		table.Append(row)
	}

	fmt.Printf("Showing XHProf data by the difference of fractions\n")
	table.Render()

	return nil
}
