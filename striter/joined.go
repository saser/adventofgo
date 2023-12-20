package striter

type Joined struct {
	iters []Iter
}

func Join(it Iter, its ...Iter) *Joined {
	return &Joined{
		iters: append([]Iter{it}, its...),
	}
}

func (j *Joined) Next() (string, bool) {
	for {
		if len(j.iters) == 0 {
			return "", false
		}
		s, ok := j.iters[0].Next()
		if !ok {
			j.iters = j.iters[1:]
			continue
		}
		return s, true
	}
}
