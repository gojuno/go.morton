package morton

import "testing"

func BenchmarkSPackUnpack2(b *testing.B) {
	m := Make64(2, 31)
	for n := 0; n < b.N; n++ {
		code := m.SPack(12345, 67890)
		m.SUnpack(code)
	}
}

func compareValues(t *testing.T, dimensions uint64, bits uint64, value uint64, unpacked uint64) {
	if unpacked != value {
		t.Errorf("%d transformed to %d after pack/unpack with %d dimensions and %d bits", value, unpacked, dimensions, bits)
	}
}

func compareSValues(t *testing.T, dimensions uint64, bits uint64, value int64, unpacked int64) {
	if unpacked != value {
		t.Errorf("%d transformed to %d after pack/unpack with %d dimensions and %d bits", value, unpacked, dimensions, bits)
	}
}

func doTestBadMake64(t *testing.T, dimenstions uint64, bits uint64) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected make with %d dimensions and %d bits", dimenstions, bits)
		}
	}()

	Make64(dimenstions, bits)
}

func TestMake64(t *testing.T) {
	doTestBadMake64(t, 0, 1)
	doTestBadMake64(t, 1, 0)
	doTestBadMake64(t, 1, 65)
}

func doTestValueBoundaries(t *testing.T, dimensions uint64, bits uint64, value uint64) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected %d packed with %d dimensions and %d bits", value, dimensions, bits)
		}
	}()

	m := Make64(dimensions, bits)
	values := make([]uint64, dimensions)
	for i := 0; i < len(values); i++ {
		values[i] = 0
	}
	values[0] = value
	m.Pack(values...)
}

func TestValueBoundaries(t *testing.T) {
	doTestValueBoundaries(t, 2, 1, 2)
	doTestValueBoundaries(t, 16, 4, 16)
}

func doTestSValueBoundaries(t *testing.T, dimensions uint64, bits uint64, value int64) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected %d packed with %d dimensions and %d bits", value, dimensions, bits)
		}
	}()

	m := Make64(dimensions, bits)
	values := make([]int64, dimensions)
	for i := 0; i < len(values); i++ {
		values[i] = 0
	}
	values[0] = value
	m.SPack(values...)
}

func TestSValueBoundaries(t *testing.T) {
	doTestSValueBoundaries(t, 2, 2, 2)
	doTestSValueBoundaries(t, 2, 2, -2)
	doTestSValueBoundaries(t, 16, 4, 8)
	doTestSValueBoundaries(t, 16, 4, -8)
}

func doTestPackUnpack(t *testing.T, dimensions uint64, bits uint64, values ...uint64) {
	m := Make64(dimensions, bits)
	code := m.Pack(values...)
	unpacked := m.Unpack(code)
	if len(values) != len(unpacked) {
		t.Errorf("%d values transformed to %d values after pack/unpack with %d dimensions and %d bits", len(values), len(unpacked), dimensions, bits)
	}
	for i := uint64(0); i < dimensions; i++ {
		compareValues(t, dimensions, bits, values[i], unpacked[i])
	}
}

func TestPackUnpackArray(t *testing.T) {
	doTestPackUnpack(t, 2, 32, 1, 2)
	doTestPackUnpack(t, 2, 32, 2, 1)
	doTestPackUnpack(t, 2, 32, (1<<32)-1, (1<<32)-1)
	doTestPackUnpack(t, 2, 1, 1, 1)

	doTestPackUnpack(t, 3, 21, 1, 2, 4)
	doTestPackUnpack(t, 3, 21, 4, 2, 1)
	doTestPackUnpack(t, 3, 21, (1<<21)-1, (1<<21)-1, (1<<21)-1)
	doTestPackUnpack(t, 3, 1, 1, 1, 1)

	doTestPackUnpack(t, 4, 16, 1, 2, 4, 8)
	doTestPackUnpack(t, 4, 16, 8, 4, 2, 1)
	doTestPackUnpack(t, 4, 16, (1<<16)-1, (1<<16)-1, (1<<16)-1, (1<<16)-1)
	doTestPackUnpack(t, 4, 1, 1, 1, 1, 1)

	doTestPackUnpack(t, 6, 10, 1, 2, 4, 8, 16, 32)
	doTestPackUnpack(t, 6, 10, 32, 16, 8, 4, 2, 1)
	doTestPackUnpack(t, 6, 10, 1023, 1023, 1023, 1023, 1023, 1023)

	values := make([]uint64, 64)
	for i := 0; i < 64; i++ {
		values[i] = 1
	}
	doTestPackUnpack(t, 64, 1, values...)
}

