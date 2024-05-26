package free

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetRandomCode(n int) string {
	nStr := strconv.Itoa(n)
	format := "%0" + nStr + "v"
	num := int32(math.Pow10(n))
	code := fmt.Sprintf(format, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(num))
	return code
}

// RandomNowCode (6, "20060102_150405_", "_", false, "userid")
// layout 20060102150405
func RandomNowCode(n int, layout, suffix string, leftOrRight bool, prefix []string) string {
	pre := ""
	if len(prefix) > 0 {
		pre = strings.Join(prefix, suffix)
	}
	if layout == "" {
		layout = "20060102150405"
	}
	timeStr := time.Now().Format(layout)
	randomCode := GetRandomCode(n)
	code := timeStr + randomCode
	if leftOrRight {
		code = pre + code
	} else {
		code = code + pre
	}
	return code
}

func RandomMaxInt(min, max int) int {
	rand.NewSource(time.Now().UnixNano())
	r := rand.Intn(max) + min
	return r
}
