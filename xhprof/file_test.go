package xhprof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPairCallMap(t *testing.T) {
	expected := &PairCallMap{
		M: map[string]*PairCall{
			"main()": &PairCall{
				WallTime: 1000,
				Count:    1,
				CpuTime:  400,
				Memory:   1500,
			},
			"main()==>foo": &PairCall{
				WallTime: 500,
				Count:    2,
				CpuTime:  200,
				Memory:   700,
			},
			"foo==>bar": &PairCall{
				WallTime: 200,
				Count:    10,
				CpuTime:  100,
				Memory:   300,
			},
		},
	}

	f := NewFile("testdata/simple.xhprof", "xhprof")
	m, err := f.GetPairCallMap()
	require.Nil(t, err)
	assert.EqualValues(t, expected, m)
}
