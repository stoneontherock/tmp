package httpapi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"strings"
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

func getDomain(c *gin.Context, domain string) (string, error) {
	dm := ctxDomain(c)
	logrus.Debugf("iDom=%s\n", dm)
	if ctxUser(c) == SA {
		dm = domain //超管不使用token中的domain
		if dm == "" {
			return "", fmt.Errorf("超管必须指定domain")
		}
	}
	return dm, nil
}

func ctxUser(c *gin.Context) string {
	return c.GetString("userName")
}
func ctxDomain(c *gin.Context) string {
	return c.GetString("domain")
}

func isName(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	name, ok := field.Interface().(string)
	if !ok {
		return false
	}

	if len(name) < 3 {
		return false
	}

	valid := true
	for _, r := range name {
		//允许小写字母 - _ 数字
		if (r >= 'a' && r <= 'z') || r == '-' || r == '_' || (r >= '0' && r <= '9') { // || unicode.IsOneOf([]*unicode.RangeTable{unicode.Han}, r) {
			continue
		}
		valid = false
	}

	return valid
}

func isPstr(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	name, ok := field.Interface().(string)
	if !ok {
		return false
	}

	if len(name) < 8 {
		return false
	}

	var upper, lower, digit, punct bool
	for _, r := range name {
		//必须是数字字母符号
		if r > '~' || r < '!' {
			return false
		}

		//密码必须包含大写字母
		if !upper && (r >= 'A' && r <= 'Z') {
			upper = true
			continue
		}

		//密码必须包含小写字母
		if !lower && (r >= 'a' && r <= 'z') {
			lower = true
			continue
		}

		//密码必须包含数字
		if !digit && (r >= '0' && r <= '9') {
			digit = true
			continue
		}

		//密码必须包含标点符号
		if !punct && !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
			punct = true
			continue
		}
	}

	println(upper, lower, digit, punct)
	return upper && lower && digit && punct
}

func isDomain(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	dom, ok := field.Interface().(string)
	if !ok {
		return false
	}

	if !strings.ContainsAny(dom, "abcdefghijklmnopqrstuvwxyz0123456789") {
		return false
	}

	for _, r := range dom {
		//域名中的字符必须是: 小写字母/数字/点号/下划线
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '.' && r != '_' {
			return false
		}
	}

	return true
}

func jwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		v := strings.Fields(c.GetHeader("Authorization"))
		if len(v) < 2 || v[0] != "Bearer" {
			respErr(c, 401, "认证失败(header)") //respErr中已经包含c.Abort(),所以这里不需要c.Abort()
			return
		}

		claims, err := ParseToken(v[1])
		if err != nil {
			respErr(c, 401, "token解析失败: "+err.Error())
			return
		}
		if !strings.HasPrefix(claims.Username, "admin@") && claims.Username != SA {
			respErr(c, 401, "普通用户不可以管理权限")
			return
		}

		if err != nil {
			errStr := "认证失败"
			if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
				errStr = "token过期"
			}

			respErr(c, 401, errStr)
			return
		}

		if claims.Domain == "" {
			respErr(c, 401, "空域")
			return
		}
		c.Set("domain", claims.Domain)
		c.Set("userName", claims.Username)

		//超管可以访问任意路径
		if claims.Username == SA {
			c.Next()
			return
		}

		//域管理员只能访问 /api/vN/da下的路径
		if strings.HasPrefix(claims.Username, "admin@") && strings.HasPrefix(c.Request.RequestURI, URI_VER+"/da") {
			c.Next()
			return
		}

		respErr(c, 403, "Request URI forbidden")
		return
	}
}
