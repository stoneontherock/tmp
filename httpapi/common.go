package httpapi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
)


//ascii 33~126: 字母数字符号
func getRandomStr(n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = byte(rand.Intn(94)) + 33
	}

	return string(s)
}

func md5sum(key string) string {
	w := md5.New()
	io.WriteString(w, key)
	return hex.EncodeToString(w.Sum(nil))
}


func getDomain(c *gin.Context,domain string) (string,error) {
	dm := ctxIDom(c)
	fmt.Printf("****************iDom=%s\n",dm)
	if ctxUser(c) == SA {
		dm = domain
		if dm == "" {
			return "",fmt.Errorf("超管必须指定domain")
		}
	}
	return dm,nil
}
