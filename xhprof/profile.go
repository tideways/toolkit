package xhprof

import (
	"sort"
)

type Profile struct {
	Calls []*Call
	Main  *Call
}

func (p *Profile) GetMain() *Call {
	return p.Main
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

func AvgProfiles(profiles []*Profile) *Profile {
	type CallSum struct {
		Call *Call
		Num  int
	}
	callMap := make(map[string]*CallSum)

	for _, p := range profiles {
		for _, c := range p.Calls {
			call, ok := callMap[c.Name]
			if !ok {
				call = &CallSum{Call: new(Call), Num: 1}
				*call.Call = *c
				callMap[call.Call.Name] = call
				continue
			}

			call.Call.Add(c)
			call.Num += 1
		}
	}

	res := new(Profile)
	calls := make([]*Call, 0, len(callMap))
	for _, call := range callMap {
		avgCall := call.Call.Divide(float32(call.Num))
		if call.Call.Name == "main()" {
			res.Main = avgCall
		}

		calls = append(calls, avgCall)
	}
	res.Calls = calls

	return res
}
