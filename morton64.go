package morton

import "fmt"

type Morton64 struct {
	dimensions uint64
	bits       uint64
	masks      []uint64
	lshifts    []uint64
	rshifts    []uint64
}

func Make64(dimensions uint64, bits uint64) *Morton64 {
	if dimensions == 0 || bits == 0 || dimensions*bits > 64 {
		panic(fmt.Sprintf("can't make morton64 with %d dimensions and %d bits", dimensions, bits))
	}

	mask := uint64((1 << bits) - 1)

	shift := dimensions * (bits - 1)
	shift |= shift >> 1
	shift |= shift >> 2
	shift |= shift >> 4
	shift |= shift >> 8
	shift |= shift >> 16
	shift |= shift >> 32
	shift -= shift >> 1

	masks := make([]uint64, 0)
	lshifts := make([]uint64, 0)

	masks = append(masks, mask)
	lshifts = append(lshifts, 0)

	for shift > 0 {
		mask = 0
		shifted := uint64(0)

		for bit := uint64(0); bit < bits; bit++ {
			distance := (dimensions * bit) - bit
			shifted |= shift & distance
			mask |= 1 << bit << (((shift - 1) ^ uint64(0xffffffffffffffff)) & distance)
		}

		if shifted != 0 {
			masks = append(masks, mask)
			lshifts = append(lshifts, shift)
		}

		shift >>= 1
	}

	rshifts := make([]uint64, len(lshifts))
	for i := 0; i < len(lshifts)-1; i++ {
		rshifts[i] = lshifts[i+1]
	}
	rshifts[len(rshifts)-1] = 0

	return &Morton64{dimensions: dimensions, bits: bits, masks: masks, lshifts: lshifts, rshifts: rshifts}
}

func (morton *Morton64) Pack(values ...uint64) int64 {
	dimensions := uint64(len(values))
	morton.dimensionsCheck(dimensions)
	for i := uint64(0); i < dimensions; i++ {
		morton.valueCheck(values[i])
	}

	code := uint64(0)
	for i := uint64(0); i < dimensions; i++ {
		code |= morton.split(values[i]) << i
	}

	return int64(code)
}

func (morton *Morton64) SPack(values ...int64) int64 {
	uvalues := make([]uint64, len(values))
	for i := 0; i < len(values); i++ {
		uvalues[i] = morton.shiftSign(values[i])
	}

	return morton.Pack(uvalues...)
}

func (morton *Morton64) Unpack(code int64) []uint64 {
	dimensions := morton.dimensions

	values := make([]uint64, dimensions, dimensions)

	for i := uint64(0); i < dimensions; i++ {
		values[i] = morton.compact(uint64(code) >> i)
	}

	return values
}

func (morton *Morton64) SUnpack(code int64) []int64 {
	uvalues := morton.Unpack(code)
	values := make([]int64, len(uvalues), len(uvalues))

	for i := 0; i < len(uvalues); i++ {
		values[i] = morton.unshiftSign(uvalues[i])
	}

	return values
}

func (morton *Morton64) dimensionsCheck(dimensions uint64) {
	if morton.dimensions != dimensions {
		panic(fmt.Sprintf("morton64 with %d dimensions received %d values", morton.dimensions, dimensions))
	}
}

func (morton *Morton64) valueCheck(value uint64) {
	if value >= (1 << morton.bits) {
		panic(fmt.Sprintf("morton64 with %d bits per dimension received %d to pack", morton.bits, value))
	}
}

func (morton *Morton64) shiftSign(value int64) uint64 {
	if value >= (1<<(morton.bits-1)) || value <= -(1<<(morton.bits-1)) {
		panic(fmt.Sprintf("morton64 with %d bits per dimension received signed %d to pack", morton.bits, value))
	}

	if value < 0 {
		value = -value
		value |= 1 << (morton.bits - 1)
	}
	return uint64(value)
}

func (morton *Morton64) unshiftSign(value uint64) int64 {
	sign := value & (1 << (morton.bits - 1))
	value &= (1 << (morton.bits - 1)) - 1
	svalue := int64(value)
	if sign != 0 {
		svalue = -svalue
	}
	return svalue
}

func (morton *Morton64) split(value uint64) uint64 {
	for o := 0; o < len(morton.masks); o++ {
		value = (value | (value << morton.lshifts[o])) & morton.masks[o]
	}

	return value
}

func (morton *Morton64) compact(code uint64) uint64 {
	for o := len(morton.masks) - 1; o >= 0; o-- {
		code = (code | (code >> morton.rshifts[o])) & morton.masks[o]
	}
	return code
}
