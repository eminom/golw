package lwm2m

func NewString(id int, s string) DataItem {
	return DataItem{
		Type: TypeString,
		ID:   uint16(id),
		raw:  s,
	}
}

func NewInteger(id int, val int) DataItem {
	return DataItem{
		Type: TypeInteger,
		ID:   uint16(id),
		raw:  int64(val),
	}
}

func NewFloat(id int, val float64) DataItem {
	return DataItem{
		Type: TypeFloat,
		ID:   uint16(id),
		raw:  val,
	}
}

func NewArray(id int, arr []DataItem) DataItem {
	return DataItem{
		Type: TypeObjectInstance,
		ID:   uint16(id),
		raw:  arr,
	}
}
