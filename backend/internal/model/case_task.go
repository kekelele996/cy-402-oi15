package model

import "time"

// CaseTask 案件待办事项实体。
type CaseTask struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	CaseID       uint64     `gorm:"not null;index" json:"case_id"`
	Title        string     `gorm:"size:200;not null" json:"title"`
	DueDate      *time.Time `json:"due_date"`
	AssigneeID   uint64     `gorm:"not null;index" json:"assignee_id"`
	Status       string     `gorm:"size:20;not null;default:pending;index" json:"status"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedBy    uint64     `gorm:"not null;default:0" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	AssigneeName string     `gorm:"-" json:"assignee_name"`
}

// TableName 指定表名。
func (CaseTask) TableName() string { return "case_tasks" }
