package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Command constants for user input
const (
	CmdReveal = "reveal"
	CmdFlag   = "flag"
	CmdQuit   = "quit"
)

// Board struct represents the game board
type Board struct {
	Width, Height int
	Cells         [][]Cell
}

// Cell struct represents a single cell on the game board
type Cell struct {
	IsMine   bool
	AdjMines int
	Revealed bool
	Flagged  bool
}

// Time complexity considerations:
// The majority of our methods are either O(1) or O(n), where n is the number of cells (width * height)
// I think complexity is mostly optimized given the constraints of the problem, excessive nested loops are avoided to prevent quadratic time complexity.
// RevealCell is O(1) if revealing a single cell, but can also be O(n) as in the worst case it can recursively reveal all adjacent cells with no adjacent mines.
// RevealCell could be further optimized to avoid deep recursion on larger boards by using a stack or queue to store cells to be revealed.

// This method creates a new board with the given width, height, and number of mines.
func NewBoard(width, height, mines int) *Board {
	board := &Board{Width: width, Height: height}
	// Create a 2D slice of cells
	board.Cells = make([][]Cell, height)
	for i := range board.Cells {
		// Initialize each cell in the board
		board.Cells[i] = make([]Cell, width)
	}
	// Place mines on the board and calculate the number of adjacent mines for each cell
	board.placeMines(mines)
	board.calculateAdjMines()
	return board
}

// placeMines places the specified number of mines randomly on the board.
func (b *Board) placeMines(mines int) {
	availableCells := b.Width * b.Height
	// Ensure the number of mines is within a reasonable range
	// Full-functionality for this would be configurable difficulty
	// For now, we'll just limit the number of mines to no more than the amount you requested (5)
	maxMines := availableCells * 5 / 9
	if mines > maxMines {
		mines = maxMines
	}

	// Create a slice of all possible positions
	positions := make([][2]int, 0, availableCells)
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			positions = append(positions, [2]int{x, y})
		}
	}

	// The first version of this code during the interview attempted to place mines by randomly selecting positions on the board and checking if a mine was already placed at that position. If it didn't, it would place a mine. This approach was inefficient and could result in an infinite loop if the number of mines was close to the total number of cells on the board. I refactored the code to shuffle the positions slice and place mines in the first N positions, where N is the number of mines. This approach guarantees that the number of mines placed is equal to the number requested and avoids the inefficiency of the original approach.
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < mines; i++ {
		x, y := positions[i][0], positions[i][1]
		b.Cells[y][x].IsMine = true
		//fmt.Printf("Mine placed at (%d, %d)\n", x, y)
	}
	//fmt.Println("Finished placing mines.")
}

// This method iterates over each cell in the board and calculates the number of adjacent mines for each cell.
// If the cell is a mine, the adjacent mines count is not calculated. The countAdjMines method is used to calculate the number of adjacent mines for each cell.
func (b *Board) calculateAdjMines() {
	for y := range b.Cells {
		for x := range b.Cells[y] {
			if !b.Cells[y][x].IsMine {
				b.Cells[y][x].AdjMines = b.countAdjMines(x, y)
			}
		}
	}
}

// countAdjMines counts the number of mines adjacent to the given cell.
// This method does a lot of necessary iteration, but I refactored this code after the interview to avoid excessive nesting and make it easier to read.
// The offsets slice contains the values -1, 0, and 1. The code then iterates over the Cartesian product of these values to get the adjacent cells.
func (b *Board) countAdjMines(x, y int) int {
	count := 0
	var offsets = []int{-1, 0, 1}

	for _, i := range offsets {
		for _, j := range offsets {
			if i == 0 && j == 0 {
				continue
			}
			// Check if the adjacent cell is a mine and increment the count if it is. The isValidCell method is used to check if the adjacent cell is within the bounds of the board.
			adjX, adjY := x+i, y+j
			if b.isValidCell(adjX, adjY) {
				if b.Cells[adjY][adjX].IsMine {
					count++
				}
			}
		}
	}
	return count
}

// This method checks if the given coordinates are within the bounds of the board.
func (b *Board) isValidCell(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// This method reveals a cell on the board. If the cell is a mine, the method returns true, indicating that the game is over.
// If the cell is not a mine and has no adjacent mines, the method recursively reveals the adjacent cells.
func (b *Board) RevealCell(x, y int) bool {
	if !b.isValidCell(x, y) || b.Cells[y][x].Revealed {
		return false
	}
	b.Cells[y][x].Revealed = true
	if b.Cells[y][x].IsMine {
		return true
	}
	if b.Cells[y][x].AdjMines == 0 {
		// Reveal adjacent cells if the current cell has no adjacent mines
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				b.RevealCell(x+i, y+j)
			}
		}
	}
	return false
}

