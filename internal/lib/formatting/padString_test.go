package formatting

import "testing"

func TestPadString(t *testing.T) {
	result := PadString("test", 10)
	expectedResult := "   test   "

	if result != expectedResult {
		t.Fatalf("Test 1. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadString2(t *testing.T) {
	result := PadString("strings", 15)
	expectedResult := "    strings    "

	if result != expectedResult {
		t.Fatalf("Test 2. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringLowest(t *testing.T) {
	result := PadString("height", 5)
	expectedResult := "height"

	if result != expectedResult {
		t.Fatalf("Test 3. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringExtra(t *testing.T) {
	result := PadString("200", 6)
	expectedResult := " 200  "

	if result != expectedResult {
		t.Fatalf("Test 4. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringExtra2(t *testing.T) {
	result := PadString("swing", 10)
	expectedResult := "  swing   "

	if result != expectedResult {
		t.Fatalf("Test 5. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringExtra3(t *testing.T) {
	result := PadString("686.125µs", 14)
	expectedResult := "  686.125µs   "

	if result != expectedResult {
		t.Fatalf("Test 6. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringExtra4(t *testing.T) {
	result := PadString("2.950375ms", 14)
	expectedResult := "  2.950375ms  "

	if result != expectedResult {
		t.Fatalf("Test 7. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}
func TestPadStringExtra5(t *testing.T) {
	result := PadString("33.417µs", 14)
	expectedResult := "   33.417µs   "

	if result != expectedResult {
		t.Fatalf("Test 8. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}

func TestPadStringExtra6(t *testing.T) {
	result := PadString("357.666µs", 14)
	expectedResult := "  357.666µs   "

	if result != expectedResult {
		t.Fatalf("Test 9. String not equals: result(%v) != expectedResult(%v)", result, expectedResult)
	}
}
