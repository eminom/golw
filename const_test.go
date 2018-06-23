package lwm2m

import (
	"testing"
)

func TestConsts(t *testing.T) {
	t.Logf("min %T: %v", MinInt8, MinInt8)
	t.Logf("max %T: %v", MaxInt8, MaxInt8)
	t.Logf("min %T: %v", MinInt16, MinInt16)
	t.Logf("max %T: %v", MaxInt16, MaxInt16)
	t.Logf("min %T: %v", MinInt32, MinInt32)
	t.Logf("max %T: %v", MaxInt32, MaxInt32)
	t.Logf("min %T: %v", MinInt64, MinInt64)
	t.Logf("max %T: %v", MaxInt64, MaxInt64)
}
