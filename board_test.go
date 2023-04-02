package main

import (
	"testing"
)

func TestCheckGameOver(t *testing.T) {
	testCases := []struct {
		grid     [rows][columns]string
		gameOver bool
		winner   string
	}{
		{
			grid: [rows][columns]string{
				{"X", "", "", "", "", "", ""},
				{"X", "", "", "", "", "", ""},
				{"X", "", "", "", "", "", ""},
				{"X", "", "", "", "", "", ""},
				{"", "", "", "", "", "", ""},
				{"", "", "", "", "", "", ""},
			},
			gameOver: true,
			winner:   "X",
		},
		{
			grid: [rows][columns]string{
				{"X", "O", "X", "O", "X", "O", ""},
				{"X", "O", "X", "O", "X", "O", ""},
				{"O", "X", "O", "X", "O", "X", ""},
				{"X", "O", "X", "O", "X", "O", ""},
				{"O", "X", "O", "X", "O", "X", ""},
				{"", "", "", "", "", "", ""},
			},
			gameOver: false,
		},
	}

	for i, testCase := range testCases {
		board := Board{
			grid:    testCase.grid,
			player1: "X",
			player2: "O",
		}

		board.CheckGameOver()

		if board.gameOver != testCase.gameOver {
			t.Errorf("Test case %d: expected gameOver to be %v, got %v", i+1, testCase.gameOver, board.gameOver)
		}

		if testCase.gameOver && board.winner != testCase.winner {
			t.Errorf("Test case %d: expected winner to be %s, got %s", i+1, testCase.winner, board.winner)
		}
	}
}
