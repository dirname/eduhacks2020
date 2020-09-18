package response

import (
	"gorm.io/gorm"
	"time"
)

// AdministrationInfo 教务的响应结构体
type AdministrationInfo struct {
	UserID    string `json:"uid"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Gender    bool   `json:"gender"`
	Phone     string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Avatar    string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}