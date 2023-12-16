// Package asciigrid provides a convenient way to work with 2D grids of ASCII
// characters, which are common in Advent of Code puzzles.
package asciigrid

import (
	"fmt"
	"strings"
)

// Grid is a 2D grid of ASCII characters.
type Grid struct {
	s string

	nRows, nCols int
}

// New parses the given string into a Grid. New assumes that the number of
// columns in the first row determines the number of columns in all rows, and
// returns an error if a row with a different number of columns is found.
func New(s string) (*Grid, error) {
	g := &Grid{s: strings.TrimSpace(s)}
	if len(g.s) == 0 {
		// Special case: this is a completely empty grid, which is valid.
		// Setting these fields to 0 is redundant, as their values already are
		// 0, but it helps readability a bit.
		g.nRows = 0
		g.nCols = 0
		return g, nil
	}
	g.nCols = strings.IndexByte(g.s, '\n')
	if g.nCols == -1 {
		// Special case: this is a grid with a single line, which is valid.
		g.nRows = 1
		g.nCols = len(g.s)
		return g, nil
	}
	for i := 0; i < len(g.s); i += g.nCols + 1 {
		rest := g.s[i:]
		newline := strings.IndexByte(rest, '\n')
		if newline == -1 {
			newline = len(rest)
		}
		if rowLen := len(rest[:newline]); rowLen != g.nCols {
			return nil, fmt.Errorf("asciigrid: row %d has %d columns, expected %d", g.nRows, rowLen, g.nCols)
		}
		g.nRows++
	}
	return g, nil
}

// MustNew is like New but panics on error.
func MustNew(s string) *Grid {
	g, err := New(s)
	if err != nil {
		panic(err)
	}
	return g
}

// NRows is the number of rows in the grid.
func (g *Grid) NRows() int {
	return g.nRows
}

// NCols is the number of columns in the grid.
func (g *Grid) NCols() int {
	return g.nCols
}

// Get returns the ASCII character at the given row and columns (0-indexed) in
// the grid. Get panics if either row or col is out of bounds.
func (g *Grid) Get(row, col int) byte {
	if row < 0 || row >= g.nRows || col < 0 || col >= g.nCols {
		panic(fmt.Errorf("asciigrid: Get(row = %d, col = %d) is out of bounds for grid with %d rows and %d cols", row, col, g.nRows, g.nCols))
	}
	// For each full row, skip over all the columns plus the newline at the end.
	// Then, skip to the right column.
	i := row*(g.nCols+1) + col
	return g.s[i]
}

// Iter represents an iterator over bytes in a string. It is intended to be very
// similar to the proposed Iter[E any] interface from
// https://github.com/golang/go/discussions/54245.
type Iter interface {
	// Next returns the next byte in the iteration if there is one, and reports
	// whether the returned value is valid. Once Next returns ok==false, the
	// iteration is over, and all subsequent calls will return ok==false.
	Next() (b byte, ok bool)
}

// RowIter is an Iter over a single row in the grid. It iterates from left to
// right.
type RowIter struct {
	// row is the row number from which the next byte will be returned from
	// Next.
	row int
	// col is the column number in the row the next byte will be returned. If
	// col >= g.NCols() then all subsequent calls to Next return ok==false.
	col int
	// g is the grid from which bytes will be returned.
	g *Grid
}

var _ Iter = (*RowIter)(nil)

func (i *RowIter) Next() (byte, bool) {
	if i.col >= i.g.NCols() {
		return 0, false
	}
	b := i.g.Get(i.row, i.col)
	i.col++
	return b, true
}

// Row returns an iterator over the given row. Row panics if row is out of bounds.
func (g *Grid) Row(row int) *RowIter {
	if row < 0 || row > g.NRows() {
		panic(fmt.Errorf("asciigrid: Row(%d) is out of bounds in a grid with %d row", row, g.nRows))
	}
	return &RowIter{
		row: row,
		col: 0,
		g:   g,
	}
}

// ColIter is an Iter over a column in the grid. It iterates from top to bottom.
type ColIter struct {
	// row is the row number from which the next byte will be returned from
	// Next. If row >= g.NRows() then all subsequent calls to Next return
	// ok==false.
	row int
	// col is the column number in the row from which the next byte will be
	// returned.
	col int
	// g is the grid from which bytes will be returned.
	g *Grid
}

var _ Iter = (*ColIter)(nil)

func (i *ColIter) Next() (byte, bool) {
	if i.row >= i.g.NRows() {
		return 0, false
	}
	b := i.g.Get(i.row, i.col)
	i.row++
	return b, true
}

// Col returns an iterator over the given column. Col panics if col is out of
// bounds.
func (g *Grid) Col(col int) *ColIter {
	if col < 0 || col >= g.NCols() {
		panic(fmt.Errorf("asciigrid: Col(%d) is out of bounds in a grid with %d columns", col, g.NCols()))
	}
	return &ColIter{
		row: 0,
		col: col,
		g:   g,
	}
}
