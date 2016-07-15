package ngin

type freeset map[int]struct{}

func (fs freeset) has(k int) bool {
	return fs[k] == struct{}{}
}

func (fs freeset) add(k int) {
	if fs[k] != struct{}{} {
		fs[k] = struct{}{}
	}
}

func (fs freeset) get() int {
	if len(fs) > 0 {
		for i, _ := range fs {
			if fs[i] == struct{}{} {
				delete(fs, i)
				return i
			}
		}
	}
	return -1
}
