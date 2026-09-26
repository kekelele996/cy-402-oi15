package service

import (
	"encoding/json"
	"log/slog"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
	"cylawcase/internal/util"
)

// CaseTodoService 案件待办业务逻辑。
type CaseTodoService struct {
	repo     *repository.CaseTodoRepository
	caseRepo *repository.CaseRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewCaseTodoService 构造案件待办服务。
func NewCaseTodoService(repo *repository.CaseTodoRepository, caseRepo *repository.CaseRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *CaseTodoService {
	return &CaseTodoService{repo: repo, caseRepo: caseRepo, userRepo: userRepo, logger: logger}
}

// Create 为案件创建待办事项，负责人必须是该案件的主办或协办人员。
func (s *CaseTodoService) Create(caseID, assigneeID uint64, title string, dueDate *time.Time) (*model.CaseTodo, error) {
	c, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "CaseTodo[case_id=%d] create: case not found", caseID)
	}
	if !isCaseStaff(c, assigneeID) {
		return nil, util.NewAppError(constants.CodeValidationFailed,
			"CaseTodo[case_id="+u64(caseID)+" assignee_id="+u64(assigneeID)+"] create: "+constants.MsgCaseTodoAssigneeInvalid)
	}
	t := &model.CaseTodo{
		CaseID:     caseID,
		Title:      title,
		DueDate:    dueDate,
		AssigneeID: assigneeID,
		Status:     constants.CaseTodoStatusPending,
	}
	if err := s.repo.Create(t); err != nil {
		s.logger.Error(constants.LogCaseTodoCreateFailed, "error", err.Error())
		return nil, util.Wrap(err, "CaseTodo[case_id=%d] create failed", caseID)
	}
	s.logger.Info(constants.LogCaseTodoCreateSuccess, "todo_id", t.ID, "case_id", caseID, "assignee_id", assigneeID)
	return t, nil
}

// Complete 完成待办事项，保留完成时间。
func (s *CaseTodoService) Complete(id uint64) (*model.CaseTodo, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "CaseTodo[id=%d] complete find failed", id)
	}
	if t.Status != constants.CaseTodoStatusPending {
		s.logger.Warn(constants.LogCaseTodoCompleteFailed, "todo_id", id, "status", t.Status)
		return nil, util.NewAppError(constants.CodeCaseTodoStatusConflict,
			"CaseTodo[id="+u64(id)+"] complete failed: status="+util.CaseTodoStatusText(t.Status))
	}
	now := time.Now()
	t.Status = constants.CaseTodoStatusDone
	t.CompletedAt = &now
	if err := s.repo.Update(t); err != nil {
		return nil, util.Wrap(err, "CaseTodo[id=%d] complete save failed", id)
	}
	s.logger.Info(constants.LogCaseTodoCompleteSuccess, "todo_id", t.ID, "case_id", t.CaseID)
	return t, nil
}

// ListByCase 查询某案件待办看板，按逾期/今日/未完成/已完成分组。
func (s *CaseTodoService) ListByCase(caseID uint64) (*model.CaseTodoBoard, error) {
	if _, err := s.caseRepo.FindByID(caseID); err != nil {
		return nil, util.Wrap(err, "CaseTodo[case_id=%d] list: case not found", caseID)
	}
	todos, err := s.repo.ListByCase(caseID)
	if err != nil {
		return nil, util.Wrap(err, "CaseTodo[case_id=%d] list failed", caseID)
	}
	names := map[uint64]string{}
	for _, t := range todos {
		if _, ok := names[t.AssigneeID]; !ok {
			if u, err := s.userRepo.FindByID(t.AssigneeID); err == nil {
				names[t.AssigneeID] = displayName(u)
			}
		}
	}
	board := GroupCaseTodos(todos, names, time.Now())
	s.logger.Info(constants.LogCaseTodoBoard, "case_id", caseID,
		util.CaseTodoGroupText(constants.CaseTodoGroupOverdue), board.Stats.OverdueCount,
		util.CaseTodoGroupText(constants.CaseTodoGroupToday), board.Stats.TodayCount,
		util.CaseTodoGroupText(constants.CaseTodoGroupPending), board.Stats.PendingCount,
		util.CaseTodoGroupText(constants.CaseTodoGroupDone), board.Stats.DoneCount)
	return &board, nil
}

// Assignees 查询某案件可担当负责人的员工（主办 + 协办）。
func (s *CaseTodoService) Assignees(caseID uint64) ([]model.User, error) {
	c, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "CaseTodo[case_id=%d] assignees: case not found", caseID)
	}
	users := make([]model.User, 0)
	for _, id := range caseStaffIDs(c) {
		u, err := s.userRepo.FindByID(id)
		if err != nil {
			s.logger.Warn(constants.LogCaseTodoAssigneeMissing, "case_id", caseID, "user_id", id)
			continue
		}
		users = append(users, *u)
	}
	return users, nil
}

// GroupCaseTodos 按逾期/今日/未完成/已完成分组并统计，已完成不计入待办统计。
func GroupCaseTodos(todos []model.CaseTodo, names map[uint64]string, now time.Time) model.CaseTodoBoard {
	board := model.CaseTodoBoard{
		Overdue: []model.CaseTodoItem{},
		Today:   []model.CaseTodoItem{},
		Pending: []model.CaseTodoItem{},
		Done:    []model.CaseTodoItem{},
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for _, t := range todos {
		item := model.CaseTodoItem{CaseTodo: t, AssigneeName: names[t.AssigneeID]}
		if t.Status == constants.CaseTodoStatusDone {
			board.Done = append(board.Done, item)
			continue
		}
		if t.DueDate != nil {
			d := t.DueDate.In(now.Location())
			day := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, now.Location())
			switch {
			case day.Before(today):
				board.Overdue = append(board.Overdue, item)
				continue
			case day.Equal(today):
				board.Today = append(board.Today, item)
				continue
			}
		}
		board.Pending = append(board.Pending, item)
	}
	board.Stats = model.CaseTodoStats{
		OverdueCount: len(board.Overdue),
		TodayCount:   len(board.Today),
		PendingCount: len(board.Pending),
		DoneCount:    len(board.Done),
	}
	board.Stats.TotalPending = board.Stats.OverdueCount + board.Stats.TodayCount + board.Stats.PendingCount
	return board
}

// caseStaffIDs 返回案件主办与协办人员 ID（去重，主办在前）。
func caseStaffIDs(c *model.Case) []uint64 {
	ids := []uint64{c.LeadLawyerID}
	seen := map[uint64]bool{c.LeadLawyerID: true}
	var co []uint64
	if len(c.CoLawyerIDs) > 0 {
		_ = json.Unmarshal(c.CoLawyerIDs, &co)
	}
	for _, id := range co {
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

// isCaseStaff 判断用户是否为该案件的主办或协办人员。
func isCaseStaff(c *model.Case, userID uint64) bool {
	for _, id := range caseStaffIDs(c) {
		if id == userID {
			return true
		}
	}
	return false
}

// displayName 用户显示名，优先真实姓名。
func displayName(u *model.User) string {
	if u.RealName != "" {
		return u.RealName
	}
	return u.Username
}
