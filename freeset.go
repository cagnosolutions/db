package db

type freeset map[int]bool

func (fs freeset) has(k int) bool {
	return fs[k] == true
}

func (fs freeset) add(k int) {
	if !fs[k] {
		fs[k] = true
	}
}

func (fs freeset) get() int {
	if len(fs) > 0 {
		for i, _ := range fs {
			if fs[i] {
				delete(fs, i)
				return i
			}
		}
	}
	return -1
}
