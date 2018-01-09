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
	}{
		{
			Name:              "main()",
			Calls:             1,
			WallTime:          1000,
			ExclusiveWallTime: 500,
			CpuTime:           900,
			ExclusiveCpuTime:  700,
		},
		{
			Name:              "foo",
			Calls:             2,
			WallTime:          500,
			ExclusiveWallTime: 300,
			CpuTime:           200,
			ExclusiveCpuTime:  100,
		},
		{
			Name:              "bar",
			Calls:             10,
			WallTime:          200,
			ExclusiveWallTime: 200,
			CpuTime:           100,
			ExclusiveCpuTime:  100,
		},
	}

	sample := map[string]Info{
		"main()": Info{
			WallTime: 1000,
			Calls:    1,
			CpuTime:  900,
		},
		"main()==>foo": Info{
			WallTime: 500,
			Calls:    2,
			CpuTime:  200,
		},
		"foo==>bar": Info{
			WallTime: 200,
			Calls:    10,
			CpuTime:  100,
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
	}
}
