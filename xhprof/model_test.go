package xhprof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	expected := []struct {
		Name              string
		Calls             int
		WallTime          float32
		ExclusiveWallTime float32
		CpuTime           float32
		ExclusiveCpuTime  float32
		Memory            float32
		ExclusiveMemory   float32
		IoTime            float32
		ExclusiveIoTime   float32
	}{
		{
			Name:              "main()",
			Calls:             1,
			WallTime:          1000,
			ExclusiveWallTime: 500,
			CpuTime:           400,
			ExclusiveCpuTime:  200,
			Memory:            1500,
			ExclusiveMemory:   800,
			IoTime:            600,
			ExclusiveIoTime:   300,
		},
		{
			Name:              "foo",
			Calls:             2,
			WallTime:          500,
			ExclusiveWallTime: 300,
			CpuTime:           200,
			ExclusiveCpuTime:  100,
			Memory:            700,
			ExclusiveMemory:   400,
			IoTime:            300,
			ExclusiveIoTime:   200,
		},
		{
			Name:              "bar",
			Calls:             10,
			WallTime:          200,
			ExclusiveWallTime: 200,
			CpuTime:           100,
			ExclusiveCpuTime:  100,
			Memory:            300,
			ExclusiveMemory:   300,
			IoTime:            100,
			ExclusiveIoTime:   100,
		},
	}

	sample := map[string]Info{
		"main()": Info{
			WallTime: 1000,
			Calls:    1,
			CpuTime:  400,
			Memory:   1500,
		},
		"main()==>foo": Info{
			WallTime: 500,
			Calls:    2,
			CpuTime:  200,
			Memory:   700,
		},
		"foo==>bar": Info{
			WallTime: 200,
			Calls:    10,
			CpuTime:  100,
			Memory:   300,
		},
	}

	profile := Flatten(sample)

	var expectedType []FlatInfo
	require.IsType(t, profile, expectedType)
	require.Len(t, profile, len(expected))

	for i, info := range profile {
		assert.Equal(t, expected[i].Name, info.Name)
		assert.Equal(t, expected[i].Calls, info.Calls)
		assert.Equal(t, expected[i].WallTime, info.WallTime)
		assert.Equal(t, expected[i].ExclusiveWallTime, info.ExclusiveWallTime)
		assert.Equal(t, expected[i].CpuTime, info.CpuTime)
		assert.Equal(t, expected[i].ExclusiveCpuTime, info.ExclusiveCpuTime)
		assert.Equal(t, expected[i].Memory, info.Memory)
		assert.Equal(t, expected[i].ExclusiveMemory, info.ExclusiveMemory)
		assert.Equal(t, expected[i].IoTime, info.IoTime)
		assert.Equal(t, expected[i].ExclusiveIoTime, info.ExclusiveIoTime)
	}
}
