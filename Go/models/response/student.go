package response

import "time"

// StudentInfo 响应体中学生的信息
type StudentInfo struct {
	UserID      string      `json:"uid"`
	Username    string      `json:"username"`
	Nickname    string      `json:"nickname"`
	Gender      bool        `json:"gender"`
	Phone       string      `json:"phone"`
	Email       string      `json:"email"`
	Avatar      string      `json:"avatar"`
	Birthday    time.Time   `json:"birthday"`
	Banned      bool        `json:"prison"`
	Class       ClassInfo   `json:"class" gorm:"embedded"`
	Major       MajorInfo   `json:"major" gorm:"embedded"`
	CollegeInfo CollegeInfo `json:"college" gorm:"embedded"`
	CreatedAt   time.Time   `json:"create"`
}

// ClassInfo 响应体中的班级信息
type ClassInfo struct {
	ClassName string `json:"name"`
	ClassID   string `json:"id"`
}

// MajorInfo 响应体中的专业信息
type MajorInfo struct {
	MajorName string `json:"name"`
	MajorID   string `json:"id"`
}

// CollegeInfo 响应体中的学院信息
type CollegeInfo struct {
	CollegeName string `json:"name"`
	CollegeID   string `json:"id"`
}
