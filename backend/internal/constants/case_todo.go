package constants

// CaseTodoStatus 案件待办状态枚举。
const (
	CaseTodoStatusPending = "pending"
	CaseTodoStatusDone    = "done"
)

// CaseTodoStatusValues 全部待办状态值。
var CaseTodoStatusValues = []string{CaseTodoStatusPending, CaseTodoStatusDone}

// CaseTodoGroup 待办分组枚举：逾期/今日/未完成/已完成。
const (
	CaseTodoGroupOverdue = "overdue"
	CaseTodoGroupToday   = "today"
	CaseTodoGroupPending = "pending"
	CaseTodoGroupDone    = "done"
)

// CaseTodoGroupValues 全部待办分组值。
var CaseTodoGroupValues = []string{CaseTodoGroupOverdue, CaseTodoGroupToday, CaseTodoGroupPending, CaseTodoGroupDone}

// IsValidCaseTodoStatus 校验待办状态。
func IsValidCaseTodoStatus(s string) bool {
	for _, v := range CaseTodoStatusValues {
		if v == s {
			return true
		}
	}
	return false
}
