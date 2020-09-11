package psql

import (
	"gorm.io/gorm"
	"time"
)

// College 学院的表结构
type College struct {
	ID          uint   `sql:"serial unique"`
	CollegeID   string `gorm:"not null;unique"`
	CollegeName string `gorm:"not null;unique"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Major 专业的表结构
type Major struct {
	ID        uint   `sql:"serial unique"`
	MajorID   string `gorm:"not null;unique"`
	MajorName string `gorm:"not null;unique"`
	CollegeID uint
	College   College `gorm:"foreignKey:CollegeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Class 班级的表结构
type Class struct {
	ID        uint   `sql:"serial unique"`
	ClassID   string `gorm:"not null;unique"`
	ClassName string `gorm:"not null;unique"`
	MajorID   uint
	Major     Major `gorm:"foreignKey:MajorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 学院的表名称
func (College) TableName() string {
	return "college.colleges"
}

// TableName 班级的表名称
func (Class) TableName() string {
	return "college.classes"
}

// TableName 专业的表名称
func (Major) TableName() string {
	return "college.majors"
}
