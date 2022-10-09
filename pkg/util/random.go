package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZwqertyhbvcsqrhybvcdfehjhnvgbvcfewwgtdvcqefrgfbvcefrthgbdvsdcfwgrfvcxfeqwdscxefrshg"

func inti() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt returns a random int64 in [min, max]
func RandInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

//RandString generates a random string of length n
func RandString(n int) string {
	var sb strings.Builder
	k := len(letterBytes)
	for i := 0; i < n; i++ {
		c := letterBytes[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandOwner generates a random owner
func RandOwner() string {
	return RandString(10)
}

//RandMoney generates a random amount of money
func RandMoney() int64 {
	return RandInt(1, 100)
}

//RandCurrency select a random currency
func RandCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CNY", "AUD", "CAD", "CHF"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandString(4))
}
