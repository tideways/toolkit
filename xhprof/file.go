package xhprof

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ParseFile(path string, callgrind bool) (profile *Profile, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}

	if callgrind {
		profile, err = ParseCallgrind(f)
		return
	}

	var rawData []byte
	if rawData, err = ioutil.ReadFile(path); err != nil {
		return
	}

	var data map[string]*PairCall
	if err = json.Unmarshal(rawData, &data); err != nil {
		return
	}

	profile, err = Flatten(data)

	return
}