func doTestSPackUnpack(t *testing.T, dimensions uint64, bits uint64, values ...int64) {
	m := Make64(dimensions, bits)
	code := m.SPack(values...)
	unpacked := m.SUnpack(code)
	if len(values) != len(unpacked) {
		t.Errorf("%d values transformed to %d values after pack/unpack with %d dimensions and %d bits", len(values), len(unpacked), dimensions, bits)
	}
	for i := uint64(0); i < dimensions; i++ {
		compareSValues(t, dimensions, bits, values[i], unpacked[i])
	}
}

func TestSPackUnpack(t *testing.T) {
	doTestSPackUnpack(t, 2, 32, 1, 2)
	doTestSPackUnpack(t, 2, 32, 2, 1)
	doTestSPackUnpack(t, 2, 32, (1<<31)-1, (1<<31)-1)
	doTestSPackUnpack(t, 2, 2, 1, 1)
	doTestSPackUnpack(t, 2, 32, -1, -2)
	doTestSPackUnpack(t, 2, 32, -2, -1)
	doTestSPackUnpack(t, 2, 32, -((1 << 31) - 1), -((1 << 31) - 1))
	doTestSPackUnpack(t, 2, 2, -1, -1)

	doTestSPackUnpack(t, 3, 21, 1, 2, 4)
	doTestSPackUnpack(t, 3, 21, 4, 2, 1)
	doTestSPackUnpack(t, 3, 21, (1<<20)-1, (1<<20)-1, (1<<20)-1)
	doTestSPackUnpack(t, 3, 2, 1, 1, 1)
	doTestSPackUnpack(t, 3, 21, -1, -2, -4)
	doTestSPackUnpack(t, 3, 21, -4, -2, -1)
	doTestSPackUnpack(t, 3, 21, -((1 << 20) - 1), -((1 << 20) - 1), -((1 << 20) - 1))
	doTestSPackUnpack(t, 3, 2, -1, -1, -1)

	doTestSPackUnpack(t, 4, 16, 1, 2, 4, 8)
	doTestSPackUnpack(t, 4, 16, 8, 4, 2, 1)
	doTestSPackUnpack(t, 4, 16, (1<<15)-1, (1<<15)-1, (1<<15)-1, (1<<15)-1)
	doTestSPackUnpack(t, 4, 2, 1, 1, 1, 1)
	doTestSPackUnpack(t, 4, 16, -1, -2, -4, -8)
	doTestSPackUnpack(t, 4, 16, -8, -4, -2, -1)
	doTestSPackUnpack(t, 4, 16, -((1 << 15) - 1), -((1 << 15) - 1), -((1 << 15) - 1), -((1 << 15) - 1))
	doTestSPackUnpack(t, 4, 2, -1, -1, -1, -1)

	doTestSPackUnpack(t, 6, 10, 1, 2, 4, 8, 16, 32)
	doTestSPackUnpack(t, 6, 10, 32, 16, 8, 4, 2, 1)
	doTestSPackUnpack(t, 6, 10, 511, 511, 511, 511, 511, 511)
	doTestSPackUnpack(t, 6, 10, -1, -2, -4, -8, -16, -32)
	doTestSPackUnpack(t, 6, 10, -32, -16, -8, -4, -2, -1)
	doTestSPackUnpack(t, 6, 10, -511, -511, -511, -511, -511, -511)

	values := make([]int64, 32)
	for i := 0; i < 32; i++ {
		values[i] = int64(1 - 2*(i%2))
	}
	doTestSPackUnpack(t, 32, 2, values...)
}

func doTestPackDimensions(t *testing.T, dimensions uint64, bits uint64, size uint64) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected pack %d values to %d dimensions", size, dimensions)
		}
	}()

	values := make([]uint64, size)
	m := Make64(dimensions, bits)
	m.Pack(values...)
}

func TestPackDimensions(t *testing.T) {
	doTestPackDimensions(t, 2, 32, 3)
	doTestPackDimensions(t, 2, 32, 1)
}
