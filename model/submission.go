package model

import "time"

type Submission struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	ProblemID  uint      `json:"problem_id"`
	Language   string    `json:"language"`              // 如 cpp, python, java
	Code       string    `gorm:"type:text" json:"code"` // 用户提交的代码
	Status     string    `json:"status"`                // Pending, Judging, Accepted, WrongAnswer, etc.
	Result     string    `json:"result"`                // 附加信息或错误输出
	Score      int       `json:"score"`                 // 得分（支持 ACM/IOI 模式）
	SubmitTime time.Time `json:"submit_time"`
}
