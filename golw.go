package lwm2m

import (
	"encoding/binary"
	"errors"
)

var (
	chunkError = errors.New("marshalling error")
	typeError  = errors.New("type error")
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
		buff, e := d.ToBuffer()
		if e != nil {
			return nil, e
		}
		hdrBuff := prvCreateHeader(isInstance, d.Type, d.ID, len(buff))
		rv = append(hdrBuff, buff...)

	case TypeInteger:
		//TODO
	case TypeFloat:
		//TODO
	case TypeBoolean:
		//TODO

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

func (d *DataItem) ToBuffer() ([]byte, error) {
	if ob, ok := d.raw.([]byte); ok {
		return ob, nil
	}
	return nil, typeError
}
