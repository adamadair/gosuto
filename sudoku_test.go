package main

import (
	"fmt"
	"strings"
	"testing"
)

var validationTests = []struct {
	grid     string
	expected bool
}{
	{`1 2 3 4 5 6 7 8 9
4 5 6 7 8 9 1 2 3
7 8 9 1 2 3 4 5 6
9 1 2 3 4 5 6 7 8
3 4 5 6 7 8 9 1 2
6 7 8 9 1 2 3 4 5
8 9 1 2 3 4 5 6 7
2 3 4 5 6 7 8 9 1
5 6 7 8 9 1 2 3 4`, true},
	{`4 3 5 2 6 9 7 8 1
6 8 2 5 7 1 4 9 3
1 9 7 8 3 4 5 6 2
8 2 6 1 9 5 3 4 7
3 7 4 6 8 2 9 1 5
9 5 1 7 4 3 6 2 8
5 1 9 3 2 6 8 7 4
2 4 8 9 5 7 1 3 6
7 6 3 4 1 8 2 5 9`, true},
	{`4 3 2 2 6 9 7 8 1
6 8 5 5 7 1 4 9 3
1 9 7 8 3 4 5 6 2
8 2 6 1 9 5 3 4 7
3 7 4 6 8 2 9 1 5
9 5 1 7 4 3 6 2 8
5 1 9 3 2 6 8 7 4
2 4 8 9 5 7 1 3 6
7 6 3 4 1 8 2 5 9`, false},
	{`4 3 5 2 6 9 7 8 1
6 8 2 5 7 1 4 9 3
1 9 7 8 3 4 5 6 2
8 2 6 1 9 5 3 4 7
3 7 4 6 8 2 9 1 5
1 5 9 7 4 3 6 2 8
5 1 9 3 2 6 8 7 4
2 4 8 9 5 7 1 3 6
7 6 3 4 1 8 2 5 9`, false},
	{`4 3 5 2 6 9 7 8 1
6 8 2 5 7 1 4 9 3
8 9 7 1 3 4 5 6 2
1 2 6 8 9 5 3 4 7
3 7 4 6 8 2 9 1 5
9 5 1 7 4 3 6 2 8
5 1 9 3 2 6 8 7 4
2 4 8 9 5 7 1 3 6
7 6 3 4 1 8 2 5 9`, false},
	{`5 9 6 1 4 2 5 3 7
6 1 4 3 5 8 2 4 8
5 6 9 4 1 2 5 3 6
1 9 5 3 6 8 4 1 6
5 9 3 6 3 4 8 2 1
5 9 5 3 2 1 4 5 6
1 3 6 4 8 6 5 2 5
4 1 2 3 6 8 4 9 2
3 6 8 7 4 1 5 6 3`, false},
	{`1 3 3 4 5 6 7 7 9
4 5 6 7 7 9 1 3 3
7 7 9 1 3 3 4 5 6
9 1 3 3 4 5 6 7 7
3 4 5 6 7 7 9 1 3
6 7 7 9 1 3 3 4 5
7 9 1 3 3 4 5 6 7
3 3 4 5 6 7 7 9 1
5 6 7 7 9 1 3 3 4`, false},
}

func TestValidation(t *testing.T) {
	for i, v := range validationTests {
		g := newSudokuGrid()
		e := g.setValuesFromString(v.grid)
		if e != nil {
			t.Errorf("error setting grid values for test case #%d: %s", i, e.Error())
			return
		}
		if b := g.validate(); b != v.expected {
			t.Errorf("test case #%d, expected %t, received %t", i, v.expected, b)
		} else {
			fmt.Printf("Successfully validated:\n%s\n\n", g.toString(true))
		}
	}
}

var solveTests = []struct {
	grid     string
	expected string
}{
	{
		`120070560
507932080
000001000
010240050
308000402
070085010
000700000
080423701
034010028`,
		`123874569
567932184
849651237
916247853
358196472
472385916
291768345
685423791
734519628`},
	{`000700040
020801900
000000173
102006097
600090001
970100405
354000000
008604030
010003000`, `531769248
427831956
869425173
182546397
645397821
973182465
354278619
798614532
216953784`},
	{`006000050
003700000
700035008
000070012
000942000
620080000
900120003
000003600
050000700`, `816294357
543718269
792635148
438576912
175942836
629381475
964127583
287453691
351869724`},
	{`800000000
003600000
070090200
050007000
000045700
000100030
001000068
008500010
090000400`, `812753649
943682175
675491283
154237896
369845721
287169534
521974368
438526917
796318452`},
}

func TestSolve(t *testing.T) {
	for i, test := range solveTests {
		g := newSudokuGrid()
		e := g.setValuesFromString(test.grid)
		if e != nil {
			t.Errorf("error setting grid values for test case #%d: %s", i, e.Error())
			return
		}
		if !g.validate() {
			t.Errorf("invalid grid values for test case #%d: %s", i, g.toString(false))
		}

		solved, e := g.solve()
		if e != nil {
			t.Errorf("error solving grid for test case #%d: %s", i, e.Error())
			return
		} else {
			solveString := solved.toString(false)
			if strings.TrimSpace(solveString) != test.expected {
				t.Errorf("testcase #%d failed. expected:\n%s\nreceived:\n%s", i, test.expected, solveString)
			}
		}
	}
}

func TestCellToGrid(t *testing.T) {
	if cellToSubGrid(1, 1) != 0 {
		t.Errorf("cell 0 error")
	}

	if cellToSubGrid(0, 3) != 1 {
		t.Errorf("cell 1 error")
	}

	if cellToSubGrid(2, 8) != 2 {
		t.Errorf("cell 2 error")
	}

	if cellToSubGrid(3, 1) != 3 {
		t.Errorf("cell 3 error")
	}

	if cellToSubGrid(3, 3) != 4 {
		t.Errorf("cell 4 error")
	}

	if cellToSubGrid(5, 8) != 5 {
		t.Errorf("cell 2 error")
	}
}
