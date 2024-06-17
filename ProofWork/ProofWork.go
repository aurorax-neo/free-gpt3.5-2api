package ProofWork

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"time"
)

var (
	numberCollisions = 500000
	cores            = []int{8, 12, 16, 24}
	screens          = []int{3000, 4000, 6000}
	timeLayout       = "Mon Jan 2 2006 15:04:05"
)

type ProofWork struct {
	Difficulty string `json:"difficulty,omitempty"`
	Required   bool   `json:"required"`
	Seed       string `json:"seed,omitempty"`
	Ospt       string `json:"-"`
}

func getParseTime() string {
	now := time.Now()
	return now.Format(timeLayout) + " GMT" + now.Format("-0700 MST (MST)")
}

func getConfig(userAgent string) []interface{} {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	core := cores[rand.Intn(4)]
	rand.New(rand.NewSource(time.Now().UnixNano()))
	screen := screens[rand.Intn(3)]
	return []interface{}{core + screen, getParseTime(), int64(4294705152), 0, userAgent}

}

func CalcProofToken(seed string, diff string, userAgent string) string {
	config := getConfig(userAgent)
	diffLen := len(diff)
	hasher := sha3.New512()
	seedBytes := []byte(seed)
	jsonDataCache := make([]byte, 0, 256)
	for i := 0; i < numberCollisions; i++ {
		config[3] = i
		jsonData, _ := json.Marshal(config)
		jsonDataCache = jsonDataCache[:0]
		jsonDataCache = append(jsonDataCache, jsonData...)
		base := base64.StdEncoding.EncodeToString(jsonDataCache)
		hasher.Write(seedBytes)
		hasher.Write([]byte(base))
		hash := hasher.Sum(nil)
		hasher.Reset()
		if hex.EncodeToString(hash[:diffLen])[:diffLen] <= diff {
			return "gAAAAAB" + base
		}
	}
	return ""
}
