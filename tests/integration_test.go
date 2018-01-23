package main

import (
	"testing"

	"github.com/tideways/toolkit/xhprof"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWPIndexXhprof(t *testing.T) {
	f := xhprof.NewFile("data/wp-index.xhprof", "xhprof")
	m, err := f.GetPairCallMap()
	require.Nil(t, err)
	require.IsType(t, m, new(xhprof.PairCallMap))
	require.NotEmpty(t, m.M)

	assert.Equal(t, m.M["main()"].WallTime, float32(60572))
	assert.Equal(t, m.M["main()"].Count, 1)
	assert.Equal(t, m.M["main()"].CpuTime, float32(54683))
	assert.Equal(t, m.M["main()"].Memory, float32(2738112))
	assert.Equal(t, m.M["main()"].PeakMemory, float32(2596544))
	assert.Equal(t, m.M["wp_set_current_user==>setup_userdata"].WallTime, float32(74))
	assert.Equal(t, m.M["wp_set_current_user==>setup_userdata"].Count, 1)
	assert.Equal(t, m.M["wp_set_current_user==>setup_userdata"].CpuTime, float32(74))
	assert.Equal(t, m.M["wp_set_current_user==>setup_userdata"].Memory, float32(4408))
	assert.Equal(t, m.M["wp_set_current_user==>setup_userdata"].PeakMemory, float32(328))

	p := m.Flatten()
	require.IsType(t, p, new(xhprof.Profile))
	require.NotEmpty(t, p.Calls)

	assert.Equal(t, p.GetMain().WallTime, float32(60572))

	c := p.GetCall("is_search")
	require.IsType(t, c, new(xhprof.Call))

	assert.Equal(t, c.Count, 5)
	assert.Equal(t, c.WallTime, float32(5))
	assert.Equal(t, c.ExclusiveWallTime, float32(4))
	assert.Equal(t, c.CpuTime, float32(5))
	assert.Equal(t, c.ExclusiveCpuTime, float32(3))
	assert.Equal(t, c.IoTime, float32(1))
	assert.Equal(t, c.ExclusiveIoTime, float32(1))
	assert.Equal(t, c.Memory, float32(672))
	assert.Equal(t, c.ExclusiveMemory, float32(560))

	c = p.GetCall("vsprintf")
	require.IsType(t, c, new(xhprof.Call))

	assert.Equal(t, c.Count, 14)
	assert.Equal(t, c.WallTime, float32(18))
	assert.Equal(t, c.ExclusiveWallTime, float32(18))
	assert.Equal(t, c.CpuTime, float32(17))
	assert.Equal(t, c.ExclusiveCpuTime, float32(17))
	assert.Equal(t, c.IoTime, float32(2))
	assert.Equal(t, c.ExclusiveIoTime, float32(2))
	assert.Equal(t, c.Memory, float32(4704))
	assert.Equal(t, c.ExclusiveMemory, float32(4704))
}

func TestParseWPIndexCallgrind(t *testing.T) {
	expected := xhprof.Profile{
		Calls: []*xhprof.Call{
			&xhprof.Call{
				Name:              "main()",
				Count:             1,
				WallTime:          305041,
				ExclusiveWallTime: 54,
			},
			&xhprof.Call{
				Name:              "require::/var/www/wordpress/wp-blog-header.php",
				Count:             1,
				WallTime:          304980,
				ExclusiveWallTime: 85,
			},
		},
	}

	f := xhprof.NewFile("data/cachegrind.out", "callgrind")
	profile, err := f.GetProfile()
	require.Nil(t, err)
	require.NotNil(t, profile)

	require.NotNil(t, profile.Main)
	assert.EqualValues(t, profile.Main, expected.Calls[0])

	for _, c := range profile.Calls {
		if c.Name == expected.Calls[1].Name {
			assert.EqualValues(t, expected.Calls[1], c)
		}
	}
}

func TestComputeNearestFamilyWPIndexXhprof(t *testing.T) {
	expected := &xhprof.NearestFamily{
		Children: &xhprof.PairCallMap{
			M: map[string]*xhprof.PairCall{},
		},
		Parents: &xhprof.PairCallMap{
			M: map[string]*xhprof.PairCall{
				"wpdb::prepare": &xhprof.PairCall{
					WallTime: float32(17),
					Count:    11,
				},
				"get_custom_header": &xhprof.PairCall{
					WallTime: float32(1),
					Count:    3,
				},
			},
		},
		ChildrenCount: 0,
		ParentsCount:  14,
	}

	f := xhprof.NewFile("data/wp-index.xhprof", "xhprof")
	m, err := f.GetPairCallMap()
	require.Nil(t, err)
	require.IsType(t, m, new(xhprof.PairCallMap))
	require.NotEmpty(t, m.M)

	family := m.ComputeNearestFamily("vsprintf")
	require.IsType(t, family, new(xhprof.NearestFamily))
	assert.EqualValues(t, expected, family)
}
