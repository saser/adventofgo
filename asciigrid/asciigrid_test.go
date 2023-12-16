package asciigrid

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func newT(t *testing.T, s string) *Grid {
	t.Helper()
	g, err := New(s)
	if err != nil {
		t.Fatalf("New(%q) failed unexpectedly: %v", s, err)
	}
	return g
}

func collect(i Iter) string {
	var bs []byte
	for b, ok := i.Next(); ok; b, ok = i.Next() {
		bs = append(bs, b)
	}
	return string(bs)
}

func TestNew(t *testing.T) {
	for _, s := range []string{
		"",
		"\n",
		"\n\n\n\n\n\n\n\n\n\n\n",
		"#",
		".",
		"######",
		"######\n######",
		"######\n######\n\n\n\n\n\n",
	} {
		if _, err := New(s); err != nil {
			t.Errorf("New(%q) err = %v; want nil", s, err)
		}
	}
}

func TestNew_Error(t *testing.T) {
	for _, s := range []string{
		"#\n#\n##",
		"#\n#\n##\n",
		"#\n\n#",
		"#\n\n#\n",
	} {
		if _, err := New(s); err == nil {
			t.Errorf("New(%q) succeeded unexpectedly", s)
		}
	}
}

func TestGrid_NRows_NCols(t *testing.T) {
	for _, tt := range []struct {
		s                  string
		wantRows, wantCols int
	}{
		{
			s:        "",
			wantRows: 0,
			wantCols: 0,
		},
		{
			s:        "\n\n\n\n\n\n",
			wantRows: 0,
			wantCols: 0,
		},
		{
			s:        "#\n\n\n\n\n\n",
			wantRows: 1,
			wantCols: 1,
		},
		{
			s:        "#",
			wantRows: 1,
			wantCols: 1,
		},
		{
			s:        "#####",
			wantRows: 1,
			wantCols: 5,
		},
		{
			s:        "#\n#\n#",
			wantRows: 3,
			wantCols: 1,
		},
		{
			s:        "##\n##\n##",
			wantRows: 3,
			wantCols: 2,
		},
	} {
		g := newT(t, tt.s)
		if got, want := g.NRows(), tt.wantRows; got != want {
			t.Errorf("New(%q).NRows() = %v; want %v", tt.s, got, want)
		}
		if got, want := g.NCols(), tt.wantCols; got != want {
			t.Errorf("New(%q).NCols() = %v; want %v", tt.s, got, want)
		}
	}
}

func TestGrid_Get(t *testing.T) {
	s := strings.TrimSpace(`
abcde
12345
!@#$%
`)
	g := newT(t, s)
	// We use a slice of strings instead of slice of bytes because that makes
	// for error messages that are easier to read -- bytes are often printed
	// with their numerical values, rather than the character they represent.
	var got []string
	for row := 0; row < g.NRows(); row++ {
		for col := 0; col < g.NCols(); col++ {
			got = append(got, string([]byte{g.Get(row, col)}))
		}
	}
	want := []string{
		"a", "b", "c", "d", "e",
		"1", "2", "3", "4", "5",
		"!", "@", "#", "$", "%",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("New(%q): unexpected result from calling Get() on each row and column (-want +got)\n%s", s, diff)
	}
}

func TestGrid_Row(t *testing.T) {
	for _, tt := range []struct {
		s    string
		want []string
	}{
		{
			s:    "",
			want: nil,
		},
		{
			s:    "#",
			want: []string{"#"},
		},
		{
			s:    "##############\n\n\n\n\n\n\n",
			want: []string{"##############"},
		},
		{
			s: strings.TrimSpace(`
abc
def
ghi
`),
			want: []string{"abc", "def", "ghi"},
		},
	} {
		g := newT(t, tt.s)
		var got []string
		for row := 0; row < g.NRows(); row++ {
			got = append(got, collect(g.Row(row)))
		}
		if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("New(%q): unexpected result from iterating over all Row() iterators (-want +got)\n%s", tt.s, diff)
		}
	}
}

func TestGrid_Col(t *testing.T) {
	for _, tt := range []struct {
		s    string
		want []string
	}{
		{
			s:    "",
			want: nil,
		},
		{
			s:    "#",
			want: []string{"#"},
		},
		{
			s:    "#####\n\n\n\n\n\n\n",
			want: []string{"#", "#", "#", "#", "#"},
		},
		{
			s: strings.TrimSpace(`
abc
def
ghi
`),
			want: []string{"adg", "beh", "cfi"},
		},
	} {
		g := newT(t, tt.s)
		var got []string
		for col := 0; col < g.NCols(); col++ {
			got = append(got, collect(g.Col(col)))
		}
		if diff := cmp.Diff(tt.want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("New(%q): unexpected result from iterating over all Row() iterators (-want +got)\n%s", tt.s, diff)
		}
	}
}
