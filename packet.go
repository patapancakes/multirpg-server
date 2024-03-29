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

package main

import (
	"bytes"
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

	if p.sender.lobby == nil {
		switch packet := packet.(type) {
		case protocol.NewLobby:
			err = p.handleNewLobby(packet)
		case protocol.JoinLobby:
			err = p.handleJoinLobby(packet)
		default:
			err = fmt.Errorf("bad packet type for server realm: %T", packet)
		}
	} else if p.sender.room == nil {
		switch packet := packet.(type) {
		case protocol.SwitchRoom:
			err = p.handleSwitchRoom(packet)
		default:
			err = fmt.Errorf("bad packet type for lobby realm: %T", packet)
		}
	} else {
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
			err = fmt.Errorf("bad packet type for room realm: %T", packet)
		}
	}
	if err != nil {
		fmt.Printf("client %d packet error: %s\n", p.sender.id, err)
	}
}

func (p *Packet) handleNewLobby(newLobby protocol.NewLobby) error {
	lobbyCode := generateLobbyCode()

	for {
		if _, ok := p.sender.server.lobbies.Load(lobbyCode); !ok {
			break
		}

		lobbyCode = generateLobbyCode()
	}

	p.sender.server.lobbies.Store(lobbyCode, p.sender.server.createLobby(newLobby.GameHash))

	p.sender.joinLobby(lobbyCode)

	packet, _ := protocol.Encode(protocol.NewLobbyR{
		LobbyCode: []byte(lobbyCode),
	})

	p.sender.sendPacket(packet)

	return nil
}

func (p *Packet) handleJoinLobby(joinLobby protocol.JoinLobby) error {
	lobby, ok := p.sender.server.lobbies.Load(string(joinLobby.LobbyCode))
	if !ok {
		return fmt.Errorf("invalid lobby code: %s", joinLobby.LobbyCode)
	}

	if gameHash := lobby.(*Lobby).gameHash; !bytes.Equal(gameHash, joinLobby.GameHash) {
		return fmt.Errorf("game hash mismatch: %s and %s", gameHash, joinLobby.GameHash)
	}

	p.sender.joinLobby(string(joinLobby.LobbyCode))

	return nil
}

func (p *Packet) handleSwitchRoom(switchRoom protocol.SwitchRoom) error {
	if _, ok := p.sender.lobby.rooms.Load(switchRoom.Id); !ok {
		p.sender.lobby.rooms.Store(switchRoom.Id, p.sender.lobby.createRoom(switchRoom.Id))
	}

	// Remove client from old room and broadcast client leave packet
	if p.sender.room != nil {
		p.sender.leaveRoom()
	}

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
