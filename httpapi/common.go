package httpapi

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
)

const SA = "super_admin"

func respErr(c *gin.Context, code int, err string) {
	logrus.Errorf("%s", err)
	c.JSON(code, struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{code, err})
	c.Abort()
}

func respOkEmpty(c *gin.Context) {
	c.JSON(200, struct {
		Code int `json:"code"`
	}{200})
	c.Abort()
}

func respOkData(c *gin.Context, data interface{}) {
	c.JSON(200, struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		200,
		data,
	})
	c.Abort()
}

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

//todo DEL
//func CasbinTest(c *gin.Context) {
//	b := Enforcer.Enforce(c.Query("n"), c.Query("p"), c.Query("m"))
//	c.JSON(200, b)
//}
