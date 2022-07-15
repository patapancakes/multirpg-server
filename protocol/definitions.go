package protocol

const (
	CONNECT     uint8 = 0x01
	DISCONNECT  uint8 = 0x02
	SWITCH_ROOM uint8 = 0x03
	SPRITE      uint8 = 0x10
	POSITION    uint8 = 0x11
	SPEED       uint8 = 0x12
)

//0x01
type Connect struct {
	Id uint16
}

//0x02
type Disconnect struct {
	Id uint16
}

//0x03
type SwitchRoom struct {
	Id uint16
}

//0x10
type Sprite struct {
	Id uint16
	Name []byte
	Index uint8
}

//0x11
type Position struct {
	Id uint16
	X uint16
	Y uint16
	Direction uint8
}

//0x12
type Speed struct {
	Id uint16
	Speed uint8
}
