// Package asciigrid provides a convenient way to work with 2D grids of ASCII
// characters, which are common in Advent of Code puzzles.
package asciigrid

import (
	"bytes"
	"fmt"
	"strconv"
)

// Direction represents the vertical, horizontal, and diagonal directions in the
// grid. A Direction's value is based on the assumption that rows increase from
// top to bottom and columns increase from left to right. To illustrate, a
// direction represents any of the positions marked 'O' from the perspective of
// the position p:
//
//	.....
//	.OOO.
//	.OpO.
//	.OOO.
//	.....
type Direction int

//go:generate go run golang.org/x/tools/cmd/stringer -type=Direction

const (
	// None is the lack of a direction. Taking a step in direction None leaves you in the same place.
	None Direction = iota
	//	.....
	//	..O..
	//	..p..
	//	.....
	//	.....
	Up
	//	.....
	//	.....
	//	..p..
	//	..O..
	//	.....
	Down
	//	.....
	//	.....
	//	.Op..
	//	.....
	//	.....
	Left
	//	.....
	//	.....
	//	..pO.
	//	.....
	//	.....
	Right
	//	.....
	//	.O...
	//	..p..
	//	.....
	//	.....
	TopLeft
	//	.....
	//	...O.
	//	..p..
	//	.....
	//	.....
	TopRight
	//	.....
	//	.....
	//	..p..
	//	.O...
	//	.....
	BottomLeft
	//	.....
	//	.....
	//	..p..
	//	...O.
	//	.....
	BottomRight
)

func (d Direction) String() string {
	switch d {
	case None:
		return "None"
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case TopLeft:
		return "TopLeft"
	case TopRight:
		return "TopRight"
	case BottomLeft:
		return "BottomLeft"
	case BottomRight:
		return "BottomRight"
	default:
		return "Direction(" + strconv.FormatInt(int64(d), 10) + ")"
	}
}

// Inverse returns the inverse direction. It has the property that for a
// position p and direction d:
//
//	p.Step(d).Step(d.Inverse()) == p
func (d Direction) Inverse() Direction {
	switch d {
	case None:
		return None
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case TopLeft:
		return BottomRight
	case TopRight:
		return BottomLeft
	case BottomLeft:
		return TopRight
	case BottomRight:
		return TopLeft
	default:
		panic(fmt.Errorf("invalid direction %d", d))
	}
}

// Pos represents a position in the grid. Rows and columns are 0-indexed. Rows
// increase from top to bottom and columns increase from left to right.
type Pos struct {
	Row, Col int
}

// Step returns the position a single step in the given direction.
func (p Pos) Step(d Direction) Pos {
	return p.StepN(d, 1)
}

// StepN is like Step but takes n steps in the given direction.
func (p Pos) StepN(d Direction, n int) Pos {
	p2 := p
	switch d {
	case None:
		// Nothing happens.
	case Up:
		p2.Row -= n
	case Down:
		p2.Row += n
	case Left:
		p2.Col -= n
	case Right:
		p2.Col += n
	case TopLeft:
		p2.Row -= n
		p2.Col -= n
	case TopRight:
		p2.Row -= n
		p2.Col += n
	case BottomLeft:
		p2.Row += n
		p2.Col -= n
	case BottomRight:
		p2.Row += n
		p2.Col += n
	}
	return p2
}

// Neighbors4 returns the four direct neighbors (marked 'O' below) to the given
// position:
//
//	.....
//	..O..
//	.OpO.
//	..O..
//	.....
//
// It has the property that:
//
//	Neighbors4()[dir] == p.Step(dir)
func (p Pos) Neighbors4() map[Direction]Pos {
	return map[Direction]Pos{
		Up:    {Row: p.Row - 1, Col: p.Col},
		Down:  {Row: p.Row + 1, Col: p.Col},
		Left:  {Row: p.Row, Col: p.Col - 1},
		Right: {Row: p.Row, Col: p.Col + 1},
	}
}

