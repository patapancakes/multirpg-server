package main

import (
	"fmt"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Packet struct {
	sender *Client
	data []byte
}

func (p *Packet) recieve() {
	packet, err := protocol.Decode(p.data)
	if err != nil {
		fmt.Println(err)
	}

	switch packet := packet.(type) {
	case protocol.SwitchRoom:
		err = p.handleSwitchRoom(packet)
	case protocol.Sprite:
		err = p.handleSprite(packet)
	case protocol.Position:
		err = p.handlePosition(packet)
	case protocol.Speed:
		err = p.handleSpeed(packet)
	default:
		err = fmt.Errorf("unknown packet type: %T", packet)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Packet) handleSwitchRoom(switchRoom protocol.SwitchRoom) error {
	if p.sender.room.server.rooms[switchRoom.Id] == nil {
		return fmt.Errorf("room not found")
	}

	p.sender.handleDisconnect() // remove from old room and broadcast disconnect packet
	p.sender.room = p.sender.room.server.rooms[switchRoom.Id] // set client room to new room
	p.sender.id = p.sender.room.getFreeId() // set client id to new room's free id
	p.sender.room.server.rooms[switchRoom.Id].clients[p.sender.id] = p.sender // add to new room
	p.sender.handleConnect() // get room data and broadcast connect packet

	return nil
}

func (p *Packet) handleSprite(sprite protocol.Sprite) error {
	if !isValidCharSet(string(sprite.Name)) {
		return nil
	}

	sprite.Id = p.sender.id

	packet, err := protocol.Encode(protocol.Sprite{})
	if err != nil {
		return err
	}

	p.sender.sprite = sprite.Name
	p.sender.spriteIndex = sprite.Index

	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handlePosition(position protocol.Position) error {
	if position.Direction > 3 {
		return fmt.Errorf("invalid direction")
	}

	position.Id = p.sender.id

	packet, err := protocol.Encode(protocol.Position{})
	if err != nil {
		return err
	}

	p.sender.x = position.X
	p.sender.y = position.Y
	p.sender.direction = position.Direction

	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handleSpeed(speed protocol.Speed) error {
	if speed.Speed > 10 {
		return fmt.Errorf("speed is too high")
	}

	speed.Id = p.sender.id

	packet, err := protocol.Encode(protocol.Speed{})
	if err != nil {
		return err
	}

	p.sender.speed = speed.Speed

	p.sender.room.broadcast(packet, p.sender)

	return nil
}
