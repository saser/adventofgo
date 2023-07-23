package geometry

// abs returns the absolute value of x.
func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

// Pos2 represents a point in 2D space with integer coordinates.
type Pos2 struct {
	X, Y int
}

// Add returns a + b.
func (a Pos2) Add(b Pos2) Pos2 {
	return Pos2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

// Sub returns a - b.
func (a Pos2) Sub(b Pos2) Pos2 {
	return Pos2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

// L1Norm returns |p.X| + |p.Y|. This is also known as the Manhattan distance to
// the origin.
func (p Pos2) L1Norm() int {
	return abs(p.X) + abs(p.Y)
}
