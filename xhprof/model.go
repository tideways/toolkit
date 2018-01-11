package xhprof

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

type Call struct {
	Name              string
	Count             int
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

func (i *Call) GetFloat32Field(field string) float32 {
	iVal := reflect.Indirect(reflect.ValueOf(i))
	return float32(iVal.FieldByName(field).Float())
}

type Profile struct {
	Calls []*Call
	Main  *Call
}

func (p *Profile) GetMain() (*Call, error) {
	if p.Main != nil {
		return p.Main, nil
	}

	for _, call := range p.Calls {
		if call.Name == "main()" {
			p.Main = call
			return call, nil
		}
	}

	return nil, errors.New("Profile has no main()")
}

type ProfileByField struct {
	Profile *Profile
	Field   string
}

func (p ProfileByField) Len() int { return len(p.Profile.Calls) }
func (p ProfileByField) Swap(i, j int) {
	p.Profile.Calls[i], p.Profile.Calls[j] = p.Profile.Calls[j], p.Profile.Calls[i]
}
func (p ProfileByField) Less(i, j int) bool {
	return p.Profile.Calls[i].GetFloat32Field(p.Field) > p.Profile.Calls[j].GetFloat32Field(p.Field)
}

func (p *Profile) SortBy(field string) error {
	params := ProfileByField{Profile: p, Field: field}
	sort.Sort(params)
	return nil
}

type PairCall struct {
	Count      int     `json:"ct"`
	WallTime   float32 `json:"wt"`
	CpuTime    float32 `json:"cpu"`
	Memory     float32 `json:"mu"`
	PeakMemory float32 `json:"pmu"`
}

func Flatten(data map[string]PairCall) *Profile {
	var parent string
	var child string

	symbols := make(map[string]*Call)
	for name, info := range data {
		fns := strings.Split(name, "==>")
		if len(fns) == 2 {
			parent = fns[0]
			child = fns[1]
		} else {
			parent = ""
			child = fns[0]
		}

		call, ok := symbols[child]
		if !ok {
			call = &Call{Name: child}
		}

		call.Count += info.Count

		call.WallTime += info.WallTime
		call.ExclusiveWallTime += info.WallTime

		call.CpuTime += info.CpuTime
		call.ExclusiveCpuTime += info.CpuTime

		call.IoTime += (info.WallTime - info.CpuTime)
		call.ExclusiveIoTime += (info.WallTime - info.CpuTime)

		call.Memory += info.Memory
		call.PeakMemory += info.PeakMemory
		call.ExclusiveMemory += info.Memory

		symbols[child] = call

		if len(parent) == 0 {
			continue
		}

		if call, ok = symbols[parent]; !ok {
			call = &Call{Name: parent}
		}

		call.ExclusiveWallTime -= info.WallTime
		call.ExclusiveCpuTime -= info.CpuTime
		call.ExclusiveMemory -= info.Memory
		call.ExclusiveIoTime -= (info.WallTime - info.CpuTime)

		symbols[parent] = call
	}

	profile := new(Profile)
	calls := make([]*Call, 0, len(symbols))
	for _, call := range symbols {
		calls = append(calls, call)
	}
	profile.Calls = calls

	return profile
}
