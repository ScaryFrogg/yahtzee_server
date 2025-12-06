package types

import (
	"fmt"
	"strings"
)

const BOARD_ROW_COUNT = 6

type Board struct {
	Rows          [BOARD_ROW_COUNT]Row `json:"rows"`
	CurrentRoll   [6]int               `json:"currentRoll"`
	RollCount     int                  `json:"rollCount"`
	Waiting       bool                 `json:"waiting"`
	CachedOptions [BOARD_ROW_COUNT]int
}

type Row struct {
	Fields    [6]int
	CurrIndex int
	Complete  bool
}

func NewBoard() *Board {
	board := &Board{
		RollCount: 0,
		Waiting:   true,
	}

	for i := range BOARD_ROW_COUNT {
		board.Rows[i] = Row{
			CurrIndex: 0,
			Complete:  false,
		}
	}

	return board
}

type RowIndex int

const (
	Row1        RowIndex = 0
	Row2        RowIndex = 1
	Row3        RowIndex = 2
	Row4        RowIndex = 3
	Row5        RowIndex = 4
	Row6        RowIndex = 5
	RowMax      RowIndex = 6
	RowMim      RowIndex = 7
	RowStraight RowIndex = 8
	RowTODO     RowIndex = 9
)

func LogPlayerBoard(player *Player) {
	playerId := player.Id
	board := player.Board
	fmt.Printf("\n┌─ Player: %s\n", playerId)
	fmt.Printf("│  Current Roll (Count: %d): ", board.RollCount)
	for i, die := range board.CurrentRoll {
		fmt.Printf("[%d]", die)
		if i < len(board.CurrentRoll)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println("│")

	for rowIdx, row := range board.Rows {
		status := " "
		if row.Complete {
			status = "✓"
		}
		fmt.Printf("│  Row %2d [%s]: ", rowIdx, status)

		for fieldIdx := range 6 {
			if fieldIdx < row.CurrIndex {
				fmt.Printf("[%2d]", row.Fields[fieldIdx])
			} else if fieldIdx == row.CurrIndex {
				fmt.Print("[>>]")
			} else {
				fmt.Print("[  ]")
			}
			if fieldIdx < 5 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println("└" + strings.Repeat("─", 70))
}