// Neighbors4 returns the eight neighbors (marked 'O' below) to the given
// position:
//
//	.....
//	.OOO.
//	.OpO.
//	.OOO.
//	.....
//
// It has the property that:
//
//	Neighbors8()[dir] == p.Step(dir)
func (p Pos) Neighbors8() map[Direction]Pos {
	return map[Direction]Pos{
		Up:          {Row: p.Row - 1, Col: p.Col},
		Down:        {Row: p.Row + 1, Col: p.Col},
		Left:        {Row: p.Row, Col: p.Col - 1},
		Right:       {Row: p.Row, Col: p.Col + 1},
		TopLeft:     {Row: p.Row - 1, Col: p.Col - 1},
		TopRight:    {Row: p.Row - 1, Col: p.Col + 1},
		BottomLeft:  {Row: p.Row + 1, Col: p.Col - 1},
		BottomRight: {Row: p.Row + 1, Col: p.Col + 1},
	}
}

// Grid is a 2D grid of ASCII characters.
type Grid struct {
	rows [][]byte

	nRows, nCols int
}

// New parses the given string into a Grid. New assumes that the number of
// columns in the first row determines the number of columns in all rows, and
// returns an error if a row with a different number of columns is found.
func New(s string) (*Grid, error) {
	bs := bytes.TrimSpace([]byte(s))
	g := &Grid{}
	if len(bs) == 0 {
		// Special case: this is a completely empty grid, which is valid.
		// Setting these fields to 0 is redundant, as their values already are
		// 0, but it helps readability a bit.
		g.nRows = 0
		g.nCols = 0
		return g, nil
	}
	g.nCols = bytes.IndexByte(bs, '\n')
	if g.nCols == -1 {
		// Special case: this is a grid with a single line, which is valid.
		g.rows = [][]byte{bs}
		g.nRows = 1
		g.nCols = len(bs)
		return g, nil
	}
	for i := 0; i < len(bs); i += g.nCols + 1 {
		rest := bs[i:]
		newline := bytes.IndexByte(rest, '\n')
		if newline == -1 {
			newline = len(rest)
		}
		row := rest[:newline]
		if rowLen := len(row); rowLen != g.nCols {
			return nil, fmt.Errorf("asciigrid: row %d has %d columns, expected %d", g.nRows, rowLen, g.nCols)
		}
		g.rows = append(g.rows, row)
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

// Index represents a position within a grid as a single integer. Use
// (*Grid).Index to construct these values and (*Grid).Pos to convert them back
// to positions.
type Index int

// Index returns the index which corresponds to p. The return value of Index(p) is
// only valid if InBounds(p) == true.
func (g *Grid) Index(p Pos) Index {
	// For each full row, skip over all the columns plus the newline at the end.
	// Then, skip to the right column.
	return Index(p.Row*g.nCols + p.Col)
}

// Pos converts the given index, assumed to be created by g.Index, into the
// corresponding position.
func (g *Grid) Pos(i Index) Pos {
	return Pos{
		Row: int(i) / g.nCols,
		Col: int(i) % g.nCols,
	}
}

// Get returns the ASCII character at the given position in the grid. Get panics
// if p is out of bounds.
func (g *Grid) Get(p Pos) byte {
	return g.rows[p.Row][p.Col]
}

// Set stores the given ASCII character at the given position. Set panics if p
// is out of bounds.
func (g *Grid) Set(p Pos, b byte) {
	g.rows[p.Row][p.Col] = b
}

// InBounds reports whether the given position is valid in the grid.
func (g *Grid) InBounds(p Pos) bool {
	return p.Row >= 0 && p.Row < g.NRows() && p.Col >= 0 && p.Col < g.NCols()
}

// String returns the string (including newlines) that the Grid currently
// represents.
func (g *Grid) String() string {
	return string(bytes.Join(g.rows, []byte{'\n'}))
}

// Iter represents an iterator over bytes in a string. It is intended to be very
// similar to the proposed Iter[E any] interface from
// https://github.com/golang/go/discussions/54245.
type Iter interface {
	// Next returns the next byte in the iteration if there is one, and reports
	// whether the returned value is valid. Once Next returns ok==false, the
	// iteration is over, and all subsequent calls will return ok==false.
	Next() (p Pos, b byte, ok bool)
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

func (i *RowIter) Next() (Pos, byte, bool) {
	if i.col >= i.g.NCols() {
		return Pos{}, 0, false
	}
	p := Pos{Row: i.row, Col: i.col}
	b := i.g.Get(p)
	i.col++
	return p, b, true
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

func (i *ColIter) Next() (Pos, byte, bool) {
	if i.row >= i.g.NRows() {
		return Pos{}, 0, false
	}
	p := Pos{Row: i.row, Col: i.col}
	b := i.g.Get(p)
	i.row++
	return p, b, true
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
