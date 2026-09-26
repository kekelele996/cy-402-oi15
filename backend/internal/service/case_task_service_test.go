package service

import (
	"testing"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
)

func taskAt(status string, due *time.Time) model.CaseTask {
	return model.CaseTask{Status: status, DueDate: due}
}

func daysFrom(now time.Time, d int) *time.Time {
	t := now.AddDate(0, 0, d)
	return &t
}

func TestSummarizeTasks(t *testing.T) {
	now := time.Now()
	tasks := []model.CaseTask{
		taskAt(constants.TaskStatusPending, daysFrom(now, -2)), // 逾期
		taskAt(constants.TaskStatusPending, daysFrom(now, 0)),  // 今日
		taskAt(constants.TaskStatusPending, daysFrom(now, 5)),  // 未完成（未来）
		taskAt(constants.TaskStatusPending, nil),               // 未完成（无截止日期）
		taskAt(constants.TaskStatusDone, daysFrom(now, -1)),    // 已完成
	}
	sum := SummarizeTasks(tasks, now)
	if sum.Overdue != 1 {
		t.Errorf("overdue = %d, want 1", sum.Overdue)
	}
	if sum.Today != 1 {
		t.Errorf("today = %d, want 1", sum.Today)
	}
	if sum.Pending != 4 {
		t.Errorf("pending = %d, want 4", sum.Pending)
	}
	if sum.Done != 1 {
		t.Errorf("done = %d, want 1", sum.Done)
	}
}

func TestFilterTasksByGroup(t *testing.T) {
	now := time.Now()
	tasks := []model.CaseTask{
		taskAt(constants.TaskStatusPending, daysFrom(now, -2)),
		taskAt(constants.TaskStatusPending, daysFrom(now, 0)),
		taskAt(constants.TaskStatusPending, nil),
		taskAt(constants.TaskStatusDone, daysFrom(now, -1)),
	}
	cases := []struct {
		group string
		want  int
	}{
		{"", 4},
		{constants.TaskGroupOverdue, 1},
		{constants.TaskGroupToday, 1},
		{constants.TaskGroupPending, 3},
		{constants.TaskGroupDone, 1},
	}
	for _, tc := range cases {
		if got := len(FilterTasksByGroup(tasks, tc.group, now)); got != tc.want {
			t.Errorf("FilterTasksByGroup(%q) = %d tasks, want %d", tc.group, got, tc.want)
		}
	}
}

func TestCaseAssigneeIDs(t *testing.T) {
	c := &model.Case{LeadLawyerID: 2, CoLawyerIDs: model.CoLawyerJSON(`[3,2,5]`)}
	ids := caseAssigneeIDs(c)
	if len(ids) != 3 || ids[0] != 2 {
		t.Fatalf("caseAssigneeIDs = %v, want lead first and deduped", ids)
	}
	if !containsU64(ids, 3) || !containsU64(ids, 5) {
		t.Errorf("caseAssigneeIDs = %v, want to contain 3 and 5", ids)
	}
	if containsU64(ids, 4) {
		t.Error("caseAssigneeIDs should not contain outsider 4")
	}
}

func TestParseCoLawyers(t *testing.T) {
	if ids := parseCoLawyers(model.CoLawyerJSON(`[1,2]`)); len(ids) != 2 {
		t.Errorf("parseCoLawyers = %v, want 2 ids", ids)
	}
	if ids := parseCoLawyers(nil); len(ids) != 0 {
		t.Errorf("parseCoLawyers(nil) = %v, want empty", ids)
	}
}
