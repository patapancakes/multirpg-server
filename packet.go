package main

import (
	"fmt"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Packet struct {
	sender *Client
	data []byte
}

func (p *Packet) receive() {
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

	p.sender.leaveRoom() // remove from old room and broadcast disconnect packet

	// Initialize client variables so other clients entering the new room don't get the old values
	// Redundant most of the time but prevents some visual weirdness
	p.sender.x = 0
	p.sender.y = 0
	p.sender.direction = 0
	p.sender.speed = 0

	p.sender.room = p.sender.room.server.rooms[switchRoom.Id] // set client room to new room
	p.sender.room.server.rooms[switchRoom.Id].clients[p.sender] = true // add to new room
	p.sender.joinRoom() // get room data and broadcast connect packet

	return nil
}

func (p *Packet) handleSprite(sprite protocol.Sprite) error {
	if !isValidCharSet(string(sprite.Name)) {
		return nil
	}

	p.sender.sprite = sprite.Name
	p.sender.spriteIndex = sprite.Index

	sprite.Id = p.sender.id
	packet, _ := protocol.Encode(protocol.Sprite{})
	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handlePosition(position protocol.Position) error {
	if position.Direction > 3 {
		return fmt.Errorf("invalid direction")
	}

	p.sender.x = position.X
	p.sender.y = position.Y
	p.sender.direction = position.Direction

	position.Id = p.sender.id
	packet, _ := protocol.Encode(protocol.Position{})
	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handleSpeed(speed protocol.Speed) error {
	if speed.Speed > 10 {
		return fmt.Errorf("speed is too high")
	}

	p.sender.speed = speed.Speed

	speed.Id = p.sender.id
	packet, _ := protocol.Encode(protocol.Speed{})
	p.sender.room.broadcast(packet, p.sender)

	return nil
}
