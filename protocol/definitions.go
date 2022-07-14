package protocol

const (
	CONNECT = 0x01
	DISCONNECT = 0x02
	SWITCH_ROOM = 0x03
	SPRITE = 0x10
	MOVE = 0x11
	SPEED = 0x12
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
type Move struct {
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
