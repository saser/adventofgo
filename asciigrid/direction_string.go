// Code generated by "stringer -type=Direction"; DO NOT EDIT.

package asciigrid

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Up-1]
	_ = x[Down-2]
	_ = x[Left-3]
	_ = x[Right-4]
	_ = x[TopLeft-5]
	_ = x[TopRight-6]
	_ = x[BottomLeft-7]
	_ = x[BottomRight-8]
}

const _Direction_name = "NoneUpDownLeftRightTopLeftTopRightBottomLeftBottomRight"

var _Direction_index = [...]uint8{0, 4, 6, 10, 14, 19, 26, 34, 44, 55}

func (i Direction) String() string {
	if i < 0 || i >= Direction(len(_Direction_index)-1) {
		return "Direction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Direction_name[_Direction_index[i]:_Direction_index[i+1]]
}