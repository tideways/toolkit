package xhprof

import (
	"errors"
	"strings"
)

type PairCall struct {
	Count      int     `json:"ct"`
	WallTime   float32 `json:"wt"`
	CpuTime    float32 `json:"cpu"`
	Memory     float32 `json:"mu"`
	PeakMemory float32 `json:"pmu"`
}

func (p *PairCall) Add(o *PairCall) *PairCall {
	p.Count += o.Count
	p.WallTime += o.WallTime
	p.CpuTime += o.CpuTime
	p.Memory += o.Memory
	p.PeakMemory += o.PeakMemory

	return p
}

func (p *PairCall) Divide(d float32) *PairCall {
	p.Count /= int(d)
	p.WallTime /= d
	p.CpuTime /= d
	p.Memory /= d
	p.PeakMemory /= d

	return p
}

type PairCallMap struct {
	M map[string]*PairCall
}

func NewPairCallMap() *PairCallMap {
	m := new(PairCallMap)
	m.M = make(map[string]*PairCall)

	return m
}

func (m *PairCallMap) Flatten() (*Profile, error) {
	var parent string
	var child string

	symbols := make(map[string]*Call)
	for name, info := range m.M {
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

		call.AddPairCall(info)
		symbols[child] = call

		if len(parent) == 0 {
			continue
		}

		if call, ok = symbols[parent]; !ok {
			call = &Call{Name: parent}
		}

		call.SubtractExcl(info)
		symbols[parent] = call
	}

	profile := new(Profile)
	calls := make([]*Call, 0, len(symbols))
	for _, call := range symbols {
		calls = append(calls, call)
	}
	profile.Calls = calls

	main, ok := symbols["main()"]
	if !ok || main == nil {
		return nil, errors.New("Profile has no main()")
	}
	profile.Main = main

	return profile, nil
}

func AvgPairCallMaps(maps []*PairCallMap) *PairCallMap {
	if len(maps) == 1 {
		return maps[0]
	}

	res := NewPairCallMap()

	for _, m := range maps {
		for k, v := range m.M {
			pairCall, ok := res.M[k]
			if !ok {
				pairCall = new(PairCall)
				*pairCall = *v
				res.M[k] = pairCall
				continue
			}

			pairCall.Add(v)
		}
	}

	num := float32(len(maps))
	for _, v := range res.M {
		v.Divide(num)
	}

	return res
}
