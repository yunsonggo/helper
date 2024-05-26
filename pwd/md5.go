package pwd

import (
	"crypto/md5"
	"fmt"
)

func GenMD5(secret, word string) string {
	str := fmt.Sprintf("%s%s", secret, word)
	has := md5.Sum([]byte(str))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
