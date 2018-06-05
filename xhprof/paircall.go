package xhprof

import (
	"strings"
)

type PairCall struct {
	Count       int     `json:"ct"`
	WallTime    float32 `json:"wt"`
	CpuTime     float32 `json:"cpu"`
	Memory      float32 `json:"mu"`
	PeakMemory  float32 `json:"pmu"`
	NumAlloc    float32 `json:"mem.na"`
	NumFree     float32 `json:"mem.nf"`
	AllocAmount float32 `json:"mem.aa"`
}

func (p *PairCall) Add(o *PairCall) *PairCall {
	p.Count += o.Count
	p.WallTime += o.WallTime
	p.CpuTime += o.CpuTime
	p.Memory += o.Memory
	p.PeakMemory += o.PeakMemory
	p.NumAlloc += o.NumAlloc
	p.NumFree += o.NumFree
	p.AllocAmount += o.AllocAmount

	return p
}

func (p *PairCall) Divide(d float32) *PairCall {
	p.Count /= int(d)
	p.WallTime /= d
	p.CpuTime /= d
	p.Memory /= d
	p.PeakMemory /= d
	p.NumAlloc /= d
	p.NumFree /= d
	p.AllocAmount /= d

	return p
}

func (p *PairCall) Subtract(o *PairCall) *PairCall {
	p.Count -= o.Count
	p.WallTime -= o.WallTime
	p.CpuTime -= o.CpuTime
	p.Memory -= o.Memory
	p.PeakMemory -= o.PeakMemory
	p.NumAlloc -= o.NumAlloc
	p.NumFree -= o.NumFree
	p.AllocAmount -= o.AllocAmount

	return p
}

type NearestFamily struct {
	Children      *PairCallMap
	Parents       *PairCallMap
	ChildrenCount int
	ParentsCount  int
}

func NewNearestFamily() *NearestFamily {
	f := new(NearestFamily)
	f.Children = NewPairCallMap()
	f.Parents = NewPairCallMap()

	return f
}

type PairCallMap struct {
	M map[string]*PairCall
}

func NewPairCallMap() *PairCallMap {
	m := new(PairCallMap)
	m.M = make(map[string]*PairCall)

	return m
}

func (m *PairCallMap) NewPairCall(name string) *PairCall {
	pc, ok := m.M[name]
	if ok {
		return pc
	}

	pc = new(PairCall)
	m.M[name] = pc

	return pc
}

func (m *PairCallMap) GetCallMap() map[string]*Call {
	symbols := make(map[string]*Call)
	for name, info := range m.M {
		parent, child := parsePairName(name)

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

	return symbols
}

func (m *PairCallMap) Flatten() *Profile {
	symbols := m.GetCallMap()

	profile := new(Profile)
	calls := make([]*Call, 0, len(symbols))
	for _, call := range symbols {
		calls = append(calls, call)
	}
	profile.Calls = calls

	main, ok := symbols["main()"]
	if ok {
		profile.Main = main
	}

	return profile
}

func (m *PairCallMap) ComputeNearestFamily(f string) *NearestFamily {
	family := NewNearestFamily()

	for name, info := range m.M {
		parent, child := parsePairName(name)
		if parent == f {
			c, ok := family.Children.M[child]
			if !ok {
				c = new(PairCall)
				family.Children.M[child] = c
			}

			c.WallTime += info.WallTime
			c.Count += info.Count
			family.ChildrenCount += info.Count
		}

		if child == f && parent != "" {
			p, ok := family.Parents.M[parent]
			if !ok {
				p = new(PairCall)
				family.Parents.M[parent] = p
			}

			p.WallTime += info.WallTime
			p.Count += info.Count
			family.ParentsCount += info.Count
		}
	}

	return family
}

func (m *PairCallMap) GetChildrenMap() map[string][]string {
	r := make(map[string][]string)

	for name, _ := range m.M {
		parent, child := parsePairName(name)
		if _, ok := r[parent]; !ok {
			r[parent] = make([]string, 0, 1)
		}

		r[parent] = append(r[parent], child)
	}

	return r
}

func (m *PairCallMap) Copy() *PairCallMap {
	r := NewPairCallMap()

	for name, info := range m.M {
		c := new(PairCall)
		*c = *info
		r.M[name] = c
	}

	return r
}

func (m *PairCallMap) Subtract(o *PairCallMap) *PairCallMap {
	r := m.Copy()

	for name, info := range o.M {
		p, ok := r.M[name]
		if !ok {
			p = new(PairCall)
			r.M[name] = p
		}

		p.Subtract(info)
	}

	return r
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

func parsePairName(name string) (parent string, child string) {
	fns := strings.Split(name, "==>")
	if len(fns) == 2 {
		parent = fns[0]
		child = fns[1]
	} else {
		child = fns[0]
	}

	return
}

func pairName(parent, child string) string {
	if parent == "" {
		return child
	} else if child == "" {
		return parent
	}

	return parent + "==>" + child
}
