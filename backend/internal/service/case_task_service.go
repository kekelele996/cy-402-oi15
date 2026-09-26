package service

import (
	"log/slog"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
	"cylawcase/internal/util"
)

// CaseTaskService 案件待办业务逻辑。
type CaseTaskService struct {
	repo     *repository.CaseTaskRepository
	caseRepo *repository.CaseRepository
	userRepo *repository.UserRepository
	logger   *slog.Logger
}

// NewCaseTaskService 构造案件待办服务。
func NewCaseTaskService(repo *repository.CaseTaskRepository, caseRepo *repository.CaseRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *CaseTaskService {
	return &CaseTaskService{repo: repo, caseRepo: caseRepo, userRepo: userRepo, logger: logger}
}

// TaskSummary 待办分组统计（逾期/今日为未完成子集，已完成不计入待办统计）。
type TaskSummary struct {
	Overdue int `json:"overdue"`
	Today   int `json:"today"`
	Pending int `json:"pending"`
	Done    int `json:"done"`
}

// Create 为案件创建待办事项，负责人必须是该案件的主办或协办人员。
func (s *CaseTaskService) Create(caseID, creatorID, assigneeID uint64, title string, dueDate *time.Time) (*model.CaseTask, error) {
	c, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "CaseTask[case_id=%d] create: case not found", caseID)
	}
	if !containsU64(caseAssigneeIDs(c), assigneeID) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "CaseTask[case_id="+u64(caseID)+"] create: assignee must be lead or co lawyer")
	}
	if _, err := s.userRepo.FindByID(assigneeID); err != nil {
		return nil, util.Wrap(err, "CaseTask[case_id=%d] create: assignee not found", caseID)
	}
	t := &model.CaseTask{
		CaseID:     caseID,
		Title:      title,
		DueDate:    dueDate,
		AssigneeID: assigneeID,
		Status:     constants.TaskStatusPending,
		CreatedBy:  creatorID,
	}
	if err := s.repo.Create(t); err != nil {
		s.logger.Error(constants.LogCaseTaskCreateFailed, "error", err.Error())
		return nil, util.Wrap(err, "CaseTask[case_id=%d] create failed", caseID)
	}
	s.logger.Info(constants.LogCaseTaskCreateSuccess, "task_id", t.ID, "case_id", caseID)
	t.AssigneeName = s.assigneeNameMap([]uint64{assigneeID})[assigneeID]
	return t, nil
}

// ListByCase 按案件查询待办，可按分组（overdue/today/pending/done）过滤，并返回分组统计。
func (s *CaseTaskService) ListByCase(caseID uint64, group string) ([]model.CaseTask, TaskSummary, error) {
	if group != "" && !constants.IsValidTaskGroup(group) {
		return nil, TaskSummary{}, util.NewAppError(constants.CodeValidationFailed, "CaseTask list: invalid group "+group)
	}
	if _, err := s.caseRepo.FindByID(caseID); err != nil {
		return nil, TaskSummary{}, util.Wrap(err, "CaseTask[case_id=%d] list: case not found", caseID)
	}
	tasks, err := s.repo.ListByCase(caseID)
	if err != nil {
		return nil, TaskSummary{}, util.Wrap(err, "CaseTask[case_id=%d] list failed", caseID)
	}
	now := time.Now()
	summary := SummarizeTasks(tasks, now)
	filtered := FilterTasksByGroup(tasks, group, now)
	names := s.assigneeNameMap(assigneeIDsOf(filtered))
	for i := range filtered {
		filtered[i].AssigneeName = names[filtered[i].AssigneeID]
	}
	return filtered, summary, nil
}

// Complete 完成待办事项：保留记录与完成时间，仅从待办统计中移走。
func (s *CaseTaskService) Complete(id uint64) (*model.CaseTask, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "CaseTask[id=%d] complete find failed", id)
	}
	if t.Status == constants.TaskStatusDone {
		return t, nil
	}
	now := time.Now()
	t.Status = constants.TaskStatusDone
	t.CompletedAt = &now
	if err := s.repo.Update(t); err != nil {
		s.logger.Error(constants.LogCaseTaskCompleteFailed, "error", err.Error())
		return nil, util.Wrap(err, "CaseTask[id=%d] complete save failed", id)
	}
	s.logger.Info(constants.LogCaseTaskCompleteSuccess, "task_id", t.ID, "case_id", t.CaseID)
	t.AssigneeName = s.assigneeNameMap([]uint64{t.AssigneeID})[t.AssigneeID]
	return t, nil
}

