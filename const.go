/* Refer to:
https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
*/

package lwm2m

const (
	MaxUint8 = ^uint8(0)
	MinUint8 = uint8(0)

	MaxInt8 = int8(MaxUint8 >> 1)
	MinInt8 = -MaxInt8 - 1

	MaxUint16 = ^uint16(0)
	MinUint16 = uint16(0)

	MaxInt16 = int16(MaxUint16 >> 1)
	MinInt16 = -MaxInt16 - 1

	MaxUint32 = ^uint32(0)
	MinUint32 = uint32(0)

	MaxInt32 = int32(MaxUint32 >> 1)
	MinInt32 = -MaxInt32 - 1

	MaxUint64 = ^uint64(0)
	MinUint64 = uint64(0)

	MaxInt64 = int64(MaxUint64 >> 1)
	MinInt64 = -MaxInt64 - 1
)
