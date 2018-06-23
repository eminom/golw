package lwm2m

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var (
	chunkError       = errors.New("marshalling error")
	typeError        = errors.New("type error")
	unknownTypeError = errors.New("unknown type error")
)

type DataType uint16

const (
	TypeUndefined DataType = iota
	TypeObject
	TypeObjectInstance
	TypeMultipleResource

	TypeString
	TypeOpaque
	TypeInteger
	TypeFloat
	TypeBoolean

	TypeObjectLink
)

type DataItem struct {
	Type DataType
	ID   uint16 // Resource ID
	raw  interface{}
}

type ObjectLink struct {
	ObjectID         uint16
	ObjectInstanceID uint16
}

func TlvMarshalItems(uri *UriT, items []DataItem) ([]byte, error) {
	if len(items) <= 0 {
		return nil, nil
	}

	var isResInstance bool
	if (uri != nil && uri.IsResourceSet()) &&
		(len(items) > 1 || uri.ResID != int(items[0].ID)) {
		isResInstance = true
	} else {
		isResInstance = false
	}

	var rv []byte
	for _, v := range items {
		c, e := v.MarshalResource(isResInstance)
		if e != nil {
			return nil, e
		}
		rv = append(rv, c...)
	}
	return rv, nil
}

/*
[]byte buffer
array of objects
objectID - objectInstanceID
*/

func (d *DataItem) MarshalResource(isResourceInstance bool) ([]byte, error) {
	var rv []byte

	isInstance := isResourceInstance

	// TypeObject is consider error
	switch d.Type {
	case TypeMultipleResource:
		isInstance = true
		fallthrough
	case TypeObjectInstance:
		items, e := d.ToDataItemArray()
		if e != nil {
			return nil, e
		}
		var xchunk []byte
		for _, v := range items {
			chunk0, e := v.MarshalResource(isInstance)
			if e != nil {
				return nil, e
			}
			xchunk = append(xchunk, chunk0...)
		}
		hdrBuff := prvCreateHeader(false, d.Type, d.ID, len(xchunk))
		rv = append(hdrBuff, xchunk...)

	case TypeObjectLink:
		ol, e := d.ToObjLink()
		if e != nil {
			return nil, e
		}
		var idBuff [4]byte
		binary.BigEndian.PutUint16(idBuff[:], ol.ObjectID)
		binary.BigEndian.PutUint16(idBuff[2:], ol.ObjectInstanceID)
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, 4)
		rv = append(hdrBuff, idBuff[:]...)

	case TypeString, TypeOpaque:
		buff, e := d.ToBufferBytes()
		if e != nil {
			return nil, e
		}
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, len(buff))
		rv = append(hdrBuff, buff...)

	case TypeInteger:
		buff, e := d.ToIntBytes()
		if e != nil {
			return nil, e
		}
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, len(buff))
		rv = append(hdrBuff, buff...)

	case TypeFloat:
		buff, e := d.ToFloatBytes()
		if e != nil {
			return nil, e
		}
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, len(buff))
		rv = append(hdrBuff, buff...)

	case TypeBoolean:
		buff, e := d.ToBooleanBytes()
		if e != nil {
			return nil, e
		}
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, len(buff))
		rv = append(hdrBuff, buff...)
	default:
		return nil, unknownTypeError
	}
	return rv, nil
}

func (d *DataItem) ToDataItemArray() ([]DataItem, error) {
	switch d.raw.(type) {
	case []DataItem:
		return d.raw.([]DataItem), nil
	}
	return nil, typeError
}

func (d *DataItem) ToObjLink() (ObjectLink, error) {
	if ol, ok := d.raw.(ObjectLink); ok {
		return ol, nil
	}
	return ObjectLink{}, typeError
}

func (d *DataItem) ToBufferBytes() ([]byte, error) {
	switch v := d.raw.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, typeError
	}
}

// Pay attention: only int64 is allowed here.
func (d *DataItem) ToIntBytes() ([]byte, error) {
	switch v := d.raw.(type) {
	case int64:
		return prvEncodeInt(v), nil
	}
	return nil, typeError
}

// Pay attention: only float64 is allowed here.
func (d *DataItem) ToFloatBytes() ([]byte, error) {
	switch v := d.raw.(type) {
	case float64:
		return prvEncodeFloat(v), nil
	}
	return nil, typeError
}

func (d *DataItem) ToBooleanBytes() ([]byte, error) {
	if b, ok := d.raw.(bool); ok {
		var v int64 = 0
		if b {
			v = 1
		}
		return prvEncodeInt(v), nil
	}
	return nil, typeError
}

func (d *DataItem) SetInteger(v int64) {
	d.Type = TypeInteger
	d.raw = v
}

func (d *DataItem) ToInteger() (int64, error) {
	if v, ok := d.raw.(int64); ok {
		return v, nil
	}
	return 0, typeError
}

func (d *DataItem) SetString(v string) {
	d.Type = TypeString
	d.raw = v
}

func (d *DataItem) ToString() (string, error) {
	switch v := d.raw.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	default:
		return "", typeError
	}
}

func (d *DataItem) SetBinary(raw []byte) {
	d.Type = TypeOpaque
	d.raw = bytes.Repeat(raw, 1)
}

func (d *DataItem) ToBinary() ([]byte, error) {
	switch d.Type {
	case TypeString, TypeOpaque:
		return d.ToBufferBytes()
	case TypeInteger:
		return d.ToIntBytes()
	case TypeFloat:
		return d.ToFloatBytes()
	case TypeBoolean:
		return d.ToBooleanBytes()

		//TODO: FIXME: support more types to binary
	default:
		return nil, typeError
	}
}
