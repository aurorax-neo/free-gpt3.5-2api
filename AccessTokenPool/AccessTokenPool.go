package AccessTokenPool

import (
	"free-gpt3.5-2api/common"
	"sync"
)

var (
	instance *AccessTokenPool
	once     sync.Once
)

const AccAuthAuthorizationPre = "ac"

type AccessTokenPool struct {
	AccessTokens []*AccessToken
	index        int
}

type AccessToken struct {
	Token    string
	CanUseAt int64
}

func newAccAuthPool() *AccessTokenPool {
	return &AccessTokenPool{
		AccessTokens: make([]*AccessToken, 0),
		index:        0,
	}
}

func GetAccAuthPoolInstance() *AccessTokenPool {
	once.Do(func() {
		instance = newAccAuthPool()
	})
	return instance
}

func (a *AccessTokenPool) AddToken(token string) {
	a.AccessTokens = append(a.AccessTokens, &AccessToken{
		Token:    token,
		CanUseAt: 0,
	})
}

func (a *AccessTokenPool) AppendTokens(tokens []string) {
	for _, accAuth := range tokens {
		a.AddToken(accAuth)
	}
}

func (a *AccessTokenPool) Size() int {
	return len(a.AccessTokens)
}

func (a *AccessTokenPool) IsEmpty() bool {
	return a.Size() == 0
}

func (a *AccessTokenPool) CanUseSize() int {
	var count int
	for _, v := range a.AccessTokens {
		if v.CanUseAt <= common.GetTimestampSecond(0) {
			count++
		}
	}
	return count
}

func (a *AccessTokenPool) GetToken() string {
	if a.IsEmpty() || a.CanUseSize() == 0 {
		return ""
	}
	a.index = (a.index + 1) % len(a.AccessTokens)
	authorization := a.AccessTokens[a.index].Token
	if a.AccessTokens[a.index].CanUseAt > common.GetTimestampSecond(0) {
		return a.GetToken()
	}
	return authorization
}

func (a *AccessTokenPool) SetCanUseAt(token string, canUseAt int64) {
	for _, v := range a.AccessTokens {
		if v.Token == token {
			v.CanUseAt = canUseAt
			break
		}
	}
}
