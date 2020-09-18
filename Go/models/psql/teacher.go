package psql

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Teacher 教师表的结构
type Teacher struct {
	ID        uint      `sql:"serial unique"`
	UserID    uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	Username  string    `gorm:"not null;unique"`
	Password  string    `gorm:"not null;"`
	Nickname  string    `gorm:"not null;"`
	Gender    bool
	Phone     string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Avatar    string
	FullTime  bool
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 教师表的表名
func (Teacher) TableName() string {
	return "teacher.users"
}
