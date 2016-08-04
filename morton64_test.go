package morton

import "testing"

func compareValues(t *testing.T, dimensions uint64, bits uint64, value uint64, unpacked uint64) {
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

func doTestValueBoundries(t *testing.T, dimensions uint64, bits uint64, value uint64) {
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
	m.Pack(values)
}

func TestValueBoundries(t *testing.T) {
	doTestValueBoundries(t, 2, 1, 2)
	doTestValueBoundries(t, 16, 4, 16)
}

func doTest2(t *testing.T, bits uint64, value0 uint64, value1 uint64) {
	m := Make64(2, bits)
	code := m.Pack2(value0, value1)
	unpacked0, unpacked1 := m.Unpack2(code)
	compareValues(t, 2, bits, value0, unpacked0)
	compareValues(t, 2, bits, value1, unpacked1)
}

func TestPackUnpack2(t *testing.T) {
	doTest2(t, 32, 1, 2)
	doTest2(t, 32, 2, 1)
	doTest2(t, 32, (1<<32)-1, (1<<32)-1)
	doTest2(t, 1, 1, 1)
}

func doTest3(t *testing.T, bits uint64, value0 uint64, value1 uint64, value2 uint64) {
	m := Make64(3, bits)
	code := m.Pack3(value0, value1, value2)
	unpacked0, unpacked1, unpacked2 := m.Unpack3(code)
	compareValues(t, 3, bits, value0, unpacked0)
	compareValues(t, 3, bits, value1, unpacked1)
	compareValues(t, 3, bits, value2, unpacked2)
}

func TestPackUnpack3(t *testing.T) {
	doTest3(t, 21, 1, 2, 4)
	doTest3(t, 21, 4, 2, 1)
	doTest3(t, 21, (1<<21)-1, (1<<21)-1, (1<<21)-1)
	doTest3(t, 1, 1, 1, 1)
}

func doTest4(t *testing.T, bits uint64, value0 uint64, value1 uint64, value2 uint64, value3 uint64) {
	m := Make64(4, bits)
	code := m.Pack4(value0, value1, value2, value3)
	unpacked0, unpacked1, unpacked2, unpacked3 := m.Unpack4(code)
	compareValues(t, 4, bits, value0, unpacked0)
	compareValues(t, 4, bits, value1, unpacked1)
	compareValues(t, 4, bits, value2, unpacked2)
	compareValues(t, 4, bits, value3, unpacked3)
}

func TestPackUnpack4(t *testing.T) {
	doTest4(t, 16, 1, 2, 4, 8)
	doTest4(t, 16, 8, 4, 2, 1)
	doTest4(t, 16, (1<<16)-1, (1<<16)-1, (1<<16)-1, (1<<16)-1)
	doTest4(t, 1, 1, 1, 1, 1)
}

func doTestArray(t *testing.T, dimensions uint64, bits uint64, values []uint64) {
	m := Make64(dimensions, bits)
	code := m.Pack(values)
	unpacked := m.Unpack(code)
	if len(values) != len(unpacked) {
		t.Errorf("%d values transformed to %d values after pack/unpack with %d dimensions and %d bits", len(values), len(unpacked), dimensions, bits)
	}
	for i := uint64(0); i < dimensions; i++ {
		compareValues(t, dimensions, bits, values[i], unpacked[i])
	}
}

func TestPackUnpackArray(t *testing.T) {
	doTestArray(t, 6, 10, []uint64{1, 2, 4, 8, 16, 32})
	doTestArray(t, 6, 10, []uint64{32, 16, 8, 4, 2, 1})
	doTestArray(t, 6, 10, []uint64{63, 63, 63, 63, 63, 63})
	values := make([]uint64, 64)
	for i := 0; i < 64; i++ {
		values[i] = 1
	}
	doTestArray(t, 64, 1, values)
}

func doTestPackArrayDimensions(t *testing.T, dimensions uint64, bits uint64, size uint64) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("unexpected pack %d values to %d dimensions", size, dimensions)
		}
	}()

	values := make([]uint64, size)
	m := Make64(dimensions, bits)
	m.Pack(values)
}

func TestPackArrayDimensions(t *testing.T) {
	doTestPackArrayDimensions(t, 2, 32, 3)
	doTestPackArrayDimensions(t, 2, 32, 1)
}
