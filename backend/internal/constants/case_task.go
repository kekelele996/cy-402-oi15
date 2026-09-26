package constants

// CaseTaskStatus 待办事项状态枚举。
const (
	TaskStatusPending = "pending"
	TaskStatusDone    = "done"
)

// CaseTaskGroup 待办事项分组枚举：逾期 / 今日 / 未完成 / 已完成。
const (
	TaskGroupOverdue = "overdue"
	TaskGroupToday   = "today"
	TaskGroupPending = "pending"
	TaskGroupDone    = "done"
)

// TaskGroupValues 全部待办分组值。
var TaskGroupValues = []string{TaskGroupOverdue, TaskGroupToday, TaskGroupPending, TaskGroupDone}

// IsValidTaskGroup 校验待办分组。
func IsValidTaskGroup(s string) bool {
	for _, v := range TaskGroupValues {
		if v == s {
			return true
		}
	}
	return false
}
