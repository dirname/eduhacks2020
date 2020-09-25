package response

import "time"

// TeacherInfo 响应体中教师的信息
type TeacherInfo struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"uid"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Gender    bool      `json:"gender"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FullTime  bool      `json:"fullTime"`
	Avatar    string    `json:"avatar"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"create"`
}
