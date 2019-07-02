package httpapi

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"os"
	"time"
)

var DB *gorm.DB

//type Model struct {
//	Name      string    `json:"name" gorm:"primary_key"`
//	CreatedAt time.Time `json:"createdAt"`
//	UpdatedAt time.Time `json:"updatedAt"`
//}

func InitModel(db *gorm.DB) {
	//todo: 去掉debug
	DB = db.Debug()
	autoMigrate()
	createSuperAdmin()
	initEnforcer()
}

func autoMigrate() {
	DB.AutoMigrate(&Resource{}, &Role{}, &User{}, &DomainAdmin{})
}

func createSuperAdmin() {
	sa := DomainAdmin{Name: "super_admin"}
	err := DB.First(&sa).Error
	if err == nil {
		fmt.Printf("super_admin: %v\n", sa)
		return
	}

	if err != gorm.ErrRecordNotFound {
		panic(fmt.Sprintf("查询超管记录时出现错误:%v", err))
	}

	pstr := os.Getenv("sa_pstr")
	if pstr == "" {
		pstr = "SuperAdmin@123"
	}

	sa.Salt = getRandomStr(8)
	sa.Pstr = md5sum(pstr + sa.Salt)
	sa.InitialDomain = "root"

	err = DB.Create(&sa).Error
	if err != nil {
		panic(fmt.Sprintf("创建超级管理员账户失败:%v", err))
	}
}

var Enforcer *casbin.Enforcer

func initEnforcer() {
	Enforcer = casbin.NewEnforcer("model.conf", false)
	err := loadAllRoleRescourcePolicy()
	if err != nil {
		panic(err)
	}
	err = loadAllUserRolePolicy()
	if err != nil {
		panic(err)
	}
}

type Role struct {
	Name          string     `gorm:"type:varchar(64);primary_key" json:"name"`
	DefaultDomain string     `gorm:"type:varchar(32)" json:"default_domain"`
	Resources     []Resource `gorm:"many2many:role_resource;" json:"resources"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type Resource struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(128)" json:"name"`
	Act       string    `gorm:"type:varchar(16)" json:"act"`
	Domain    string    `gorm:"type:varchar(32)" json:"domain"`
	Comment   string    `gorm:"type:varchar(128)" json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	Name          string    `gorm:"type:varchar(16);primary_key" json:"name"`
	Pstr          string    `gorm:"type:varchar(32)" json:"-"`
	Salt          string    `gorm:"type:varchar(8)" json:"-"`
	DefaultDomain string    `gorm:"type:varchar(32)" json:"-"`
	//Creator       string    `gorm:"type:varchar(38)" json:"-"`
	Roles         []Role    `gorm:"many2many:user_role;" json:"-" `
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type DomainAdmin struct {
	Name          string    `gorm:"type:varchar(38);primary_key" json:"name"`
	Pstr          string    `gorm:"type:varchar(32)" json:"-"`
	Salt          string    `gorm:"type:varchar(8)" json:"-"`
	InitialDomain string    `gorm:"type:varchar(32)" json:"-"`
	JoinedDomain  string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type userRole struct {
	UserName string
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

func getJoinedDomainByInitialDomain(initialDomain string ) (string ,error){
	var da DomainAdmin
	err := DB.Where(`initial_domain = ?`, initialDomain).First(&da).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "",err
	}
	return da.JoinedDomain,nil
}