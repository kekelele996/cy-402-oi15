package repository

import (
	"errors"
	"fmt"

	"cylawcase/internal/model"

	"gorm.io/gorm"
)

// CaseTodoRepository 案件待办仓储。
type CaseTodoRepository struct {
	db *gorm.DB
}

// NewCaseTodoRepository 构造案件待办仓储。
func NewCaseTodoRepository(db *gorm.DB) *CaseTodoRepository {
	return &CaseTodoRepository{db: db}
}

// Create 创建待办事项。
func (r *CaseTodoRepository) Create(t *model.CaseTodo) error {
	if err := r.db.Create(t).Error; err != nil {
		return fmt.Errorf("create case todo: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询待办事项。
func (r *CaseTodoRepository) FindByID(id uint64) (*model.CaseTodo, error) {
	var t model.CaseTodo
	if err := r.db.First(&t, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find case todo by id: %w", err)
	}
	return &t, nil
}

// ListByCase 查询某案件的全部待办（含已完成，完成时间保留可查）。
func (r *CaseTodoRepository) ListByCase(caseID uint64) ([]model.CaseTodo, error) {
	var list []model.CaseTodo
	if err := r.db.Where("case_id = ?", caseID).Order("status ASC, due_date ASC, id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list case todos by case: %w", err)
	}
	return list, nil
}

// Update 更新待办事项。
func (r *CaseTodoRepository) Update(t *model.CaseTodo) error {
	if err := r.db.Save(t).Error; err != nil {
		return fmt.Errorf("update case todo: %w", err)
	}
	return nil
}

// CountPendingByCase 统计某案件未完成待办数量。
func (r *CaseTodoRepository) CountPendingByCase(caseID uint64) (int64, error) {
	var n int64
	if err := r.db.Model(&model.CaseTodo{}).
		Where("case_id = ? AND status = ?", caseID, "pending").
		Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count pending case todos: %w", err)
	}
	return n, nil
}

// EarliestPendingByCase 查询某案件截止最早的未完成待办。
func (r *CaseTodoRepository) EarliestPendingByCase(caseID uint64) (*model.CaseTodo, error) {
	var t model.CaseTodo
	if err := r.db.Where("case_id = ? AND status = ?", caseID, "pending").
		Order("due_date ASC, id ASC").First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find earliest pending case todo: %w", err)
	}
	return &t, nil
}
