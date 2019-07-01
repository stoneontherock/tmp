package httpapi

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type adaIn struct {
	InitialDomain string   `json:"initialDomain" binding:"required"` //数字或字母或点号
	JoinedDomain  []string `json:"joinedDomain"`
	Pstr          string   `json:"pstr" binding:"required"`
}

const rule = "域只能由小写字母、数字、点号、下划线组成，且不能是全部是点号或下划线"

func AddDomainAdmin(c *gin.Context) {
	var ai adaIn
	err := c.BindJSON(&ai)
	if err != nil {
		respErr(c, 400, err.Error())
		return
	}

	var dmList = []string{ai.InitialDomain}
	dmList = append(dmList, ai.JoinedDomain...)
	if !validateDomain(dmList) {
		respErr(c, 400, rule)
		return
	}

	salt := getRandomStr(8)
	var da = DomainAdmin{
		Name:          "admin@" + ai.InitialDomain,
		Pstr:          md5sum(ai.Pstr + salt),
		Salt:          salt,
		InitialDomain: ai.InitialDomain,
		JoinedDomain:  strings.Join(ai.JoinedDomain, ","),
	}
	err = DB.Create(&da).Error
	if err != nil {
		respErr(c, 500, "新增域管理者失败:"+err.Error())
		return
	}

	respOkEmpty(c)
}

func validateDomain(dmList []string) bool {
	for _, dm := range dmList {
		if !strings.ContainsAny(dm, "abcdefghijklmnopqrstuvwxyz0123456789") {
			return false
		}

		for _, r := range dm {
			if (r > 'a' && r < 'z') && (r > '0' && r < '9') && r != '.' && r != '_' {
				return false
			}
		}
	}

	return true
}
