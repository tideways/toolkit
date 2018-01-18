package xhprof

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	expected := &Call{
		Count:             7,
		WallTime:          1000,
		ExclusiveWallTime: 600,
		CpuTime:           500,
		ExclusiveCpuTime:  300,
		IoTime:            500,
		ExclusiveIoTime:   300,
		Memory:            1024,
		ExclusiveMemory:   512,
	}

	c1 := &Call{
		Count:             2,
		WallTime:          300,
		ExclusiveWallTime: 200,
		CpuTime:           100,
		ExclusiveCpuTime:  50,
		IoTime:            200,
		ExclusiveIoTime:   150,
		Memory:            256,
		ExclusiveMemory:   128,
	}

	c2 := &Call{
		Count:             3,
		WallTime:          400,
		ExclusiveWallTime: 200,
		CpuTime:           200,
		ExclusiveCpuTime:  150,
		IoTime:            200,
		ExclusiveIoTime:   50,
		Memory:            256,
		ExclusiveMemory:   128,
	}

	c3 := &Call{
		Count:             2,
		WallTime:          300,
		ExclusiveWallTime: 200,
		CpuTime:           200,
		ExclusiveCpuTime:  100,
		IoTime:            100,
		ExclusiveIoTime:   100,
		Memory:            512,
		ExclusiveMemory:   256,
	}

	c1.Add(c2).Add(c3)

	assert.EqualValues(t, expected, c1)
}

func TestAddPairCall(t *testing.T) {
	expected := &Call{
		Count:             5,
		WallTime:          700,
		ExclusiveWallTime: 600,
		CpuTime:           300,
		ExclusiveCpuTime:  250,
		IoTime:            400,
		ExclusiveIoTime:   350,
		Memory:            512,
		PeakMemory:        300,
		ExclusiveMemory:   384,
	}

	c := &Call{
		Count:             2,
		WallTime:          300,
		ExclusiveWallTime: 200,
		CpuTime:           100,
		ExclusiveCpuTime:  50,
		IoTime:            200,
		ExclusiveIoTime:   150,
		Memory:            256,
		ExclusiveMemory:   128,
	}

	p := &PairCall{
		Count:      3,
		WallTime:   400,
		CpuTime:    200,
		Memory:     256,
		PeakMemory: 300,
	}

	c.AddPairCall(p)

	assert.EqualValues(t, expected, c)
}

func TestSubtractExcl(t *testing.T) {
	expected := &Call{
		Count:             4,
		WallTime:          500,
		ExclusiveWallTime: 200,
		CpuTime:           200,
		ExclusiveCpuTime:  100,
		IoTime:            300,
		ExclusiveIoTime:   100,
		Memory:            512,
		ExclusiveMemory:   128,
	}

	c := &Call{
		Count:             4,
		WallTime:          500,
		ExclusiveWallTime: 300,
		CpuTime:           200,
		ExclusiveCpuTime:  150,
		IoTime:            300,
		ExclusiveIoTime:   150,
		Memory:            512,
		ExclusiveMemory:   256,
	}

	p := &PairCall{
		Count:    1,
		WallTime: 100,
		CpuTime:  50,
		Memory:   128,
	}

	c.SubtractExcl(p)

	assert.EqualValues(t, expected, c)
}

func TestDivide(t *testing.T) {
	expected := &Call{
		Count:             3,
		WallTime:          900,
		ExclusiveWallTime: 690,
		CpuTime:           400,
		ExclusiveCpuTime:  300,
		IoTime:            500,
		ExclusiveIoTime:   390,
		Memory:            1024,
		ExclusiveMemory:   512,
	}

	c1 := &Call{
		Count:             10,
		WallTime:          2700,
		ExclusiveWallTime: 2070,
		CpuTime:           1200,
		ExclusiveCpuTime:  900,
		IoTime:            1500,
		ExclusiveIoTime:   1170,
		Memory:            3072,
		ExclusiveMemory:   1536,
	}

	c1.Divide(3)

	assert.Equal(t, expected, c1)
}
