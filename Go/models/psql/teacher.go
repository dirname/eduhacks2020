package psql

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

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
	Birthday  time.Time
	MajorID   uint
	Major     Major          `gorm:"foreignKey:MajorID"`
	Classes   pq.StringArray `gorm:"type:text[]"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
