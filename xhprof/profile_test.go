package xhprof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvgProfiles(t *testing.T) {
	expected := Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             2,
				WallTime:          300,
				ExclusiveWallTime: 200,
				CpuTime:           150,
				ExclusiveCpuTime:  100,
				IoTime:            150,
				ExclusiveIoTime:   100,
				Memory:            256,
				ExclusiveMemory:   128,
			},
		},
	}

	p1 := &Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             2,
				WallTime:          200,
				ExclusiveWallTime: 200,
				CpuTime:           100,
				ExclusiveCpuTime:  50,
				IoTime:            100,
				ExclusiveIoTime:   150,
				Memory:            256,
				ExclusiveMemory:   128,
			},
		},
	}

	p2 := &Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             3,
				WallTime:          400,
				ExclusiveWallTime: 200,
				CpuTime:           200,
				ExclusiveCpuTime:  150,
				IoTime:            200,
				ExclusiveIoTime:   50,
				Memory:            256,
				ExclusiveMemory:   128,
			},
		},
	}

	p3 := &Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             2,
				WallTime:          300,
				ExclusiveWallTime: 200,
				CpuTime:           150,
				ExclusiveCpuTime:  100,
				IoTime:            150,
				ExclusiveIoTime:   100,
				Memory:            256,
				ExclusiveMemory:   128,
			},
		},
	}

	p := AvgProfiles([]*Profile{p1, p2, p3})

	require.Len(t, p.Calls, 1)
	assert.Equal(t, expected.Calls[0].Count, p.Calls[0].Count)
	assert.Equal(t, expected.Calls[0].WallTime, p.Calls[0].WallTime)
	assert.Equal(t, expected.Calls[0].ExclusiveWallTime, p.Calls[0].ExclusiveWallTime)
	assert.Equal(t, expected.Calls[0].CpuTime, p.Calls[0].CpuTime)
	assert.Equal(t, expected.Calls[0].ExclusiveCpuTime, p.Calls[0].ExclusiveCpuTime)
	assert.Equal(t, expected.Calls[0].IoTime, p.Calls[0].IoTime)
	assert.Equal(t, expected.Calls[0].ExclusiveIoTime, p.Calls[0].ExclusiveIoTime)
	assert.Equal(t, expected.Calls[0].Memory, p.Calls[0].Memory)
	assert.Equal(t, expected.Calls[0].ExclusiveMemory, p.Calls[0].ExclusiveMemory)
}
