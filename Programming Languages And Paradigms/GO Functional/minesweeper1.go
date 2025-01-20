package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Board [][]string

// createBlankBoard creates an empty board filled with "+" symbols.
func createBlankBoard(size int) Board {
	return createBoard(size, "+")
}

// createBoard creates a board filled with a specified string.
func createBoard(size int, fill string) Board {
	board := make(Board, size)
	for i := range board {
		board[i] = make([]string, size)
		for j := range board[i] {
			board[i][j] = fill
		}
	}
	return board
}

// createComparisonBoard generates the internal board with randomly placed mines.
func createComparisonBoard(size, mines int) [][]int {
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}

	mineCount := 0
	rand.Seed(time.Now().UnixNano())
	for mineCount < mines {
		row := rand.Intn(size)
		col := rand.Intn(size)
		if board[row][col] == 0 {
			board[row][col] = 1
			mineCount++
		}
	}
	return board
}

// displayBoard prints the user's board with column and row headers.
func displayBoard(board Board, size, mines, colWidth int) {
	fmt.Printf("There are %d mines to be found!\n\n", mines)
	fmt.Println("  " + headerRow(size, colWidth))
	for i, row := range board {
		fmt.Printf("%2d  %s\n", i, formatRow(row, colWidth))
	}
}

// headerRow generates the header row for the board display.
func headerRow(size, colWidth int) string {
	header := "  "
	for i := 0; i < size; i++ {
		header += fmt.Sprintf("%*d ", colWidth, i)
	}
	return header
}

// formatRow formats a single row of the board for display.
func formatRow(row []string, colWidth int) string {
	formatted := ""
	for _, cell := range row {
		formatted += fmt.Sprintf("%*s ", colWidth, cell)
	}
	return formatted
}

// updateBoard processes the user's input and updates the game state accordingly.
func updateBoard(userBoard Board, compareBoard [][]int, size, mines, unTouchedSpaces int, shield bool, userInput string) (Board, [][]int, int, bool, bool) {
	if unTouchedSpaces == mines {
		fmt.Println("You have won!")
		return resetGame(size, mines)
	}

	// Check for special commands
	switch userInput {
	case "S":
		fmt.Println("Shield activated! Avoiding the next mine.")
		return userBoard, compareBoard, unTouchedSpaces, true, false
	case "H":
		userBoard = giveHint(userBoard, compareBoard, size)
		return userBoard, compareBoard, unTouchedSpaces, shield, false
	}

	// Attempt to parse input as row and column
	inputRow, inputCol, err := parseInput(userInput)
	if err != nil || !validInput(inputRow, inputCol, size) {
		fmt.Println("Invalid input. Try again.")
		return userBoard, compareBoard, unTouchedSpaces, shield, false
	}

	// Process the input
	if compareBoard[inputRow][inputCol] == 1 {
		if shield {
			userBoard[inputRow][inputCol] = "S"
			fmt.Println("Shield used! You avoided the mine!")
			return userBoard, compareBoard, unTouchedSpaces, false, false
		} else {
			fmt.Println("You hit a mine! Game over.")
			return resetGame(size, mines)
		}
	} else {
		userBoard, unTouchedSpaces = revealTiles(userBoard, compareBoard, size, inputRow, inputCol, unTouchedSpaces)
	}
	return userBoard, compareBoard, unTouchedSpaces, shield, false
}

// resetGame resets the game state to its initial configuration.
func resetGame(size, mines int) (Board, [][]int, int, bool, bool) {
	return createBlankBoard(size), createComparisonBoard(size, mines), size * size, false, false
}

// parseInput parses the user's input into row and column indices.
func parseInput(input string) (int, int, error) {
	var row, col int
	_, err := fmt.Sscanf(input, "%d %d", &row, &col)
	return row, col, err
}

// validInput checks whether the given row and column are within the board's bounds.
func validInput(row, col, size int) bool {
	return row >= 0 && row < size && col >= 0 && col < size
}

// revealTiles uncovers tiles recursively based on surrounding mine counts.
func revealTiles(userBoard Board, compareBoard [][]int, size, x, y, unTouchedSpaces int) (Board, int) {
	if userBoard[x][y] != "+" {
		return userBoard, unTouchedSpaces
	}

	unTouchedSpaces--
	surroundingMines := checkSurround(compareBoard, x, y, size)

	if surroundingMines == 0 {
		userBoard[x][y] = " "
		neighbors := []struct{ dx, dy int }{
			{-1, -1}, {-1, 0}, {-1, 1},
			{0, -1} /*      */, {0, 1},
			{1, -1}, {1, 0}, {1, 1},
		}
		for _, neighbor := range neighbors {
			nx, ny := x+neighbor.dx, y+neighbor.dy
			if validInput(nx, ny, size) {
				userBoard, unTouchedSpaces = revealTiles(userBoard, compareBoard, size, nx, ny, unTouchedSpaces)
			}
		}
	} else {
		userBoard[x][y] = fmt.Sprintf("%d", surroundingMines)
	}
	return userBoard, unTouchedSpaces
}

// checkSurround counts the number of mines surrounding a specific tile.
func checkSurround(board [][]int, x, y, size int) int {
	count := 0
	neighbors := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1},
	}
	for _, neighbor := range neighbors {
		nx, ny := x+neighbor.dx, y+neighbor.dy
		if validInput(nx, ny, size) {
			count += board[nx][ny]
		}
	}
	return count
}

// giveHint reveals a random mine temporarily to assist the player.
func giveHint(userBoard Board, compareBoard [][]int, size int) Board {
	revealed := []struct{ x, y int }{}
	for i := 0; i < rand.Intn(3)+1; i++ {
		for {
			row := rand.Intn(size)
			col := rand.Intn(size)
			if compareBoard[row][col] == 1 {
				revealed = append(revealed, struct{ x, y int }{row, col})
				userBoard[row][col] = "M"
				break
			}
		}
	}
	displayBoard(userBoard, size, 0, 2)
	time.Sleep(3 * time.Second)
	for _, cell := range revealed {
		userBoard[cell.x][cell.y] = "+"
	}
	return userBoard
}

// main is the entry point for the game, handling user input and game loop.
func main() {
	size := 10
	mines := 1
	userBoard, compareBoard, unTouchedSpaces, shield, _ := resetGame(size, mines)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		displayBoard(userBoard, size, mines, 2)
		fmt.Print("Enter row and column (or 'S' for shield, 'H' for hint): ")
		if scanner.Scan() {
			userInput := strings.TrimSpace(scanner.Text())
			userBoard, compareBoard, unTouchedSpaces, shield, _ = updateBoard(userBoard, compareBoard, size, mines, unTouchedSpaces, shield, userInput)
		} else {
			fmt.Println("Follow input requirements! e.g. 1 4")
		}
	}
}
