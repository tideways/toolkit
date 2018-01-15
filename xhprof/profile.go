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
	callMap := make(map[string]*Call)
	for _, p := range profiles {
		for _, c := range p.Calls {
			call, ok := callMap[c.Name]
			if !ok {
				call = new(Call)
				*call = *c
				callMap[call.Name] = call
				continue
			}

			call.Add(c)
		}
	}

	num := float32(len(profiles))
	res := new(Profile)
	calls := make([]*Call, 0, len(callMap))
	for _, call := range callMap {
		avgCall := call.Divide(num)
		if call.Name == "main()" {
			res.Main = avgCall
		}

		calls = append(calls, avgCall)
	}
	res.Calls = calls

	return res
}
