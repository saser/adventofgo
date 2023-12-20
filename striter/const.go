package striter

type Const struct {
	s  string
	ok bool
}

func Of(s string) *Const {
	return &Const{
		s:  s,
		ok: true,
	}
}

func (c *Const) Next() (string, bool) {
	if c.ok {
		c.ok = false
		return c.s, true
	}
	return "", false
}
