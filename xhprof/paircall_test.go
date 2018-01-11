package xhprof

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	expected := []struct {
		Name              string
		Count             int
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
			Count:             1,
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
			Count:             2,
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
			Count:             10,
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

	sample := map[string]PairCall{
		"main()": PairCall{
			WallTime: 1000,
			Count:    1,
			CpuTime:  400,
			Memory:   1500,
		},
		"main()==>foo": PairCall{
			WallTime: 500,
			Count:    2,
			CpuTime:  200,
			Memory:   700,
		},
		"foo==>bar": PairCall{
			WallTime: 200,
			Count:    10,
			CpuTime:  100,
			Memory:   300,
		},
	}

	profile, err := Flatten(sample)
	profile.SortBy("WallTime")

	var expectedType *Profile
	require.Nil(t, err)
	require.IsType(t, profile, expectedType)
	require.Len(t, profile.Calls, len(expected))
	assert.Equal(t, float32(1000), profile.Main.WallTime)

	for i, call := range profile.Calls {
		assert.Equal(t, expected[i].Name, call.Name)
		assert.Equal(t, expected[i].Count, call.Count)
		assert.Equal(t, expected[i].WallTime, call.WallTime)
		assert.Equal(t, expected[i].ExclusiveWallTime, call.ExclusiveWallTime)
		assert.Equal(t, expected[i].CpuTime, call.CpuTime)
		assert.Equal(t, expected[i].ExclusiveCpuTime, call.ExclusiveCpuTime)
		assert.Equal(t, expected[i].Memory, call.Memory)
		assert.Equal(t, expected[i].ExclusiveMemory, call.ExclusiveMemory)
		assert.Equal(t, expected[i].IoTime, call.IoTime)
		assert.Equal(t, expected[i].ExclusiveIoTime, call.ExclusiveIoTime)
	}
}
