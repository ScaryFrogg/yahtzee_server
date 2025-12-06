package service

import (
	"math/rand/v2"

	"github.com/ScaryFrogg/yahtzee_server/internal/types"
)

func Roll(board *types.Board, changes [6]bool) [6]int {
	if board.RollCount > 3 {
		return board.CurrentRoll
	}
	if board.RollCount == 0 {
		board.CurrentRoll = [6]int{rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1, rand.IntN(6) + 1}
	} else {
		for i, b := range changes {
			if b {
				board.CurrentRoll[i] = rand.IntN(6) + 1
			}
		}

	}
	board.RollCount++

	return board.CurrentRoll
}

func Calculate(board *types.Board) [types.BOARD_ROW_COUNT]int {
	result := [types.BOARD_ROW_COUNT]int{}
	for _, v := range board.CurrentRoll {
		//calculate sum fields
		result[v-1] = result[v-1] + v

		//TODO calculate other
	}
	board.CachedOptions = result

	return result
}

func Commit(board *types.Board, commitIndex int) {
	commitValue := board.CachedOptions[commitIndex]
	commitRow := board.Rows[commitIndex]

	if commitRow.Complete {
		return
	}

	commitRow.Fields[commitRow.CurrIndex] = commitValue
	commitRow.CurrIndex = commitRow.CurrIndex + 1
	if commitRow.CurrIndex > 5 {
		commitRow.Complete = true
	}
}
