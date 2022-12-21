package schedule

import "errors"

type Member struct {
	name string
	times int
}

func (t *Member) GetOneTime() (string, error) {
	if t.times > 0 {
		t.times--
		return t.name, nil
	} else {
		return "", errors.New("已经没有了")
	}
}

func (t *Member) GetTimes() int {
	return t.times
}