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
	rshifts := make([]uint64, 0)

	masks = append(masks, mask)
	lshifts = append(lshifts, 0)
	rshifts = append(rshifts, shift>>1)

	for shift > 0 {
		mask = 0
		shifted := uint64(0)

		for bit := uint64(0); bit < bits; bit++ {
			distance := (dimensions * bit) - bit
			shifted |= shift & distance
			mask |= 1 << bit << (((shift - 1) ^ 0xffffffffffffffff) & distance)
		}

		if shifted != 0 {
			masks = append(masks, mask)
			lshifts = append(lshifts, shift)
			rshifts = append(rshifts, (shift >> 1))
		}

		shift >>= 1
	}

	rshifts[(len(rshifts) - 1)] = 0

	return &Morton64{dimensions: dimensions, bits: bits, masks: masks, lshifts: lshifts, rshifts: rshifts}
}

func (morton *Morton64) Pack(values []uint64) uint64 {
	dimensions := uint64(len(values))
	morton.dimensionsCheck(dimensions)
	for i := uint64(0); i < dimensions; i++ {
		morton.valueCheck(values[i])
	}

	code := uint64(0)
	for i := uint64(0); i < dimensions; i++ {
		code |= morton.split(values[i]) << i
	}

	return code
}

func (morton *Morton64) Pack2(value0 uint64, value1 uint64) uint64 {
	morton.dimensionsCheck(2)
	morton.valueCheck(value0)
	morton.valueCheck(value1)

	return morton.split(value0) | (morton.split(value1) << 1)
}

func (morton *Morton64) Pack3(value0 uint64, value1 uint64, value2 uint64) uint64 {
	morton.dimensionsCheck(3)
	morton.valueCheck(value0)
	morton.valueCheck(value1)
	morton.valueCheck(value2)

	return morton.split(value0) | (morton.split(value1) << 1) | (morton.split(value2) << 2)
}

func (morton *Morton64) Pack4(value0 uint64, value1 uint64, value2 uint64, value3 uint64) uint64 {
	morton.dimensionsCheck(4)
	morton.valueCheck(value0)
	morton.valueCheck(value1)
	morton.valueCheck(value2)
	morton.valueCheck(value3)

	return morton.split(value0) | (morton.split(value1) << 1) | (morton.split(value2) << 2) | (morton.split(value3) << 3)
}

func (morton *Morton64) Unpack(code uint64) []uint64 {
	dimensions := morton.dimensions

	values := make([]uint64, dimensions, dimensions)

	for i := uint64(0); i < dimensions; i++ {
		values[i] = morton.compact(code >> i)
	}

	return values
}

func (morton *Morton64) Unpack2(code uint64) (uint64, uint64) {
	morton.dimensionsCheck(2)

	value0 := morton.compact(code)
	value1 := morton.compact(code >> 1)

	return value0, value1
}

func (morton *Morton64) Unpack3(code uint64) (uint64, uint64, uint64) {
	morton.dimensionsCheck(3)

	value0 := morton.compact(code)
	value1 := morton.compact(code >> 1)
	value2 := morton.compact(code >> 2)

	return value0, value1, value2
}

func (morton *Morton64) Unpack4(code uint64) (uint64, uint64, uint64, uint64) {
	morton.dimensionsCheck(4)

	value0 := morton.compact(code)
	value1 := morton.compact(code >> 1)
	value2 := morton.compact(code >> 2)
	value3 := morton.compact(code >> 3)

	return value0, value1, value2, value3
}

func (morton *Morton64) dimensionsCheck(dimensions uint64) {
	if morton.dimensions != dimensions {
		panic(fmt.Sprintf("morton with %d dimensions received %d values", morton.dimensions, dimensions))
	}
}

func (morton *Morton64) valueCheck(value uint64) {
	if value >= (1 << morton.bits) {
		panic(fmt.Sprintf("morton with %d bits per dimension received %d to pack", morton.bits, value))
	}
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