// AssigneeCandidates 返回可负责该案件待办的人员（主办 + 协办）。
func (s *CaseTaskService) AssigneeCandidates(caseID uint64) ([]model.User, error) {
	c, err := s.caseRepo.FindByID(caseID)
	if err != nil {
		return nil, util.Wrap(err, "CaseTask[case_id=%d] assignees: case not found", caseID)
	}
	users, err := s.userRepo.ListByIDs(caseAssigneeIDs(c))
	if err != nil {
		return nil, util.Wrap(err, "CaseTask[case_id=%d] assignees query failed", caseID)
	}
	return users, nil
}

// assigneeNameMap 查询用户 ID -> 展示名映射（real_name 为空时回退 username）。
func (s *CaseTaskService) assigneeNameMap(ids []uint64) map[uint64]string {
	names := map[uint64]string{}
	users, err := s.userRepo.ListByIDs(ids)
	if err != nil {
		s.logger.Warn("case task assignee names query failed", "error", err.Error())
		return names
	}
	for _, u := range users {
		names[u.ID] = displayName(u)
	}
	return names
}

// assigneeIDsOf 提取待办列表中的负责人 ID（去重）。
func assigneeIDsOf(tasks []model.CaseTask) []uint64 {
	seen := map[uint64]bool{}
	ids := make([]uint64, 0, len(tasks))
	for _, t := range tasks {
		if !seen[t.AssigneeID] {
			seen[t.AssigneeID] = true
			ids = append(ids, t.AssigneeID)
		}
	}
	return ids
}

// SummarizeTasks 统计待办分组数量：逾期/今日为未完成子集，已完成单独统计。
func SummarizeTasks(tasks []model.CaseTask, now time.Time) TaskSummary {
	var sum TaskSummary
	todayStart, todayEnd := dayRange(now)
	for i := range tasks {
		t := &tasks[i]
		if t.Status == constants.TaskStatusDone {
			sum.Done++
			continue
		}
		sum.Pending++
		if t.DueDate != nil {
			if t.DueDate.Before(todayStart) {
				sum.Overdue++
			} else if t.DueDate.Before(todayEnd) {
				sum.Today++
			}
		}
	}
	return sum
}

// FilterTasksByGroup 按分组过滤待办；空分组返回全部。
func FilterTasksByGroup(tasks []model.CaseTask, group string, now time.Time) []model.CaseTask {
	if group == "" {
		return tasks
	}
	todayStart, todayEnd := dayRange(now)
	out := make([]model.CaseTask, 0, len(tasks))
	for i := range tasks {
		t := &tasks[i]
		switch group {
		case constants.TaskGroupDone:
			if t.Status == constants.TaskStatusDone {
				out = append(out, *t)
			}
		case constants.TaskGroupPending:
			if t.Status != constants.TaskStatusDone {
				out = append(out, *t)
			}
		case constants.TaskGroupOverdue:
			if t.Status != constants.TaskStatusDone && t.DueDate != nil && t.DueDate.Before(todayStart) {
				out = append(out, *t)
			}
		case constants.TaskGroupToday:
			if t.Status != constants.TaskStatusDone && t.DueDate != nil &&
				!t.DueDate.Before(todayStart) && t.DueDate.Before(todayEnd) {
				out = append(out, *t)
			}
		}
	}
	return out
}

// dayRange 返回当日 [起, 止) 时间范围。
func dayRange(now time.Time) (time.Time, time.Time) {
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return start, start.AddDate(0, 0, 1)
}

// caseAssigneeIDs 案件主办与协办人员 ID 集合（去重）。
func caseAssigneeIDs(c *model.Case) []uint64 {
	ids := []uint64{c.LeadLawyerID}
	for _, id := range parseCoLawyers(c.CoLawyerIDs) {
		if !containsU64(ids, id) {
			ids = append(ids, id)
		}
	}
	return ids
}

// displayName 用户展示名。
func displayName(u model.User) string {
	if u.RealName != "" {
		return u.RealName
	}
	return u.Username
}

func containsU64(list []uint64, v uint64) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
