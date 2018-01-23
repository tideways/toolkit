package xhprof

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCallgrind(t *testing.T) {
	expected := &PairCallMap{
		M: map[string]*PairCall{
			"main()": &PairCall{
				Count:    1,
				WallTime: 820,
			},
			"main()==>func2": &PairCall{
				Count:    3,
				WallTime: 400,
			},
			"main()==>func1": &PairCall{
				Count:    1,
				WallTime: 400,
			},
			"func1==>func2": &PairCall{
				Count:    2,
				WallTime: 300,
			},
		},
	}

	f, err := os.Open("testdata/callgrind-simple.out")
	require.Nil(t, err)

	m, err := ParseCallgrind(f)
	require.Nil(t, err)

	assert.EqualValues(t, expected, m)
}
