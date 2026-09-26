package repository

import (
	"errors"
	"fmt"

	"cylawcase/internal/model"

	"gorm.io/gorm"
)

// CaseTaskRepository 案件待办仓储。
type CaseTaskRepository struct {
	db *gorm.DB
}

// NewCaseTaskRepository 构造案件待办仓储。
func NewCaseTaskRepository(db *gorm.DB) *CaseTaskRepository {
	return &CaseTaskRepository{db: db}
}

// Create 创建待办事项。
func (r *CaseTaskRepository) Create(t *model.CaseTask) error {
	if err := r.db.Create(t).Error; err != nil {
		return fmt.Errorf("create case task: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询待办事项。
func (r *CaseTaskRepository) FindByID(id uint64) (*model.CaseTask, error) {
	var t model.CaseTask
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find case task by id: %w", err)
	}
	return &t, nil
}

// ListByCase 查询某案件全部待办：未完成在前（按截止日期升序，无截止日期排后），已完成在后。
func (r *CaseTaskRepository) ListByCase(caseID uint64) ([]model.CaseTask, error) {
	var list []model.CaseTask
	if err := r.db.Where("case_id = ?", caseID).
		Order("CASE WHEN status = 'pending' THEN 0 ELSE 1 END").
		Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END").
		Order("due_date ASC").
		Order("id DESC").
		Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list case tasks by case: %w", err)
	}
	return list, nil
}

// ListPendingByCase 查询某案件未完成待办（按截止日期升序，无截止日期排后）。
func (r *CaseTaskRepository) ListPendingByCase(caseID uint64) ([]model.CaseTask, error) {
	var list []model.CaseTask
	if err := r.db.Where("case_id = ? AND status = 'pending'", caseID).
		Order("CASE WHEN due_date IS NULL THEN 1 ELSE 0 END").
		Order("due_date ASC").
		Order("id ASC").
		Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list pending case tasks: %w", err)
	}
	return list, nil
}

// Update 更新待办事项。
func (r *CaseTaskRepository) Update(t *model.CaseTask) error {
	if err := r.db.Save(t).Error; err != nil {
		return fmt.Errorf("update case task: %w", err)
	}
	return nil
}
