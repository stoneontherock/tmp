package httpapi

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"strings"
)

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

//func isAction(
//	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
//	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
//) bool {
//	action, ok := field.Interface().(string)
//	if !ok {
//		return false
//	}
//
//	return action == "add" || action == "del"
//}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		v := strings.Fields(c.GetHeader("Authorization"))
		if len(v) < 2 || v[0] != "Bearer" {
			respErr(c, 401, "认证失败(header)")
			c.Abort()
			return
		}

		claims, err := parseToken(v[1])
		if err == nil {
			if claims.InitDom == "" {
				respErr(c, 401, "空初始域")
				c.Abort()
				return
			}
			c.Set("initDomain", claims.InitDom)
			if claims.Username == SA {
				c.Set("initDomain", "root")
			}

			if strings.HasPrefix(claims.Username, "admin@") && !strings.HasPrefix(c.Request.RequestURI, URI_VER+"/da") {
				respErr(c, 403, "uri forbidden")
				c.Abort()
				return
			}

			if strings.HasPrefix(claims.Username, SA) && !strings.HasPrefix(c.Request.RequestURI, URI_VER+"/sa") {
				respErr(c, 403, "uri forbidden")
				c.Abort()
				return
			}

			c.Next()
			return
		}

		errStr := "认证失败"
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			errStr = "token过期"
		}

		respErr(c, 401, errStr)
		c.Abort()
	}
}
