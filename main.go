package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	rows    = 6
	columns = 7
)

type Board struct {
	grid     [rows][columns]string
	current  string
	player1  string
	player2  string
	gameOver bool
	winner   string
}

func main() {
	board := Board{
		player1: "X",
		player2: "O",
		current: "X",
	}

	board.PrintBoard()

	reader := bufio.NewReader(os.Stdin)

	for !board.gameOver {
		fmt.Printf("Player %s, enter column number (1-%d): ", board.current, columns)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		col, err := strconv.Atoi(text)
		if err != nil || col < 1 || col > columns {
			fmt.Println("Invalid input, please try again.")
			continue
		}

		err = board.PlayMove(col - 1)
		if err != nil {
			fmt.Println(err)
			continue
		}

		board.PrintBoard()
		board.CheckGameOver()
		board.SwitchPlayer()
	}

	fmt.Printf("Game over! %s wins!\n", board.winner)
}

func (b *Board) PlayMove(col int) error {
	for row := rows - 1; row >= 0; row-- {
		if b.grid[row][col] == "" {
			b.grid[row][col] = b.current
			return nil
		}
	}
	return fmt.Errorf("column %d is full", col+1)
}

func (b *Board) SwitchPlayer() {
	if b.current == b.player1 {
		b.current = b.player2
	} else {
		b.current = b.player1
	}
}

func (b *Board) PrintBoard() {
	fmt.Println()
	fmt.Print("|1||2||3||4||5||6||7|")
	fmt.Println()
	for _, row := range b.grid {
		for _, cell := range row {
			if cell == "" {
				fmt.Print("| |")
			} else {
				fmt.Printf("|%s|", cell)
			}
		}
		fmt.Println()
		fmt.Println(strings.Repeat("---", columns))
	}
	fmt.Print("|1||2||3||4||5||6||7|")
	fmt.Println()
}

func (b *Board) CheckGameOver() {
	// Check for win conditions here
	// If a player has won, set b.gameOver = true and b.winner to the winning player
}
