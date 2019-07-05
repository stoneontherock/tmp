package httpapi

import (
	"aa/config"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type loginIn struct {
	UserName string `json:"userName" binding:"required"`
	Pstr     string `json:"pstr" binding:"required"`
}

var ErrGenToken = errors.New("token生成失败")

func login(c *gin.Context) {
	var li loginIn
	err := c.ShouldBindJSON(&li)
	if err != nil {
		respErr(c, 400, "参数验证失败:"+err.Error())
		return
	}

	token, err := CheckUser(li.UserName, li.Pstr)
	if err == nil {
		c.JSON(200, gin.H{"code": 200, "token": token})
		return
	}

	if err == ErrGenToken {
		respErr(c, 500, err.Error())
		return
	}

	respErr(c, 400, err.Error())
	return
}

func CheckUser(userName, pstr string) (string, error) {
	var user User
	err := DB.Where(&User{Name: userName}).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("用户%s不存在\n", userName)
		}
		return "", err
	}

	sum := md5sum(pstr + user.Salt)
	logrus.Debugf("pstr=%s salt=%s hashedPstr=%s sum=%s\n", pstr, user.Salt, user.Pstr, sum)
	if sum != user.Pstr {
		logrus.Warnf("token生成失败,%v", err)
		return "", fmt.Errorf("用户名或密码错误\n")
	}

	token, err := genJWTToken(userName, user.Domain)
	if err != nil {
		logrus.Errorf("token生成失败,%v", err)
		return "", ErrGenToken
	}

	return token, nil
}

type Claims struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	jwt.StandardClaims
}

func genJWTToken(username, domain string) (string, error) {
	claims := Claims{
		username,
		domain,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * config.C.JWT.Expire).Unix(),
			Issuer:    "aa",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) { return jwtSecret, nil })

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
