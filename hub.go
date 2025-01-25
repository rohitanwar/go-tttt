// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Hub struct {
	// Registered clients.
	//clients map[*Client]int

	games []*Game

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (h *Hub) findEmptyGame() int {
	for i := 0; i < len(h.games); i++ {
		if (h.games[i].status == Waiting) || (h.games[i].status == Empty) {
			return i
		}
	}
	return -1
}

func newHub() *Hub {
	var games []*Game

	for i := 0; i < 100; i++ {
		games = append(games, newGame())
	}
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		games:      games,
		//clients:    make(map[*Client]int),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			{
				gameId := h.findEmptyGame()
				if gameId == -1 {
					panic("Hub Full!!")
				}
				// h.clients[client] = gameId
				client.gameId = gameId
				h.games[gameId].registerPlayer(client)

				select {
				case client.send <- []byte("connected " + strconv.Itoa(gameId) + " " + strconv.Itoa(client.player)):
				default:
					close(client.send)
				}
			}
		}
	}
}
