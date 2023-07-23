package striter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLines(t *testing.T) {
	for _, test := range []struct {
		s    string
		want []string
	}{
		{
			s:    "",
			want: []string{""},
		},
		{
			s:    "\n",
			want: []string{"", ""},
		},
		{
			s:    "\n\n\n",
			want: []string{"", "", "", ""},
		},
		{
			s:    "foo",
			want: []string{"foo"},
		},
		{
			s:    "foo\nbar",
			want: []string{"foo", "bar"},
		},
		{
			s:    "foo\nbar\n",
			want: []string{"foo", "bar", ""},
		},
		{
			s:    "foo\n\nbar",
			want: []string{"foo", "", "bar"},
		},
	} {
		got := Collect(OverLines(test.s))
		if diff := cmp.Diff(test.want, got, cmpopts.EquateEmpty()); diff != "" {
			t.Errorf("Collect(Lines(%q)) returned unexpected result (-want +got)\n%s", test.s, diff)
		}
	}
}
