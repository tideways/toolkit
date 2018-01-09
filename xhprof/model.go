package xhprof

import (
	"reflect"
	"sort"
	"strings"
)

type Info struct {
	Calls      int     `json:"ct"`
	WallTime   float32 `json:"wt"`
	CpuTime    float32 `json:"cpu"`
	Memory     float32 `json:"mu"`
	PeakMemory float32 `json:"pmu"`
}

type FlatInfo struct {
	Name              string
	Calls             int
	WallTime          float32
	CpuTime           float32
	IoTime            float32
	Memory            float32
	PeakMemory        float32
	ExclusiveWallTime float32
	ExclusiveCpuTime  float32
	ExclusiveMemory   float32
	ExclusiveIoTime   float32
}

type SortParams struct {
	Items []FlatInfo
	Field string
}

type ByField SortParams

func (a ByField) Len() int      { return len(a.Items) }
func (a ByField) Swap(i, j int) { a.Items[i], a.Items[j] = a.Items[j], a.Items[i] }
func (a ByField) Less(i, j int) bool {
	iVal := reflect.Indirect(reflect.ValueOf(a.Items[i]))
	jVal := reflect.Indirect(reflect.ValueOf(a.Items[j]))
	return iVal.FieldByName(a.Field).Float() > jVal.FieldByName(a.Field).Float()
}

func SortBy(items []FlatInfo, field string) error {
	params := SortParams{Items: items, Field: field}
	sort.Sort(ByField(params))
	return nil
}

func Flatten(data map[string]Info) []FlatInfo {
	var parent string
	var child string

	symbols := make(map[string]FlatInfo)
	for call, info := range data {
		var flatInfo FlatInfo
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
			flatInfo = FlatInfo{Name: child}
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
			flatInfo = FlatInfo{Name: parent}
		}

		flatInfo.ExclusiveWallTime -= info.WallTime
		flatInfo.ExclusiveCpuTime -= info.CpuTime
		flatInfo.ExclusiveMemory -= info.Memory
		flatInfo.ExclusiveIoTime -= (info.WallTime - info.CpuTime)

		symbols[parent] = flatInfo
	}

	profile := make([]FlatInfo, 0, len(symbols))
	for _, flatInfo := range symbols {
		profile = append(profile, flatInfo)
	}

	return profile
}
