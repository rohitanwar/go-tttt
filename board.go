package main

import (
	"errors"
	"fmt"
)

type Player int
type CellStates int
type Grid [9]CellStates

type Move struct {
	globalIndex int
	localIndex  int
}

const (
	PlayedX   CellStates = CellStates(PlayerX)
	PlayedY   CellStates = CellStates(PlayerY)
	CellEmpty CellStates = 0
)

const (
	PlayerX Player = -1
	PlayerY Player = 1
)

type Board struct {
	globalBoard  Grid
	localBoards  [9]Grid
	next_player  Player
	active_board int
	moveHistory  []Move
}

// Returns all the possible win conditions on a `Grid`.
func getWinStates() [12][3]int32 {
	return [12][3]int32{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{1, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{1, 4, 8},
		{2, 4, 6},
		{3, 7, 2},
		{7, 5, 2},
		{1, 5, 6},
		{1, 3, 8},
	}
	// `getWinStates` is a function due to no constant arrays in Go.
}

// Takes a `Grid` and checks if it is winning for a `Player`
func (g *Grid) isWinning(player Player) bool {
	winStates := getWinStates()
	for _, winState := range winStates {
		if (g[winState[0]] == g[winState[1]]) &&
			(g[winState[1]] == g[winState[2]]) &&
			(g[winState[2]] == CellStates(player)) {
			return true
		}
	}
	return false
}

func (b *Board) makeMove(player Player, globalIndex int, localIndex int) error {

	if globalIndex < 0 || globalIndex > 8 {
		return errors.New("invalid global index")
	}

	if localIndex < 0 || localIndex > 8 {
		return errors.New("invalid local index")
	}

	if b.next_player != player {
		return errors.New("wrong player")
	}

	if (b.active_board != -1) && (b.active_board != globalIndex) {
		return errors.New("can't play here")
	}

	if b.globalBoard[globalIndex] != CellEmpty {
		return errors.New("can't play here")
	}

	if b.localBoards[globalIndex][localIndex] != CellEmpty {
		return errors.New("can't play here")
	}

	// Make the move
	b.localBoards[globalIndex][localIndex] = CellStates(player)

	// Check if current move led to a win and update the `globalBoard`
	if b.localBoards[localIndex].isWinning(player) {
		b.globalBoard[localIndex] = CellStates(player)
	}

	// Update the `active_board`
	if b.globalBoard[localIndex] != CellEmpty {
		b.active_board = -1
	} else {
		b.active_board = localIndex
	}

	b.next_player *= -1

	b.moveHistory = append(b.moveHistory, Move{globalIndex: globalIndex, localIndex: localIndex})

	return nil
}

func (b *Board) getState() []byte {
	s := make([]byte, 0)
	s = append(s, []byte("state: ")...)
	for i := 0; i < 9; i++ {
		t := fmt.Sprintf("%02d", b.globalBoard[i])
		s = append(s, []byte(t)...)
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			t := fmt.Sprintf("%02d", b.localBoards[i][j])
			s = append(s, []byte(t)...)
		}
	}

	t := fmt.Sprintf("%02d", b.next_player)
	s = append(s, []byte(t)...)

	t = fmt.Sprintf("%02d", b.active_board)
	s = append(s, []byte(t)...)

	fmt.Println(s)
	return s
}

func (b *Board) reset() {
	b.globalBoard.reset()
	for _, lb := range b.localBoards {
		lb.reset()
	}

	b.next_player = PlayerX
	b.active_board = -1
	b.moveHistory = make([]Move, 0, 100)
}

func (g *Grid) reset() {
	for i := 0; i < len(g); i++ {
		g[i] = CellEmpty
	}
}

func newBoard() Board {
	return Board{
		globalBoard: newGrid(),
		localBoards: [9]Grid{
			newGrid(), newGrid(), newGrid(),
			newGrid(), newGrid(), newGrid(),
			newGrid(), newGrid(), newGrid(),
		},
		next_player:  PlayerX,
		active_board: -1,
		moveHistory:  make([]Move, 0, 100),
	}
}

func newGrid() Grid {
	return Grid{
		CellEmpty, CellEmpty, CellEmpty,
		CellEmpty, CellEmpty, CellEmpty,
		CellEmpty, CellEmpty, CellEmpty,
	}
}
