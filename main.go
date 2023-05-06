package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	rows           = 6
	columns        = 7
	pause_autoplay = 0
)

type Board struct {
	grid     [rows][columns]string
	current  string
	player1  string
	player2  string
	gameOver bool
	winner   string
}

type MoveData struct {
	Player string
	Move   int
	State  [][]int
}

type GameData struct {
	Winner   string
	MoveData *MoveData
}

func main() {

	rand.Seed(time.Now().UnixNano())

	autoplay := flag.Bool("autoplay", false, "Enable autoplay mode")
	flag.Parse()

	board := Board{}

	if *autoplay {
		board.Autoplay()
	} else {
		board.Play()
	}

}

func (b *Board) Play() {
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
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			if b.grid[row][col] == "" {
				continue
			}

			player := b.grid[row][col]
			if b.checkWin(row, col, player) {
				b.gameOver = true
				b.winner = player
				return
			}
		}
	}
}

func (b *Board) checkWin(row, col int, player string) bool {
	directions := [][]int{
		{1, 0},  // Vertical
		{0, 1},  // Horizontal
		{1, 1},  // Diagonal up-right
		{-1, 1}, // Diagonal up-left
	}

	for _, direction := range directions {
		count := b.countInDirection(row, col, player, direction)
		if count >= 4 {
			return true
		}
	}
	return false
}

func (b *Board) countInDirection(row, col int, player string, direction []int) int {
	count := 1
	for i := 1; i < 4; i++ {
		newRow := row + i*direction[0]
		newCol := col + i*direction[1]

		if newRow < 0 || newRow >= rows || newCol < 0 || newCol >= columns {
			break
		}

		if b.grid[newRow][newCol] == player {
			count++
		} else {
			break
		}
	}
	return count
}

func (b *Board) Autoplay() {
	b.player1 = "X"
	b.player2 = "O"
	b.current = "X"

	gameDataList := []GameData{}

	b.PrintBoard()

	for !b.gameOver {
		move := b.AutoplayMove()
    b.CheckGameOver()
		moveData := &MoveData{
			Player: b.current,
			Move:   move,
			//Board:  b.grid,
      State:  b.GameState(),
		}

		gameData := GameData{
			Winner:   b.winner,
			MoveData: moveData,
		}
		gameDataList = append(gameDataList, gameData)

		b.SwitchPlayer()

		b.PrintBoard()
		time.Sleep(pause_autoplay * time.Millisecond)
	}

	if b.winner != "" {
		fmt.Printf("Game Over! Player %s wins!\n", b.winner)
	} else {
		fmt.Println("Game Over! It's a draw!")
	}
    
  // Replace the ioutil.WriteFile code with the following:
  file, err := os.OpenFile("game_data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
  	fmt.Println("Error opening file:", err)
  	return
  }
  defer file.Close()
  
  encoder := json.NewEncoder(file)
  err = encoder.Encode(gameDataList)
  if err != nil {
  	fmt.Println("Error writing JSON data to file:", err)
  	return
  }
}

func (b *Board) GameState() [][]int {
	state := make([][]int, rows)
	for i := range state {
		state[i] = make([]int, columns)
	}
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			switch b.grid[row][col] {
			case "":
				state[row][col] = 0
			case b.player1:
				state[row][col] = -1
			case b.player2:
				state[row][col] = 1
			}
		}
	}
	return state
}

func (b *Board) AutoplayMove() int {
	for {
		col := rand.Intn(columns)

		if b.IsValidMove(col) {
			err := b.PlayMove(col)
			if err != nil {
				fmt.Println(err)
				continue
			}
			return col
		}
	}
}

func (b *Board) IsValidMove(col int) bool {
	return b.grid[0][col] == ""
}

func (b *Board) PlaceToken(col int) {
	for row := rows - 1; row >= 0; row-- {
		if b.grid[row][col] == "" {
			b.grid[row][col] = b.current
			break
		}
	}
}
