package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	_ = iota
	Easy
	Medium
	Hard
)

func GenerateBoard(lvl int) *[9][9]int {
	var b [9][9]int
	Solve(0, 0, &b)
	deleteRandomCells(&b, lvl)
	return &b
}

func deleteRandomCells(board *[9][9]int, lvl int) {
	emptyCells := lvl * 20
	for i := 0; i < emptyCells; i++ {
		ranRow := rand.Intn(9)
		ranCol := rand.Intn(9)
		board[ranRow][ranCol] = 0
	}
}

func Solve(row, col int, board *[9][9]int) bool {
	if col == len(board[row]) {
		col = 0
		row++
	}
	if row == len(board) {
		return true
	}
	if board[row][col] != 0 {
		return Solve(row, col+1, board)
	}
	rand.Seed(time.Now().UnixNano())
	moves := rand.Perm(10)[:10]
	for i := range moves {
		if canPlace(moves[i], row, col, board) {
			board[row][col] = moves[i]
			if Solve(row, col+1, board) {
				return true
			}
		}
	}
	board[row][col] = 0
	return false
}

func canPlace(num, row, col int, board *[9][9]int) bool {
	// row constraint
	for r := range board {
		if board[r][col] == num {
			return false
		}
	}
	// col constraint
	for c := range board[row] {
		if board[row][c] == num {
			return false
		}
	}
	// square constraint
	sqSize := int(math.Sqrt(float64(len(board))))
	sqY := row / sqSize
	sqX := col / sqSize
	initRow := sqY * sqSize
	initCol := sqX * sqSize
	for r := 0; r < sqSize; r++ {
		for c := 0; c < sqSize; c++ {
			if num == board[initRow+r][initCol+c] {
				return false
			}
		}
	}
	return true
}

func Correct(board [9][9]int) bool {
	//row constraint
	for row := range board {
		seen := map[int]bool{}
		for col := range board {
			if seen[board[row][col]] || board[row][col] == 0 {
				return false
			}
			seen[board[row][col]] = true

		}
	}
	//col constraint
	for col := range board {
		seen := map[int]bool{}
		for row := range board {
			if seen[board[row][col]] {
				return false
			}
			seen[board[row][col]] = true

		}
	}

	//square constraint
	for col := range board {
		for row := range board {
			seen := map[int]bool{}
			sqSize := int(math.Sqrt(float64(len(board))))
			sqY := row / sqSize
			sqX := col / sqSize
			initRow := sqY * sqSize
			initCol := sqX * sqSize
			for r := 0; r < sqSize; r++ {
				for c := 0; c < sqSize; c++ {
					if seen[board[initRow+r][initCol+c]] {
						return false
					}
					seen[board[initRow+r][initCol+c]] = true
				}
			}
		}
	}
	return true
}
