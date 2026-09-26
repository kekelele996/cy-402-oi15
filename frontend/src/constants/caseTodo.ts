// 案件待办状态/分组枚举（与后端 backend/internal/constants/case_todo.go 保持一致）
export const CaseTodoStatus = {
  PENDING: 'pending',
  DONE: 'done',
} as const

export const CaseTodoStatusText: Record<string, string> = {
  [CaseTodoStatus.PENDING]: '未完成',
  [CaseTodoStatus.DONE]: '已完成',
}

export const CaseTodoGroup = {
  OVERDUE: 'overdue',
  TODAY: 'today',
  PENDING: 'pending',
  DONE: 'done',
} as const

export const CaseTodoGroupText: Record<string, string> = {
  [CaseTodoGroup.OVERDUE]: '逾期',
  [CaseTodoGroup.TODAY]: '今日',
  [CaseTodoGroup.PENDING]: '未完成',
  [CaseTodoGroup.DONE]: '已完成',
}

export const CaseTodoGroupColor: Record<string, string> = {
  [CaseTodoGroup.OVERDUE]: 'red',
  [CaseTodoGroup.TODAY]: 'orange',
  [CaseTodoGroup.PENDING]: 'blue',
  [CaseTodoGroup.DONE]: 'green',
}
