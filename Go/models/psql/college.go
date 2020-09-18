package psql

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// College 学院的表结构
type College struct {
	ID          uint      `sql:"serial unique"`
	CollegeID   uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	CollegeName string    `gorm:"not null;unique"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Major 专业的表结构
type Major struct {
	ID        uint      `sql:"serial unique"`
	MajorID   uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	MajorName string    `gorm:"not null;unique"`
	CollegeID uint
	College   College `gorm:"foreignKey:CollegeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Class 班级的表结构
type Class struct {
	ID        uint      `sql:"serial unique"`
	ClassID   uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	ClassName string    `gorm:"not null;unique"`
	MajorID   uint
	Major     Major `gorm:"foreignKey:MajorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Administration 教务的架构表
type Administration struct {
	ID        uint      `sql:"serial unique"`
	UserID    uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null;"`
	Nickname  string    `gorm:"not null;"`
	CollegeID uint
	College   College `gorm:"foreignKey:CollegeID"`
	Gender    bool
	Phone     string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Avatar    string
	Birthday  time.Time
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

// TableName 教务的表名称
func (Administration) TableName() string {
	return "college.users"
}
