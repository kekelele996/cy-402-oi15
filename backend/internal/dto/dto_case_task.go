package dto

import "time"

// CaseTaskCreateRequest 创建案件待办请求。
type CaseTaskCreateRequest struct {
	Title      string  `json:"title" binding:"required,max=200"`
	DueDate    *string `json:"due_date"`
	AssigneeID uint64  `json:"assignee_id" binding:"required"`
}

// ParseDueDate 解析截止日期字符串（格式 2006-01-02，可空）。
func ParseDueDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
