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
	//return PrvTlvTypeUnknown
}

func prvCreateHeader(isInstance bool, kind DataType, id uint16, dataLength int) []byte {

	hdrBuff := prvCreateHeaderBuffer(id, dataLength)
	var hdrType PrvType
	if isInstance {
		hdrType = PrvTlvTypeResourceInstance
	} else {
		hdrType = mapToHeaderType(kind)
	}

	hdrBuff[0] |= byte(hdrType & PrvTlvTypeMask)

	//TODO:FIXME
	return nil
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
		length += 3
	} else if dataLength > 0xFF {
		length += 2
	} else if dataLength > 7 {
		length += 1
	}
	return length
}
