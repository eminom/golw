package lwm2m

import (
	"encoding/binary"
	"math"
)

func prvEncodeInt(v int64) []byte {
	var buf [8]byte
	if v >= int64(MinInt8) && v <= int64(MaxInt8) {
		buf[0] = byte(v)
		return buf[:1]
	} else if v >= int64(MinInt16) && v <= int64(MaxInt16) {
		buf[0] = byte(v >> 8)
		buf[1] = byte(v)
		return buf[:2]
	} else if v >= int64(MinInt32) && v <= int64(MaxInt32) {
		buf[0] = byte(v >> 24)
		buf[1] = byte(v >> 16)
		buf[2] = byte(v >> 8)
		buf[3] = byte(v)
		return buf[:4]
	}

	var i uint
	for i = 0; i < 8; i++ {
		buf[i] = byte(v >> ((7 - i) * 8))
	}
	return buf[:]
}

func prvEncodeFloat(v float64) []byte {
	if v < 0.0-math.MaxFloat32 || v > math.MaxFloat32 {
		v64 := math.Float64bits(v)
		var ob [8]byte
		binary.BigEndian.PutUint64(ob[:], v64)
		return ob[:]
	}

	v32 := math.Float32bits(float32(v))
	var ob [4]byte
	binary.BigEndian.PutUint32(ob[:], v32)
	return ob[:]
}
