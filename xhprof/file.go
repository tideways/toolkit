package xhprof

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type File struct {
	Path   string
	Format string
}

func NewFile(path, format string) (f *File) {
	f = new(File)
	f.Path = path
	f.Format = format

	return
}

func (f *File) GetProfile() (*Profile, error) {
	if f.Format == "callgrind" {
		fh, err := os.Open(f.Path)
		if err != nil {
			return nil, err
		}

		return ParseCallgrind(fh)
	}

	m, err := f.GetPairCallMap()
	if err != nil {
		return nil, err
	}

	return m.Flatten(), nil
}

func (f *File) GetPairCallMap() (m *PairCallMap, err error) {
	var rawData []byte
	if rawData, err = ioutil.ReadFile(f.Path); err != nil {
		return
	}

	m = new(PairCallMap)
	if err = json.Unmarshal(rawData, &m.M); err != nil {
		return
	}

	return
}

func (f *File) WritePairCallMap(m *PairCallMap) error {
	data, err := json.Marshal(m.M)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.Path, data, 0755)
}
