package schedule

import (
	"math/rand"
)

type Schedule struct {
	member         []string
	personEveryDay int
	totalDay       int
}

func (t *Schedule) Create(member []string, personEveryDay int, totalDay int) []interface{} {
	t.member, t.personEveryDay, t.totalDay = member, personEveryDay, totalDay
	return t.create()
}

func (t *Schedule) create() []interface{} {
	memberNum := len(t.member)                  // 成员数量
	personNums := t.totalDay * t.personEveryDay // 总人次
	timesEveryPerson := personNums / memberNum  //每人总执班天数

	memberBucket := make([]Member, memberNum)
	for i, name := range t.member {
		memberBucket[i] = Member{name, timesEveryPerson}
	}

	var result = make([]interface{}, t.totalDay)
	for index := 0; index < t.totalDay; index++ {
		result[index] = t.pick(memberBucket)

	}

	return result
}

// IntersectArray 求两个切片的交集
func IntersectArray(a []string, b []string) []string {
	var inter []string
	mp := make(map[string]bool)

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}

func (t *Schedule) pick(memberBucket []Member) []string {
	var pickedNum int   // 已提取数量
	var pickedKey []int // 已提取过的member key
	var result []string // 提取结果

	var maxkey int    // 剩余次数最大的那个
	var maxkeys []int // 剩余次数最大的可能不止一个

	// 每次最少会挑选一个成员
	for i := 0; i < t.personEveryDay; i++ {
		if pickedNum == t.personEveryDay {
			break
		} else {
			// 次数剩余最大的key
			for key, member := range memberBucket {
				if member.GetTimes() > memberBucket[maxkey].GetTimes() {
					maxkey = key
				}
			}

			// 次数剩余最大的key是不是有多个
			for key, member := range memberBucket {
				if member.GetTimes() == memberBucket[maxkey].GetTimes() && member.GetTimes() > 0 && notPicked(pickedKey, key) {
					maxkeys = append(maxkeys, key)
				}
			}

			// 次数剩余最大的key全部都提取也不会超过本次挑选的限额
			if len(maxkeys) <= t.personEveryDay-pickedNum {
				for _, key := range maxkeys {
					name, err := memberBucket[key].GetOneTime()
					if err != nil {
						continue
					} else {
						result = append(result, name)
						pickedKey = append(pickedKey, key)
						pickedNum++
					}
				}
			} else {
				for pickedNum < t.personEveryDay {
					randomKey := rand.Intn(len(t.member))
					if has(maxkeys, randomKey) && notPicked(pickedKey, randomKey) {
						name, err := memberBucket[randomKey].GetOneTime()
						if err != nil {
							continue
						} else {
							result = append(result, name)
							pickedKey = append(pickedKey, randomKey)
							pickedNum++
						}
					}
				}
			}
		}
	}

	return result
}

func notPicked(source []int, val int) bool {
	flag := true
	for _, v := range source {
		if v == val {
			flag = false
			break
		}
	}
	return flag
}

func has(source []int, val int) bool {
	flag := false
	for _, v := range source {
		if v == val {
			flag = true
			break
		}
	}
	return flag
}
