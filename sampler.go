package optcgo

import "math"

type Sampler interface {
	Sample(id uint64, ratio float64) bool
	Ratio() float64
}

const knuthFactor = uint64(1111111111111111111)

type CommonSampler float64

func (cs CommonSampler) Sample(id uint64, ratio float64) bool {
	if ratio < 1 {
		return id*knuthFactor < uint64(math.MaxUint64*ratio)
	}

	return true
}

func (cs CommonSampler) Ratio() float64 {
	return float64(cs)
}
