package httpapi

import (
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

//type Model struct {
//	Name      string    `json:"name" gorm:"primary_key"`
//	CreatedAt time.Time `json:"createdAt"`
//	UpdatedAt time.Time `json:"updatedAt"`
//}

type Role struct {
	Name      string     `gorm:"type:varchar(64);primary_key" json:"name"`
	Domain    string     `gorm:"type:varchar(32);index" json:"domain"`
	Resources []Resource `gorm:"many2many:role_resource;" json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

type Resource struct {
	ID        uint      `gorm:"primary_key" json:"ID"`
	Name      string    `gorm:"type:varchar(128)" json:"name"`
	Action    string    `gorm:"type:varchar(16)" json:"action"`
	Domain    string    `gorm:"type:varchar(32);index" json:"-"`
	Comment   string    `gorm:"type:varchar(128)" json:"comment"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	Name   string `gorm:"type:varchar(32);primary_key" json:"name"`
	Pstr   string `gorm:"type:varchar(32)" json:"-"`
	Salt   string `gorm:"type:varchar(8)" json:"-"`
	Domain string `gorm:"type:varchar(32);index" json:"domain"`
	//Creator       string    `gorm:"type:varchar(38)" json:"-"`
	Roles     []Role    `gorm:"many2many:user_role;" json:"-" `
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type userRole struct {
	UserName string
	RoleName string
}

//role-resource表
type roleResource struct {
	RoleName string
}

func txCommit(db *gorm.DB, commit *bool) (*gorm.DB, func()) {
	tx := db.Begin()
	return tx, func() {
		if *commit {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}

//func getJoinedDomainByInitialDomain(initialDomain string) (string, error) {
//	var da DomainAdmin
//	err := DB.Where(`initial_domain = ?`, initialDomain).First(&da).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		return "", err
//	}
//	return da.JoinedDomain, nil
//}
