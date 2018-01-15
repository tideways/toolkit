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

func (p *Profile) Subtract(o *Profile) *ProfileDiff {
	d := new(ProfileDiff)
	diff := make(map[string]*CallDiff)

	oCalls := make(map[string]*Call)
	for _, c := range o.Calls {
		oCalls[c.Name] = c
	}

	for _, c := range p.Calls {
		if c == p.Main {
			continue
		}

		oCall, ok := oCalls[c.Name]
		if !ok {
			callDiff := &CallDiff{
				Name:           c.Name,
				WallTime:       c.WallTime,
				CpuTime:        c.CpuTime,
				FractionWtFrom: c.WallTime / p.Main.WallTime,
				FractionWtTo:   0,
			}
			diff[c.Name] = callDiff
			continue
		}

		var wtChange float32
		var cpuChange float32
		if c.WallTime != oCall.WallTime {
			wtChange = oCall.WallTime - c.WallTime
		}
		if c.CpuTime != oCall.CpuTime {
			cpuChange = oCall.CpuTime - c.CpuTime
		}

		if wtChange != 0 || cpuChange != 0 {
			callDiff := &CallDiff{
				Name:           c.Name,
				WallTime:       wtChange,
				CpuTime:        cpuChange,
				FractionWtFrom: c.WallTime / p.Main.WallTime,
				FractionWtTo:   oCall.WallTime / o.Main.WallTime,
			}
			diff[c.Name] = callDiff
		}

		delete(oCalls, c.Name)
	}

	for _, c := range oCalls {
		diff[c.Name] = &CallDiff{
			Name:           c.Name,
			WallTime:       c.WallTime,
			CpuTime:        c.WallTime,
			FractionWtFrom: 0,
			FractionWtTo:   c.WallTime / o.Main.WallTime,
		}
	}

	d.Calls = make([]*CallDiff, 0, len(diff))
	for _, c := range diff {
		d.Calls = append(d.Calls, c)
	}

	return d
}

type ProfileDiff struct {
	Calls []*CallDiff
}

type ProfileDiffRelative ProfileDiff

func (d ProfileDiffRelative) Len() int { return len(d.Calls) }
func (d ProfileDiffRelative) Swap(i, j int) {
	d.Calls[i], d.Calls[j] = d.Calls[j], d.Calls[i]
}
func (d ProfileDiffRelative) Less(i, j int) bool {
	iFractionDiff := d.Calls[i].FractionWtFrom - d.Calls[i].FractionWtTo
	jFractionDiff := d.Calls[j].FractionWtFrom - d.Calls[j].FractionWtTo

	return iFractionDiff > jFractionDiff
}

func (d *ProfileDiff) Sort() {
	params := ProfileDiffRelative(*d)
	sort.Sort(params)
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
