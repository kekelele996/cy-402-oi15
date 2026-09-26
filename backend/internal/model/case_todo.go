package model

import "time"

// CaseTodo 案件待办事项实体。
type CaseTodo struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	CaseID      uint64     `gorm:"not null;index" json:"case_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  uint64     `gorm:"not null;index" json:"assignee_id"`
	Status      string     `gorm:"size:30;not null;default:pending;index" json:"status"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TableName 指定表名。
func (CaseTodo) TableName() string { return "case_todos" }

// CaseTodoItem 带负责人姓名的待办视图。
type CaseTodoItem struct {
	CaseTodo
	AssigneeName string `json:"assignee_name"`
}

// CaseTodoBoard 待办看板：按逾期/今日/未完成/已完成分组，并附待办统计。
type CaseTodoBoard struct {
	Overdue []CaseTodoItem `json:"overdue"`
	Today   []CaseTodoItem `json:"today"`
	Pending []CaseTodoItem `json:"pending"`
	Done    []CaseTodoItem `json:"done"`
	Stats   CaseTodoStats  `json:"stats"`
}

// CaseTodoStats 待办统计：已完成事项不计入待办。
type CaseTodoStats struct {
	TotalPending int `json:"total_pending"`
	OverdueCount int `json:"overdue_count"`
	TodayCount   int `json:"today_count"`
	PendingCount int `json:"pending_count"`
	DoneCount    int `json:"done_count"`
}
