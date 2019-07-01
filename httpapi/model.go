package httpapi

import (
	"github.com/casbin/casbin"
	"github.com/jinzhu/gorm"
	"time"
)

var DB *gorm.DB

type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func InitModel(db *gorm.DB) {
	//todo: 去掉debug
	DB = db.Debug()
	AutoMigrate()
	InitEnforcer()
}

func AutoMigrate() {
	DB.AutoMigrate(&Resource{}, &Role{}, &User{},&DomainAdmin{})
}

var Enforcer *casbin.Enforcer

func InitEnforcer() {
	Enforcer = casbin.NewEnforcer("model.conf", false)
	err := loadAllRoleRescourcePolicy(Enforcer)
	if err != nil {
		panic(err)
	}
	err = loadAllRoleUserPolicy(Enforcer)
	if err != nil {
		panic(err)
	}
}


type Role struct {
	Model
	Name      string     `gorm:"type:varchar(25);unique" json:"roleName"`
	Resources []Resource `gorm:"many2many:role_resource;" json:"resources"`
}

type Resource struct {
	Model
	Act     string `gorm:"type:varchar(16)"json:"act"`
	Domain  string `gorm:"type:varchar(32)" json:"domain"`
	Name    string `gorm:"type:varchar(64)" json:"name"`
	Comment string `gorm:"type:varchar(64)" json:"comment"`
}

type User struct {
	Model
	Name          string `json:"username"`
	Pstr          string `json:"-"`
	Salt          string `json:"-"`
	DefaultDomain string `json:"-"`
	Roles         []Role `json:"-" gorm:"many2many:user_role;"`
}

type DomainAdmin struct {
	Model
	Name string
	Pstr string
	Salt string
	InitialDomain string
	JoinedDomain string
}