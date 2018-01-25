package xhprof

import (
	"reflect"
)

type Call struct {
	Name              string
	Count             int
	WallTime          float32
	CpuTime           float32
	IoTime            float32
	Memory            float32
	PeakMemory        float32
	ExclusiveWallTime float32
	ExclusiveCpuTime  float32
	ExclusiveMemory   float32
	ExclusiveIoTime   float32

	graphvizId int
}

func (c *Call) GetFloat32Field(field string) float32 {
	cVal := reflect.Indirect(reflect.ValueOf(c))
	return float32(cVal.FieldByName(field).Float())
}

func (c *Call) Add(o *Call) *Call {
	c.Count += o.Count
	c.WallTime += o.WallTime
	c.ExclusiveWallTime += o.ExclusiveWallTime
	c.CpuTime += o.CpuTime
	c.ExclusiveCpuTime += o.ExclusiveCpuTime
	c.Memory += o.Memory
	c.PeakMemory += o.PeakMemory
	c.ExclusiveMemory += o.ExclusiveMemory
	c.IoTime += o.IoTime
	c.ExclusiveIoTime += o.ExclusiveIoTime

	return c
}

func (c *Call) AddPairCall(p *PairCall) *Call {
	c.Count += p.Count
	c.WallTime += p.WallTime
	c.ExclusiveWallTime += p.WallTime
	c.CpuTime += p.CpuTime
	c.ExclusiveCpuTime += p.CpuTime

	io := p.WallTime - p.CpuTime
	if io < 0 {
		io = 0
	}

	c.IoTime += io
	c.ExclusiveIoTime += io

	c.Memory += p.Memory
	c.PeakMemory += p.PeakMemory
	c.ExclusiveMemory += p.Memory

	return c
}

func (c *Call) SubtractExcl(p *PairCall) *Call {
	c.ExclusiveWallTime -= p.WallTime
	c.ExclusiveCpuTime -= p.CpuTime
	c.ExclusiveMemory -= p.Memory

	io := p.WallTime - p.CpuTime
	if io < 0 {
		io = 0
	}

	c.ExclusiveIoTime -= io

	return c
}

func (c *Call) Divide(d float32) *Call {
	c.Count /= int(d)
	c.WallTime /= d
	c.ExclusiveWallTime /= d
	c.CpuTime /= d
	c.ExclusiveCpuTime /= d
	c.Memory /= d
	c.PeakMemory /= d
	c.ExclusiveMemory /= d
	c.IoTime /= d
	c.ExclusiveIoTime /= d

	return c
}

type CallDiff struct {
	Name           string
	WallTime       float32
	Count          int
	FractionWtFrom float32
	FractionWtTo   float32
}
