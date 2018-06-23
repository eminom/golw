package lwm2m

type UriT struct {
	Flag   int
	ObjID  int
	InstID int
	ResID  int
}

const (
	ResourceBit = 0x01
	InstanceBit = 0x02
	ObjectBit   = 0x04
)

// Same implementation as liblwm2m.h
func (u *UriT) IsResourceSet() bool {
	return (u.Flag & ResourceBit) != 0
}

func (u *UriT) IsInstanceSet() bool {
	return (u.Flag & InstanceBit) != 0
}

func (u *UriT) IsObjectSet() bool {
	return (u.Flag & ObjectBit) != 0
}

func (u *UriT) SetResourceID(id int) {
	u.Flag |= ResourceBit
	u.ResID = id
}

func (u *UriT) SetInstanceID(id int) {
	u.Flag |= InstanceBit
	u.InstID = id
}

func (u *UriT) SetObjectID(id int) {
	u.Flag |= ObjectBit
	u.ObjID = id
}
