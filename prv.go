package lwm2m

type PrvType uint8

const (
	PrvTlvTypeUnknown          PrvType = 0xFF
	PrvTlvTypeObject                   = 0x10
	PrvTlvTypeObjectInstance           = 0x00
	PrvTlvTypeResource                 = 0xC0
	PrvTlvTypeMultipleResource         = 0x80
	PrvTlvTypeResourceInstance         = 0x40
)

const (
	PrvTlvHeaderMaxLength = 6
	PrvTlvTypeMask        = 0xC0 // 1100-0000
)

func mapToHeaderType(kind DataType) PrvType {
	switch kind {
	case TypeObject:
		return PrvTlvTypeObject
	case TypeObjectInstance:
		return PrvTlvTypeObjectInstance
	case TypeMultipleResource:
		return PrvTlvTypeMultipleResource

	case TypeString, TypeInteger, TypeFloat, TypeBoolean, TypeOpaque, TypeObjectLink:
		return PrvTlvTypeResource
	default:
		fallthrough
	case TypeUndefined:
		return PrvTlvTypeUnknown
	}
}

// normally: Mask 0xC0
func getDataType(kind uint8) DataType {
	cat := PrvType(kind)
	switch cat {
	case PrvTlvTypeObject:
		//0x10, which is not valid after masked
		//somebody FIXME??
		return TypeObject

	case PrvTlvTypeObjectInstance:
		//0x00
		return TypeObjectInstance

	case PrvTlvTypeMultipleResource:
		//0x80
		return TypeMultipleResource

	case PrvTlvTypeResource, PrvTlvTypeResourceInstance:
		//0xC0, 0x40
		return TypeOpaque

	default:
		return TypeUndefined
	}
}

func prvCreateHeader(isInstance bool, kind DataType, id uint16, dataLength int) []byte {

	ob := prvCreateHeaderBuffer(id, dataLength)
	var hdrType PrvType
	if isInstance {
		hdrType = PrvTlvTypeResourceInstance
	} else {
		hdrType = mapToHeaderType(kind)
	}

	var offset int

	ob[0] |= byte(hdrType & PrvTlvTypeMask)
	if id > 0xFF {
		ob[0] |= 0x20
		ob[1] = byte((id >> 8) & 0xFF)
		ob[2] = byte(id & 0xFF)
		offset = 3

	} else {
		ob[1] = byte(id)
		offset = 2
	}

	if dataLength <= 7 {
		// 0000 0000 (00)
		ob[0] += byte(dataLength)

	} else if dataLength <= 0xFF {
		// 0000 1000 (01)
		ob[0] |= 0x08
		ob[offset] = byte(dataLength)

	} else if dataLength <= 0xFFFF {
		// 0001 0000 (10)
		ob[0] |= 0x10
		ob[offset] = byte((dataLength >> 8) & 0xFF)
		ob[offset+1] = byte(dataLength & 0xFF)

	} else if dataLength <= 0xFFFFFF {
		// 0001 1000 (11)
		ob[0] |= 0x18
		ob[offset] = byte((dataLength >> 16) & 0xFF)
		ob[offset+1] = byte((dataLength >> 8) & 0xFF)
		ob[offset+2] = byte(dataLength & 0xFF)

	}

	return ob
}

func prvCreateHeaderBuffer(id uint16, dataLength int) []byte {
	return make([]byte, prvGetHeaderLength(id, dataLength))
}

func prvGetHeaderLength(id uint16, dataLength int) int {
	length := 2
	if id > 0xFF {
		length += 1
	}
	if dataLength > 0xFFFF {
		// three bytes: 01-00-00 to ff-ff-ff
		length += 3
	} else if dataLength > 0xFF {
		// two bytes: 01-00 to ff-ff
		length += 2
	} else if dataLength > 7 {
		// one byte: 00-10 to ff
		length += 1
	}
	return length
}
