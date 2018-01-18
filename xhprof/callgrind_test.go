package xhprof

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCallgrind(t *testing.T) {
	expected := Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             1,
				WallTime:          820,
				ExclusiveWallTime: 20,
			},
			&Call{
				Name:              "func2",
				Count:             5,
				WallTime:          700,
				ExclusiveWallTime: 700,
			},
			&Call{
				Name:              "func1",
				Count:             1,
				WallTime:          400,
				ExclusiveWallTime: 100,
			},
		},
	}

	f, err := os.Open("testdata/callgrind-simple.out")

	require.Nil(t, err)

	profile, err := ParseCallgrind(f)

	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Len(t, profile.Calls, len(expected.Calls))

	require.NotNil(t, profile.Main)
	assert.EqualValues(t, profile.Main, expected.Calls[0])

	profile.SortBy("WallTime")

	for i, c := range profile.Calls {
		assert.EqualValues(t, expected.Calls[i], c)
	}
}

func TestParseCallgrindRealisticSample(t *testing.T) {
	expected := Profile{
		Calls: []*Call{
			&Call{
				Name:              "main()",
				Count:             1,
				WallTime:          305041,
				ExclusiveWallTime: 54,
			},
			&Call{
				Name:              "require::/var/www/wordpress/wp-blog-header.php",
				Count:             1,
				WallTime:          304980,
				ExclusiveWallTime: 85,
			},
		},
	}

	f, err := os.Open("../tests/data/cachegrind.out")

	require.Nil(t, err)

	profile, err := ParseCallgrind(f)

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
