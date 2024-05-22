package AccAuthPool

import (
	"free-gpt3.5-2api/common"
	"sync"
)

var (
	instance *AccAuthPool
	once     sync.Once
)

const AccAuthAuthorizationPre = "ac"

type AccAuthPool struct {
	AccAuths []*AccAuth
	index    int
}

type AccAuth struct {
	Authorization string
	CanUseAt      int64
}

func newAccAuthPool() *AccAuthPool {
	return &AccAuthPool{
		AccAuths: make([]*AccAuth, 0),
		index:    0,
	}
}

func GetAccAuthPoolInstance() *AccAuthPool {
	once.Do(func() {
		instance = newAccAuthPool()
	})
	return instance
}

func (a *AccAuthPool) AddAccAuth(accAuth string) {
	a.AccAuths = append(a.AccAuths, &AccAuth{
		Authorization: accAuth,
		CanUseAt:      0,
	})
}

func (a *AccAuthPool) AppendAccAuths(accAuths []string) {
	for _, accAuth := range accAuths {
		a.AddAccAuth(accAuth)
	}
}

func (a *AccAuthPool) Size() int {
	return len(a.AccAuths)
}

func (a *AccAuthPool) IsEmpty() bool {
	return a.Size() == 0
}

func (a *AccAuthPool) CanUseSize() int {
	var count int
	for _, v := range a.AccAuths {
		if v.CanUseAt <= common.GetTimestampSecond(0) {
			count++
		}
	}
	return count
}

func (a *AccAuthPool) GetAccAuth() string {
	if a.IsEmpty() || a.CanUseSize() == 0 {
		return ""
	}
	a.index = (a.index + 1) % len(a.AccAuths)
	authorization := a.AccAuths[a.index].Authorization
	if a.AccAuths[a.index].CanUseAt > common.GetTimestampSecond(0) {
		return a.GetAccAuth()
	}
	return authorization
}

func (a *AccAuthPool) SetCanUseAt(accAuth string, canUseAt int64) {
	for _, v := range a.AccAuths {
		if v.Authorization == accAuth {
			v.CanUseAt = canUseAt
			break
		}
	}
}
