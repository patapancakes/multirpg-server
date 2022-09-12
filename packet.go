package main

/*
multirpg-server
https://github.com/Gamizard/multirpg-server

Copyright (C) 2022 azarashi <azarashi@majestaria.fun>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"fmt"

	"github.com/Gamizard/multirpg-server/protocol"
)

type Packet struct {
	sender *Client
	data   []byte
}

func (p *Packet) process() {
	packet, err := protocol.Decode(p.data)
	if err != nil {
		fmt.Printf("client %d packet error: %s\n", p.sender.id, err)
		return
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
		fmt.Printf("client %d packet error: %s\n", p.sender.id, err)
	}
}

func (p *Packet) handleSwitchRoom(switchRoom protocol.SwitchRoom) error {
	if p.sender.room.server.rooms[switchRoom.Id] == nil {
		p.sender.room.server.createRoom(switchRoom.Id)
	}

	// Remove client from old room and broadcast disconnect packet
	p.sender.leaveRoom()

	// Initialize client variables so other clients entering the new room don't get the old values
	// Redundant most of the time but prevents some visual weirdness
	p.sender.x = 0
	p.sender.y = 0
	p.sender.direction = 0
	p.sender.speed = 0

	p.sender.joinRoom(switchRoom.Id)

	return nil
}

func (p *Packet) handleSprite(sprite protocol.Sprite) error {
	p.sender.sprite = sprite.Name
	p.sender.spriteIndex = sprite.Index

	sprite.Id = p.sender.id
	packet, _ := protocol.Encode(sprite)
	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handlePosition(position protocol.Position) error {
	if position.Direction > 3 {
		return fmt.Errorf("invalid direction: %d", position.Direction)
	}

	p.sender.x = position.X
	p.sender.y = position.Y
	p.sender.direction = position.Direction

	position.Id = p.sender.id
	packet, _ := protocol.Encode(position)
	p.sender.room.broadcast(packet, p.sender)

	return nil
}

func (p *Packet) handleSpeed(speed protocol.Speed) error {
	if speed.Speed > 10 {
		return fmt.Errorf("speed is too high: %d", speed.Speed)
	}

	p.sender.speed = speed.Speed

	speed.Id = p.sender.id
	packet, _ := protocol.Encode(speed)
	p.sender.room.broadcast(packet, p.sender)

	return nil
}
