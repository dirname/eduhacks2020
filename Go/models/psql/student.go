package psql

import (
	"gorm.io/gorm"
	"time"
)

//type Ban struct {
//	ID        uint   `sql:"serial unique"`
//	BanID     string `gorm:"not null;unique"`
//	Message   string
//	CreatedAt time.Time
//	UpdatedAt time.Time
//	DeletedAt gorm.DeletedAt `gorm:"index"`
//}

// Student 学生表的结构
type Student struct {
	ID       uint   `sql:"serial unique"`
	UserID   string `gorm:"not null;unique"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null;"`
	Nickname string `gorm:"not null;"`
	Gender   bool
	Phone    string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Avatar   string
	Birthday time.Time
	Banned   bool `gorm:"default:false"`
	//BannedID  uint
	//Ban       Ban `gorm:"foreignKey:BannedID"`
	ClassID   uint
	Class     Class `gorm:"foreignKey:ClassID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 学生表的表名
func (Student) TableName() string {
	return "student.users"
}
