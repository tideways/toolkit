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

func Flatten(data map[string]*PairCall) (*Profile, error) {
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
