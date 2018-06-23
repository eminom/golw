package lwm2m

import (
	"math"
	"testing"
)

func runUno(t *testing.T, v interface{}, expectedLength int) {
	var val int64
	switch v.(type) {
	case int8:
		val = int64(v.(int8))
	case int16:
		val = int64(v.(int16))
	case int32:
		val = int64(v.(int32))
	case int64:
		val = int64(v.(int64))
	default:
		t.Fatalf("unknown type: %T", v)
	}
	res := prvEncodeInt(val)
	t.Logf("%T: %v", v, res)
	if len(res) != expectedLength {
		t.Fatalf("error: unexpected length for %T: %v,  %v but %v", v, v, expectedLength, len(res))
	}
}

func TestPrvEncodeInt(t *testing.T) {
	runUno(t, MaxInt8, 1)
	runUno(t, MinInt8, 1)
	runUno(t, MaxInt16, 2)
	runUno(t, MinInt16, 2)
	runUno(t, MaxInt32, 4)
	runUno(t, MinInt32, 4)
	runUno(t, MaxInt64, 8)
	runUno(t, MinInt64, 8)
}

func runDos(t *testing.T, a interface{}, expectedLength int) {
	var val float64
	switch a.(type) {
	case float32:
		val = float64(a.(float32))
	case float64:
		val = a.(float64)
	default:
		t.Fatalf("unknown type: %T", a)
	}
	res := prvEncodeFloat(val)
	t.Logf("%T: %v", a, res)
	if len(res) != expectedLength {
		t.Fatalf("error: unexpected length for %T: %v, %v but got %v", a, a, expectedLength, len(res))
	}
}

func TestPrvEncodeFloat(t *testing.T) {
	runDos(t, float32(math.MaxFloat32), 4)
	runDos(t, math.MaxFloat32*10, 8)
	runDos(t, math.MaxFloat32+1e25, 8) // Not so sure. But it is OK.
	runDos(t, float32(-math.MaxFloat32), 4)
	runDos(t, math.MaxFloat64, 8)
	runDos(t, -math.MaxFloat64, 8)
}
