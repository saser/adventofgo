package striter

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type emptyIter struct{}

func (emptyIter) Next() (string, bool) { return "", false }

type sliceIter struct {
	ss []string
	i  int
}

type constIter struct {
	s       string
	yielded bool
}

func (i *constIter) Next() (string, bool) {
	ok := !i.yielded
	if !i.yielded {
		i.yielded = true
	}
	return i.s, ok
}

func (i *sliceIter) Next() (string, bool) {
	if i.i >= len(i.ss) {
		return "", false
	}
	s := i.ss[i.i]
	i.i++
	return s, true
}

func TestCollect(t *testing.T) {
	for _, test := range []struct {
		name string
		iter Iter
		want []string
	}{
		{
			name: "empty",
			iter: emptyIter{},
			want: nil,
		},
		{
			name: "single",
			iter: &constIter{s: "foo"},
			want: []string{"foo"},
		},
		{
			name: "slice",
			iter: &sliceIter{ss: []string{"foo", "bar", "quux"}},
			want: []string{"foo", "bar", "quux"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := Collect(test.iter)
			if diff := cmp.Diff(test.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Collect() returned unexpected results (-want +got)\n%s", diff)
			}
		})
	}
}
