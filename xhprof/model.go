package xhprof

import (
	"reflect"
	"sort"
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
