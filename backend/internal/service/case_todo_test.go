package service

import (
	"testing"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
)

func day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 10, 0, 0, 0, time.Local)
}

func TestGroupCaseTodos(t *testing.T) {
	now := day(2026, 9, 25)
	overdue := day(2026, 9, 20)
	today := day(2026, 9, 25)
	future := day(2026, 9, 30)
	done := day(2026, 9, 18)
	todos := []model.CaseTodo{
		{ID: 1, Title: "开庭", DueDate: &overdue, AssigneeID: 2, Status: constants.CaseTodoStatusPending},
		{ID: 2, Title: "补证", DueDate: &today, AssigneeID: 3, Status: constants.CaseTodoStatusPending},
		{ID: 3, Title: "回访", DueDate: &future, AssigneeID: 2, Status: constants.CaseTodoStatusPending},
		{ID: 4, Title: "已办", DueDate: &overdue, AssigneeID: 3, Status: constants.CaseTodoStatusDone, CompletedAt: &done},
		{ID: 5, Title: "无截止", AssigneeID: 2, Status: constants.CaseTodoStatusPending},
	}
	board := GroupCaseTodos(todos, map[uint64]string{2: "张律师", 3: "李助理"}, now)
	if len(board.Overdue) != 1 || board.Overdue[0].ID != 1 {
		t.Errorf("overdue group = %+v, want only todo 1", board.Overdue)
	}
	if len(board.Today) != 1 || board.Today[0].ID != 2 {
		t.Errorf("today group = %+v, want only todo 2", board.Today)
	}
	if len(board.Pending) != 2 {
		t.Errorf("pending group size = %d, want 2 (future + no due date)", len(board.Pending))
	}
	if len(board.Done) != 1 || board.Done[0].CompletedAt == nil {
		t.Errorf("done group = %+v, want todo 4 with completed_at kept", board.Done)
	}
	if board.Stats.TotalPending != 4 {
		t.Errorf("total pending = %d, want 4 (done excluded)", board.Stats.TotalPending)
	}
	if board.Stats.DoneCount != 1 || board.Stats.OverdueCount != 1 || board.Stats.TodayCount != 1 || board.Stats.PendingCount != 2 {
		t.Errorf("stats = %+v, want {1 1 1 2}", board.Stats)
	}
	if board.Overdue[0].AssigneeName != "张律师" {
		t.Errorf("assignee name = %q, want 张律师", board.Overdue[0].AssigneeName)
	}
}

func TestCaseStaffIDs(t *testing.T) {
	c := &model.Case{LeadLawyerID: 2, CoLawyerIDs: model.CoLawyerJSON([]byte("[3,2,5]"))}
	ids := caseStaffIDs(c)
	if len(ids) != 3 || ids[0] != 2 || ids[1] != 3 || ids[2] != 5 {
		t.Errorf("caseStaffIDs = %v, want [2 3 5] (lead first, dedup)", ids)
	}
	if !isCaseStaff(c, 3) || isCaseStaff(c, 9) {
		t.Error("isCaseStaff mismatch")
	}
	empty := &model.Case{LeadLawyerID: 7}
	if got := caseStaffIDs(empty); len(got) != 1 || got[0] != 7 {
		t.Errorf("caseStaffIDs(no co) = %v, want [7]", got)
	}
}

func TestCaseTodoStatusValidator(t *testing.T) {
	if !constants.IsValidCaseTodoStatus(constants.CaseTodoStatusDone) {
		t.Error("done should be valid")
	}
	if constants.IsValidCaseTodoStatus("bogus") {
		t.Error("bogus should be invalid")
	}
}
