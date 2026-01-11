package main

import (
	"fmt"
)

type squareState int
type player int

const (
	none   = iota // 0
	cross         // 1
	circle        // 2
)

type gameState struct {
	board      [3][3]squareState
	turnPlayer player
}

func (e player) String() string {
	switch e {
	case none:
		return "none"
	case cross:
		return "cross"
	case circle:
		return "circle"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func (state *gameState) drawBoard() {
	for i, row := range state.board {
		for j, square := range row {
			fmt.Print(" ")
			switch square {
			case none:
				fmt.Print(" ")
			case cross:
				fmt.Print("X")
			case circle:
				fmt.Print("O")
			}
			if j != len(row)-1 {
				fmt.Print(" |")
			}
		}
		if i != len(state.board)-1 {
			fmt.Print("\n------------")
		}
		fmt.Print("\n")
	}
}

// error types
type markAlreadyExistsError struct {
	row    int
	column int
}

type positionOutOfBoundsError struct {
	row    int
	column int
}

// error method implementation
func (e *markAlreadyExistsError) Error() string {
	return fmt.Sprintf("position (%d,%d) already has a mark on it.", e.row, e.column)
}

func (e *positionOutOfBoundsError) Error() string {
	return fmt.Sprintf("position (%d,%d) is out of bounds.", e.row, e.column)
}

// placing mark at a position
func (state *gameState) placeMark(row int, column int) error {
	if row < 0 || column < 0 || row >= len(state.board) || column >= len(state.board[row]) {
		return &positionOutOfBoundsError{row, column}
	}
	if state.board[row][column] != none {
		return &markAlreadyExistsError{row, column}
	}

	state.board[row][column] = squareState(state.turnPlayer)

	return nil
}

type gameResult int

const (
	noWinner = iota
	crossWon
	circleWon
	draw
)

func (state *gameState) decideNext() player {
	return state.turnPlayer
}

func (state *gameState) nextTurn() {
	if state.turnPlayer == cross {
		state.turnPlayer = circle
	} else {
		state.turnPlayer = cross
	}
}

func (state *gameState) checkForWinner() gameResult {
	boardSize := len(state.board)

	// lambda func
	checkLine := func(startRow int, startColumn int, deltaRow int, deltaColumn int) gameResult {
		var lastSquare squareState = state.board[startRow][startColumn]
		row, column := startRow+deltaRow, startColumn+deltaColumn

		for row >= 0 && column >= 0 && row < boardSize && column < boardSize {

			if state.board[row][column] == none {
				return noWinner
			}

			if lastSquare != state.board[row][column] {
				return noWinner
			}

			lastSquare = state.board[row][column]
			row, column = row+deltaRow, column+deltaColumn
		}

		switch lastSquare {
		case cross:
			return crossWon
		case circle:
			return circleWon
		}

		return noWinner
	}

	// check horizontal rows
	for row := 0; row < boardSize; row++ {
		if result := checkLine(row, 0, 0, 1); result != noWinner {
			return result
		}
	}
	// check vertical columns
	for column := 0; column < boardSize; column++ {
		if result := checkLine(column, 0, 0, 1); result != noWinner {
			return result
		}
	}

	// check top-left to bottom-right diagonal
	if result := checkLine(0, 0, 1, 1); result != noWinner {
		return result
	}
	// check top-right to bottom-left diagonal
	if result := checkLine(0, boardSize-1, 1, -1); result != noWinner {
		return result
	}

	// check for draw
	for _, row := range state.board {
		for _, square := range row {
			if square == none {
				return noWinner
			}
		}
	}
	// if no one wins yet, but none of the squares are empty
	return draw
}

func main() {
	state := gameState{}
	state.turnPlayer = cross

	var result gameResult = noWinner

	for {
		fmt.Printf("next player to place a mark is: %v\n", state.decideNext())

		state.drawBoard()

		for {
			var row, column int

			fmt.Scan(&row, &column)

			e := state.placeMark(row, column)

			if e == nil {
				break
			}

			fmt.Println(e)
			fmt.Printf("please re-enter a position:\n> ")
		}

		result = state.checkForWinner()
		if result != noWinner {
			break
		}

		state.nextTurn()

		fmt.Println()
	}

	state.drawBoard()

	switch result {
	case crossWon:
		fmt.Printf("cross won the game!\n")
	case circleWon:
		fmt.Printf("circle won the game!\n")
	case draw:
		fmt.Printf("the game has ended with a draw!\n")
	}
}
