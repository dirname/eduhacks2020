package response

import "time"

// TeacherInfo 响应体中教师的信息
type TeacherInfo struct {
	UserID    string    `json:"uid"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Gender    bool      `json:"gender"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FullTime  bool      `json:"full-time"`
	Avatar    string    `json:"avatar"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"create"`
}
