package schedule

import (
	"math/rand"
	"time"
)

type Schedule struct {
	member         []string
	personEveryDay int
	totalDay       int

	_lastPick          map[string]interface{} // 用来检查上次是否有重复
	_memberBucketIndex int
}

func (t *Schedule) Create(member []string, personEveryDay int, totalDay int) []interface{} {
	t.member, t.personEveryDay, t.totalDay = member, personEveryDay, totalDay
	t._lastPick = make(map[string]interface{})
	return t.create()
}

func (t *Schedule) create() []interface{} {

	rand.Seed(time.Now().UnixNano())
	memberBucket := make(map[int][]string)
	memberBucket[0] = t.member

	var result = make([]interface{}, t.totalDay)
	for index := 0; index < t.totalDay; index++ {
		result[index] = t.pick(memberBucket)

	}

	return result
}

func (t *Schedule) pick(memberBucket map[int][]string) []string {
	var (
		result       []string
		_currentPick = make(map[string]interface{}) // 用来检查本次是否有重复
	)

	for i := 0; i < t.personEveryDay; i++ {
		for {
			member := memberBucket[t._memberBucketIndex]
			if len(member) > 0 {
				randomKey := rand.Intn(len(member))

				// 检查是否有两天连续
				if _, ok := t._lastPick[member[randomKey]]; ok {
					continue
				}
				// 检查是否有重复
				if _, ok := _currentPick[member[randomKey]]; ok {
					continue
				}
				// 追加到结果
				result = append(result, member[randomKey])
				// 标记为已选
				_currentPick[member[randomKey]] = nil
				// 从memberBucket中移动到下一个桶
				memberBucket[t._memberBucketIndex+1] = append(memberBucket[t._memberBucketIndex+1], member[randomKey])
				memberBucket[t._memberBucketIndex] = append(memberBucket[t._memberBucketIndex][:randomKey], memberBucket[t._memberBucketIndex][randomKey+1:]...)
				break
			}
			t._memberBucketIndex++
		}
	}

	t._lastPick = make(map[string]interface{})
	for i := range result {
		t._lastPick[result[i]] = nil
	}

	return result
}
