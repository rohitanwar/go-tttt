package main

import "errors"

type GameStatus int

const (
	Full    GameStatus = 2
	Waiting GameStatus = 1
	Empty   GameStatus = 0
)

type Game struct {
	status  GameStatus
	playerX *Client
	playerY *Client
	board   Board
}

func (g *Game) makeMove(player Player, globalIndex int, localIndex int) error {
	if g.status != Full {
		return errors.New("game not started")
	}

	err := g.board.makeMove(player, globalIndex, localIndex)

	g.broadcastState()

	if err != nil {
		return err
	}

	return nil
}

func (g *Game) registerPlayer(client *Client) {
	if g.status == Empty {
		g.playerX = client
		g.status = Waiting
		client.player = -1
	} else if g.status == Waiting {
		g.playerY = client
		g.status = Full
		client.player = 1
	} else {
		panic("Unreachable State! Too many players in one game!")
	}
}

func newGame() *Game {
	return &Game{
		status:  Empty,
		playerX: nil,
		playerY: nil,
		board:   newBoard(),
	}
}

func (g *Game) broadcastState() {
	select {
	case g.playerX.send <- g.board.getState():
	default:
		close(g.playerX.send)
	}
}
