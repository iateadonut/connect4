package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)


type BoardStateMessage struct {
	BoardState [][]int `json:"board_state"`
	Message    string  `json:"message"`
  YourTurn   bool
  Done       bool
  Winner     bool
}

func main() {

    _ = godotenv.Load(".connect4.env") // Error ignored if the file doesn't exist

  // Get the URL and port from environment variables with default values
  url := os.Getenv("URL")
  if url == "" {
    url = "connect4.100wires.com"
  }
  port := os.Getenv("PORT")
  if port == "" {
    port = "51234"
  }

  conn, err := net.Dial("tcp", url+":"+port)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
  fmt.Println("Connected to " + url + ":" + port)
	defer conn.Close()

	go listenForServerMessages(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		_, err = conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Failed to send data to server:", err)
			os.Exit(1)
		}
	}
}


func listenForServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection lost:", err)
			os.Exit(1)
		}

		var boardStateMessage BoardStateMessage
		err = json.Unmarshal([]byte(message), &boardStateMessage)
		if err != nil {
			fmt.Println("Failed to parse server message:", err)
      fmt.Println(string(message))
			continue
		}

		printBoardState(boardStateMessage.BoardState)
		fmt.Println(boardStateMessage.Message)

		// if strings.HasPrefix(boardStateMessage.Message, "Player") {
		// 	fmt.Print("Enter column number (1-7): ")
		// }
	}
}

func printBoardState(boardState [][]int) {
  fmt.Println("|1||2||3||4||5||6||7|")
	for _, row := range boardState {
		for _, cell := range row {
			if cell == 0 {
				fmt.Print("| |")
			} else if cell == -1 {
				fmt.Print("|X|")
			} else {
				fmt.Print("|O|")
			}
		}
		fmt.Println()
		fmt.Println(strings.Repeat("---", len(row)))
	}
  fmt.Println("|1||2||3||4||5||6||7|")
	fmt.Println()
}

