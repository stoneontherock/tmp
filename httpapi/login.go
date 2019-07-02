package httpapi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type loginIn struct {
	UserName string `json:"userName" binding:"required"`
	Pstr     string `json:"pstr" binding:"required"`
}

func Login(c *gin.Context) {
	var li loginIn
	err := c.BindJSON(&li)
	if err != nil {
		respErr(c, 400, "参数验证失败:"+err.Error())
		return
	}

	if strings.HasPrefix(li.UserName, "admin@") || li.UserName == SA {
		var da DomainAdmin
		err = DB.Where(&DomainAdmin{Name: li.UserName}).First(&da).Error
		if err == nil {
			checkUser(c, da.Name, li.Pstr, da.Salt, da.Pstr, da.InitialDomain, da.JoinedDomain)
			return
		}
	} else {
		var user User
		err = DB.Where(&User{Name: li.UserName}).First(&user).Error
		if err == nil {
			checkUser(c, user.Name, li.Pstr, user.Salt, user.Pstr, user.DefaultDomain, "")
			return
		}
	}

	if err == gorm.ErrRecordNotFound {
		respErr(c, 400, "用户不存在")
		return
	}

	respErr(c, 500, err.Error())
	return
}

func checkUser(c *gin.Context, userName, pstr, salt, hashedPstr, initialDomain, joinedDomain string) {
	sum := md5sum(pstr + salt)
	fmt.Printf("******************* pstr=%s salt=%s hashedPstr=%s sum=%s\n", pstr, salt, hashedPstr, sum)
	if sum != hashedPstr {
		c.JSON(401, gin.H{"code": 401, "msg": "用户名或密码错误"})
		return
	}

	token, err := genJWTToken(userName, initialDomain, joinedDomain)
	if err != nil {
		respErr(c, 500, "token生成失败:"+err.Error())
		return
	}

	c.JSON(200, gin.H{"code": 200, "token": token})
}

type Claims struct {
	Username      string `json:"username"`
	InitialDomain string `json:"initialDomain"`
	JoinedDomain  string `json:"joinedDomain"`
	jwt.StandardClaims
}

var JWTSecret = []byte("1234567890")

func genJWTToken(username, iDom, jDom string) (string, error) {
	claims := Claims{
		username,
		iDom,
		jDom,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), //todo 过期时间
			Issuer:    "aa",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTSecret) //todo 配置jwt密码
	return token, err
}

func parseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