// This method toggles the flag on a cell. If the cell is already revealed, the flag is not toggled.
func (b *Board) FlagCell(x, y int) {
	if !b.isValidCell(x, y) || b.Cells[y][x].Revealed {
		return
	}
	b.Cells[y][x].Flagged = !b.Cells[y][x].Flagged
}

// This method checks if the player has won the game. If all safe cells are revealed, the player wins.
func (b *Board) CheckWin() bool {
	for _, row := range b.Cells {
		for _, cell := range row {
			// If a safe cell is not revealed, the game continues
			if !cell.IsMine && !cell.Revealed {
				return false
			}
		}
	}
	// All safe cells are revealed, victory!
	return true
}

// This method is called when the game is over. It prints the final state of the board, revealing all mines.
func (b *Board) GameOver(showMines bool) {
	fmt.Println("Game over!")
	b.PrintBoard(showMines)
}

// This method prints the current state of the board. If showMines is true, all mines are revealed.
// The board is printed with the following symbols:
// * for mines
// F for flagged cells
// . for unrevealed safe cells
// An int representing adjacent mine count for revealed safe cells
// The board is printed row by row, with each cell separated by a space.
func (b *Board) PrintBoard(showMines bool) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Revealed {
				if cell.IsMine {
					fmt.Print("* ")
				} else {
					fmt.Printf("%d ", cell.AdjMines)
				}
			} else if cell.Flagged {
				fmt.Print("F ")
			} else if showMines && cell.IsMine {
				fmt.Print("M ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

// Debug method to print the board with mines and adjacent mine counts
// This was used for testing the board generation functions
// placeMines() and calculateAdjMines()
// It's also useful as a guide to validate the win condition
// I left this in for reference, in case you want to verify yourself
func (b *Board) PrintBoardDebug() {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.IsMine {
				fmt.Print("* ")
			} else {
				fmt.Printf("%d ", cell.AdjMines)
			}
		}
		fmt.Println()
	}
}

func main() {
	// Given a board of size 3x3 with 5 mines
	width, height, mines := 3, 3, 5
	board := NewBoard(width, height, mines)

	// Goals of printing: count the mines & ensure 'mines' amount, check adj. counts
	// DEBUG FUNCTION
	//board.PrintBoardDebug()

	// Game loop
	// Read user input via the console and execute commands
	// Initially, I used fmt.Scan to read user input, but this method was blocking and doesn't allow for easy exit. It also was less robust for handling inputs. I switched to bufio.Scanner to allow for non-blocking input and added a quit command to exit the game.
	scanner := bufio.NewScanner(os.Stdin)

	// Start the timer
	startTime := time.Now()

	for {
		board.PrintBoard(false)
		fmt.Println("Coordinates are a 1-based index. (1, 1) is the top-left corner.")
		fmt.Println("Enter your move in the format 'cmd x y' (cmd: reveal, flag), or type 'quit' to exit:")

		scanner.Scan()
		input := scanner.Text()
		parts := strings.Fields(input) // Split input based on whitespace

		// Ensure we have some input
		if len(parts) == 0 {
			continue
		}

		// Extract the command and coordinates
		cmd := parts[0]

		if cmd == CmdQuit {
			board.PrintBoard(true)
			fmt.Println("Quit game.")
			goto End
		}

		// Ensure we have the correct number of arguments
		if len(parts) != 3 {
			fmt.Println("Invalid input. Please enter a command followed by two integers.")
			continue
		}

		// Check errX and errY separately immediately after conversion
		// Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
		// The bitSize argument specifies the integer type that the result must fit into.
		// In Go, it is idiomatic to check the error as soon as possible, and in this case we should check each conversion separately as strconv.Atoi() is a function call that can return an error.
		x, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid x coordinate. Please enter an integer.")
			continue
		}
		y, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Invalid y coordinate. Please enter an integer.")
			continue
		}

		// Assuming user input is a 1-based index as this is a bit more intuitive for the user
		// We convert the input to a 0-based index to match the array
		x -= 1
		y -= 1

		if !board.isValidCell(x, y) {
			fmt.Println("Invalid coordinates. Please try again.")
			continue
		}

		// Switch case for our commands
		switch cmd {
		case CmdReveal:
			if board.RevealCell(x, y) {
				board.PrintBoard(true)
				fmt.Println("You hit a mine! Game over!")
				goto End
			}
			if board.CheckWin() {
				board.PrintBoard(true)
				fmt.Println("Congratulations, you won!")
				goto End
			}
		case CmdFlag:
			board.FlagCell(x, y)
		default:
			fmt.Println("Invalid command. Please use 'reveal' or 'flag'.")
		}
	}

	// End the timer
End:
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Game duration: %.2f seconds\n", duration.Seconds())
}
