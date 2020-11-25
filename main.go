package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var prettyPrint = flag.Bool("p", false, "pretty print sudoku grids")
var getHint = flag.String("hint", "row,col", "get hint only, provide row col")
var validatePuzzle = flag.Bool("validate", false, "validate puzzle only")

func main() {
	flag.Parse()
	if *validatePuzzle {
		validate()
	} else if *getHint != "row,col" {
		hint()
	} else {
		solve()
	}
}

func solve() {
	grid := getPuzzleFromCommandLine()
	if grid == nil {
		return
	}
	solved, err := grid.solve()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "solving sudoku grid:", err)
	} else {
		fmt.Println(solved.toString(*prettyPrint))
	}
}

func validate() {
	grid := getPuzzleFromCommandLine()
	if grid == nil {
		return
	}
	if grid.validate() {
		fmt.Println("valid")
	} else {
		fmt.Println("invalid")
	}
}

func hint() {
	coords := strings.Split(*getHint,",")
	if len(coords) != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "hint requires row and col coordinates in 'r,c' format")
		return
	}
	row, err := strconv.Atoi(coords[0])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "converting row:", err)
	}
	col, err := strconv.Atoi(coords[1])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "converting col:", err)
	}
	if row < 1 || row > 9 {
		_, _ = fmt.Fprintln(os.Stderr, "row must be in range 1..9")
	}
	if col < 1 || col > 9 {
		_, _ = fmt.Fprintln(os.Stderr, "column must be in range 1..9")
	}

	grid := getPuzzleFromCommandLine()
	if grid == nil {
		return
	}
	solved, err := grid.solve()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "solving sudoku grid:", err)
	} else {
		fmt.Println(solved.grid[row-1][col-1])
	}
}

func getPuzzleFromCommandLine() *SudokuGrid {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1000000), 1000000)
	var rows []string
	for i := 0; i < 9; i++ {
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.ReplaceAll(line, " ", "")
			rows = append(rows, line)
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return nil
	}
	g := newSudokuGrid()
	err := g.setValuesFromStringArray(rows)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "constructing sudoku grid:", err)
		return nil
	}
	return &g
}
