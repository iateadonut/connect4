package main

import (
	"bufio"
	"encoding/json"

	// "flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	rows           = 6
	columns        = 7
	pause_autoplay = 0
)

type Board struct {
	grid     [rows][columns]string
	current  string
	player1  Player
	player2  Player
	gameOver bool
	winner   string
}

type BoardStateMessage struct {
	BoardState [][]int `json:"board_state"`
	Message    string  `json:"message"`
}

type Player struct {
    Symbol string
    State  string
    Conn   net.Conn
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

// func main() {

// 	rand.Seed(time.Now().UnixNano())

// 	autoplay := flag.Bool("autoplay", false, "Enable autoplay mode")
//   api := flag.Bool("api", false, "Return state as array")
// 	flag.Parse()

//   board := Board{api:*api}
//   // fmt.Println(board.api)

// 	if *autoplay {
// 		board.Autoplay()
// 	} else {
// 		board.Play()
// 	}

// }

func main() {
    _ = godotenv.Load(".connect4.env") // Error ignored if the file doesn't exist

    // Get the URL and port from environment variables with default values
    url := os.Getenv("URL")
    if url == "" {
      url = "localhost"
    }
    port := os.Getenv("PORT")
    if port == "" {
      port = "51234"
    }

    // Use the variables in your program...
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
      log.Fatalf("Could not start server: %s", err.Error())
    }
    defer ln.Close()

    log.Printf("Server is running on %s:%s", url, port)

    var player1 Player
    var player2 Player

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }

        if player1.Conn == nil {
            player1 = Player{
                Symbol: "X",
                State:  "waiting",
                Conn:   conn,
            }
        } else if player2.Conn == nil {
            player2 = Player{
                Symbol: "O",
                State:  "waiting",
                Conn:   conn,
            }

            go startNewGame(player1, player2)
        }
    }
}

func startNewGame(player1, player2 Player) {

	board := Board{
		// grid:     [rows][columns]string,
		current:  player1.Symbol,
		player1:  player1,
		player2:  player2,
		gameOver: false,
		// winner:   "",
	}

	board.PrintBoard()

	for !board.gameOver {
		err := board.PlayTurn()
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

func (b *Board) PlayTurn() error {
	var currentPlayer Player
	if b.current == b.player1.Symbol {
		currentPlayer = b.player1
	} else {
		currentPlayer = b.player2
	}

	boardState := b.GameState()
	message := fmt.Sprintf("Player %s, enter column number (1-%d):", b.current, columns)
	boardStateMessage := BoardStateMessage{BoardState: boardState, Message: message}

	jsonData, err := json.Marshal(boardStateMessage)
	if err != nil {
		return fmt.Errorf("Error marshaling board state message: %w", err)
	}

	_, err = currentPlayer.Conn.Write(append(jsonData, '\n'))
	if err != nil {
		return fmt.Errorf("Error writing to player connection: %w", err)
	}

	// fmt.Fprintf(currentPlayer.Conn, "Player %s, enter column number (1-%d): \n", b.current, columns)
	
	reader := bufio.NewReader(currentPlayer.Conn)
	text, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Error reading from player: %w", err)
	}

	text = strings.TrimSpace(text)
	col, err := strconv.Atoi(text)
	if err != nil || col < 1 || col > columns {
		fmt.Fprintln(currentPlayer.Conn, "Invalid input, please try again.")
		return fmt.Errorf("Invalid input")
	}

	err = b.PlayMove(col - 1)
	if err != nil {
		fmt.Fprintln(currentPlayer.Conn, err)
		return err
	}

	// Send the board state back to the current player without a message
  boardState = b.GameState()
	boardStateMessage = BoardStateMessage{BoardState: boardState, Message: ""}

  jsonData, err = json.Marshal(boardStateMessage)
	if err != nil {
		return fmt.Errorf("Error marshaling board state message: %w", err)
	}

	_, err = currentPlayer.Conn.Write(append(jsonData, '\n'))
	if err != nil {
		return fmt.Errorf("Error writing to player connection: %w", err)
	}

	return nil
}

// func (b *Board) Play(conn net.Conn) {
// 	// board := Board{
// 	// 	player1: "X",
// 	// 	player2: "O",
// 	// 	current: "X",
//  //    api: b.api,
// 	// }

// 	// board.PrintBoard()

// 	reader := bufio.NewReader(conn)

// 	for !board.gameOver {
// 		fmt.Printf("Player %s, enter column number (1-%d): ", board.current, columns)
// 		text, _ := reader.ReadString('\n')
// 		text = strings.TrimSpace(text)
// 		col, err := strconv.Atoi(text)
// 		if err != nil || col < 1 || col > columns {
// 			fmt.Println("Invalid input, please try again.")
// 			continue
// 		}

// 		err = board.PlayMove(col - 1)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}

// 		board.PrintBoard()
// 		board.CheckGameOver()
// 		board.SwitchPlayer()
// 	}

// 	fmt.Printf("Game over! %s wins!\n", board.winner)
// }

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
	if b.current == b.player1.Symbol {
		b.current = b.player2.Symbol
	} else {
		b.current = b.player1.Symbol
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
  filledCells := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			if b.grid[row][col] == "" {
				continue
			}
      filledCells++

			player := b.grid[row][col]
			if b.checkWin(row, col, player) {
				b.gameOver = true
				b.winner = player
				return
			}
		}
	}

  if filledCells == rows*columns {
    b.gameOver = true
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
	b.player1.Symbol = "X"
	b.player2.Symbol = "O"
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
			case b.player1.Symbol:
				state[row][col] = -1
			case b.player2.Symbol:
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
