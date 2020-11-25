package main

import (
	"fmt"
	"strings"
)

// we use 0 in the grid to indicate a cell that has not been set yet.
// These would be the blank values in
const UnsetValue = 0

type SudokuGrid struct {
	grid [9][9]int
}

// initialize a new grid, set all values
func newSudokuGrid() SudokuGrid {
	var s SudokuGrid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			s.grid[i][j] = UnsetValue
		}
	}
	return s
}

// set the sudoku grid values from a single string.
func (g *SudokuGrid) setValuesFromString(s string) error {
	s = strings.ReplaceAll(s, " ", "")
	rows := strings.Split(s, "\n")
	return g.setValuesFromStringArray(rows)
}

// set the sudoku grid values from a slice of strings, each slice member
// represents a row of the grid.
func (g *SudokuGrid) setValuesFromStringArray(rows []string) error {
	l := len(rows)
	if l != 9 {
		return fmt.Errorf("soduko grid requires 9 rows of digits")
	}
	for i := 0; i < 9; i++ {
		rowText := rows[i]
		if len(rowText) != 9 {
			return fmt.Errorf("soduko grid row requires 9 characters")
		}
		for j := 0; j < 9; j++ {
			g.grid[i][j] = runeToInt(rune(rowText[j]))
		}
	}
	return nil
}

// return true if there are no duplicate digits in any row, column, or sub-grid.
func (g *SudokuGrid) validate() bool {
	for i := 0; i < 9; i++ {
		if !(g.validateRow(i) && g.validateColumn(i) && g.validateSubGrid(i)) {
			return false
		}
	}
	return true
}

// returns true if row has no duplicate values
// r must be 0..8
func (g *SudokuGrid) validateRow(r int) bool {
	pos := [9]bool{false, false, false, false, false, false, false, false, false}
	for i := 0; i < 9; i++ {
		v := g.grid[r][i]
		if v >= 1 && v <= 9 {
			if pos[v-1] == true {
				return false
			}
			pos[v-1] = true
		}
	}
	return true
}

// return true if column has no duplicates
// c must be 0..8
func (g *SudokuGrid) validateColumn(c int) bool {
	pos := [9]bool{false, false, false, false, false, false, false, false, false}
	for i := 0; i < 9; i++ {
		v := g.grid[i][c]
		if v >= 1 && v <= 9 {
			if pos[v-1] == true {
				return false
			}
			pos[v-1] = true
		}
	}
	return true
}

// return true if sub-grid has no duplicates
// s must be 0..9
func (g *SudokuGrid) validateSubGrid(s int) bool {
	pos := [9]bool{false, false, false, false, false, false, false, false, false}
	row := s / 3
	col := s % 3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			v := g.grid[3*row+r][3*col+c]
			if v >= 1 && v <= 9 {
				if pos[v-1] == true {
					return false
				}
				pos[v-1] = true
			}
		}
	}
	return true
}

// get string representation of grid. set makePretty to true for printing to screen.
// set makePretty to false for file or clipboard.
func (g *SudokuGrid) toString(makePretty bool) string {
	var sb strings.Builder
	if makePretty {
		sb.WriteString("+---+---+---+\n")
	}
	for i := 0; i < 9; i++ {
		if makePretty {
			sb.WriteString(fmt.Sprintf("|%d%d%d|%d%d%d|%d%d%d|\n", g.grid[i][0], g.grid[i][1], g.grid[i][2], g.grid[i][3], g.grid[i][4], g.grid[i][5], g.grid[i][6], g.grid[i][7], g.grid[i][8]))
			if (i+1)%3 == 0 {
				sb.WriteString("+---+---+---+\n")
			}
		} else {
			sb.WriteString(fmt.Sprintf("%d%d%d%d%d%d%d%d%d\n", g.grid[i][0], g.grid[i][1], g.grid[i][2], g.grid[i][3], g.grid[i][4], g.grid[i][5], g.grid[i][6], g.grid[i][7], g.grid[i][8]))
		}
	}
	return sb.String()
}

func (g *SudokuGrid) copy() SudokuGrid {
	r := newSudokuGrid()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			r.grid[i][j] = g.grid[i][j]
		}
	}
	return r
}

// solve the grid, returns a pointer to a new SudokuGrid
func (g *SudokuGrid) solve() (*SudokuGrid, error) {
	r := g.copy()
	if !r.fill(0) {
		return nil, fmt.Errorf("grid was unsolvable")
	}
	return &r, nil
}

func (g *SudokuGrid) fill(n int) bool {
	if n >= 81 {
		return true
	}
	row := n / 9
	col := n % 9
	v := g.grid[row][col]
	if v >= 1 && v <= 9 {
		return g.fill(n + 1)
	}
	opts := g.getCellOptions(row, col)
	if len(opts) == 0 {
		return false
	}
	for _, v := range opts {
		g.grid[row][col] = v
		if g.fill(n + 1) {
			return true
		}
		g.grid[row][col] = UnsetValue
	}
	return false
}

func (g *SudokuGrid) getCellOptions(r, c int) []int {
	var response []int

	row := g.getRowChecklist(r)
	col := g.getColumnChecklist(c)
	sub := g.getSubGridChecklist(cellToSubGrid(r, c))
	for i := 0; i < 9; i++ {
		if !row[i] && !col[i] && !sub[i] {
			response = append(response, i+1)
		}
	}
	return response
}

func (g *SudokuGrid) getRowChecklist(n int) []bool {
	r := []bool{false, false, false, false, false, false, false, false, false}
	for i := 0; i < 9; i++ {
		v := g.grid[n][i]
		if v >= 1 && v <= 9 {
			r[v-1] = true
		}
	}
	return r
}

func (g *SudokuGrid) getColumnChecklist(n int) []bool {
	r := []bool{false, false, false, false, false, false, false, false, false}
	for i := 0; i < 9; i++ {
		v := g.grid[i][n]
		if v >= 1 && v <= 9 {
			r[v-1] = true
		}
	}
	return r
}

func (g *SudokuGrid) getSubGridChecklist(n int) []bool {
	ret := []bool{false, false, false, false, false, false, false, false, false}
	row := n / 3
	col := n % 3
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			v := g.grid[3*row+r][3*col+c]
			if v >= 1 && v <= 9 {
				ret[v-1] = true
			}
		}
	}
	return ret
}

func runeToInt(r rune) int {
	if r >= '1' && r <= '9' {
		return int(r - '0')
	}
	return UnsetValue
}

// determine the sub-grid for provided grid coordinates
// r, c => integer in the range 0..8
func cellToSubGrid(r, c int) int {
	return (r/3)*3 + c/3
}
