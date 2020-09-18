package psql

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Course 课程表的结构
type Course struct {
	ID         uint      `sql:"serial unique"`
	CourseID   uuid.UUID `gorm:"not null;unique;type:uuid;default:uuid_generate_v4()"`
	CourseName string    `gorm:"not null"`
	MajorID    uint
	Major      Major `gorm:"foreignKey:MajorID"`
	TeacherID  uint
	Teacher    Teacher `gorm:"foreignKey:TeacherID"`
	Status     bool
	Open       bool
	StartAt    time.Time
	EndAt      time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName 课程表的表名
func (Course) TableName() string {
	return "resource.courses"
}
