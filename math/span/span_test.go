package span

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	for _, tt := range []struct {
		start, end int
		want       Span[int]
	}{
		{
			start: 0,
			end:   0,
			want:  Empty[int](),
		},
		{
			start: 1337,
			end:   1337,
			want:  Empty[int](),
		},
		{
			start: 0,
			end:   1337,
			want: Span[int]{
				Start: 0,
				End:   1337,
			},
		},
		{
			start: 1337,
			end:   0,
			want:  Empty[int](),
		},
	} {
		got := New(tt.start, tt.end)
		if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Span[int]{})); diff != "" {
			t.Errorf("New(%d, %d) returned unexpected result (-want +got)\n%s", tt.start, tt.end, diff)
		}
	}
}

func TestIntersection(t *testing.T) {
	for _, tt := range []struct {
		a, b Span[int]
		want Span[int]
	}{
		{
			a:    New(0, 5), // [01234]
			b:    New(0, 3), // [012]
			want: New(0, 3), // [012]
		},
		{
			a:    New(0, 5), // [01234]
			b:    New(3, 7), //    [3456]
			want: New(3, 5), //    [34]
		},
		{
			a:    New(0, 5),    // [01234]
			b:    Empty[int](), // []
			want: Empty[int](), // []
		},
		{
			a:    New(0, 5), // [01234]
			b:    New(1, 4), //  [123]
			want: New(1, 4), //  [123]
		},
	} {
		if got := Intersection(tt.a, tt.b); got != tt.want {
			t.Errorf("Intersection(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
		if got := Intersection(tt.b, tt.a); got != tt.want {
			t.Errorf("Intersection(%v, %v) = %v; want %v", tt.b, tt.a, got, tt.want)
		}
	}
}

func TestUnion(t *testing.T) {
	for _, tt := range []struct {
		a, b Span[int]
		want Span[int]
	}{
		{
			a:    New(0, 5), // [01234]
			b:    New(0, 3), // [012]
			want: New(0, 5), // [01234]
		},
		{
			a:    New(0, 5), // [01234]
			b:    New(3, 7), //    [3456]
			want: New(0, 7), // [0123456]
		},
		{
			a:    New(0, 5),    // [01234]
			b:    Empty[int](), // []
			want: New(0, 5),    // [01234]
		},
		{
			a:    New(0, 5), // [01234]
			b:    New(1, 4), //  [123]
			want: New(0, 5), // [01234]
		},
	} {
		if got := Union(tt.a, tt.b); got != tt.want {
			t.Errorf("Union(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
		}
		if got := Union(tt.b, tt.a); got != tt.want {
			t.Errorf("Union(%v, %v) = %v; want %v", tt.b, tt.a, got, tt.want)
		}
	}
}

func TestSpan_Len(t *testing.T) {
	for _, tt := range []struct {
		s    Span[int]
		want int
	}{
		{
			s:    New(0, 10),
			want: 10,
		},
		{
			s:    New(0, 0),
			want: 0,
		},
		{
			s:    Empty[int](),
			want: 0,
		},
		{
			s:    New(-1, 0),
			want: 1,
		},
	} {
		if got := tt.s.Len(); got != tt.want {
			t.Errorf("(%v).Len() = %v; want %v", tt.s, got, tt.want)
		}
	}
}

func TestSpan_Contains(t *testing.T) {
	for _, tt := range []struct {
		s    Span[int]
		v    int
		want int
	}{
		{
			s:    New(0, 10),
			v:    -1,
			want: -1,
		},
		{
			s:    New(0, 10),
			v:    0,
			want: 0,
		},
		{
			s:    New(0, 10),
			v:    5,
			want: 0,
		},
		{
			s:    New(0, 10),
			v:    10,
			want: +1,
		},
	} {
		sgn := func(i int) int {
			if i < 0 {
				return -1
			}
			if i > 0 {
				return +1
			}
			return 0
		}
		got := tt.s.Contains(tt.v)
		if gotSgn, wantSgn := sgn(got), sgn(tt.want); gotSgn != wantSgn {
			t.Errorf("%v.Contains(%v) = %+d with sgn %+d; want sgn %+d", tt.s, tt.v, got, gotSgn, wantSgn)
		}
	}
}
