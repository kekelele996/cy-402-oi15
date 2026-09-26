// 案件状态/类型枚举（与后端 backend/internal/constants/case.go 保持一致）
export const CaseStatus = {
  FILED: 'filed',
  INVESTIGATING: 'investigating',
  HEARING: 'hearing',
  CLOSED: 'closed',
  ARCHIVED: 'archived',
} as const

export const CaseStatusText: Record<string, string> = {
  [CaseStatus.FILED]: '已立案',
  [CaseStatus.INVESTIGATING]: '调查取证',
  [CaseStatus.HEARING]: '庭审中',
  [CaseStatus.CLOSED]: '已结案',
  [CaseStatus.ARCHIVED]: '已归档',
}

export const CaseType = {
  CIVIL: 'civil',
  CRIMINAL: 'criminal',
  ADMINISTRATIVE: 'administrative',
  COMMERCIAL: 'commercial',
  LABOR: 'labor',
} as const

export const CaseTypeText: Record<string, string> = {
  [CaseType.CIVIL]: '民事',
  [CaseType.CRIMINAL]: '刑事',
  [CaseType.ADMINISTRATIVE]: '行政',
  [CaseType.COMMERCIAL]: '商事',
  [CaseType.LABOR]: '劳动',
}

export const CaseStatusOptions = Object.entries(CaseStatusText).map(([value, label]) => ({ label, value }))
export const CaseTypeOptions = Object.entries(CaseTypeText).map(([value, label]) => ({ label, value }))

// 案件待办状态/分组（与后端 backend/internal/constants/case_task.go 保持一致）
export const CaseTaskStatus = {
  PENDING: 'pending',
  DONE: 'done',
} as const

export const CaseTaskGroup = {
  OVERDUE: 'overdue',
  TODAY: 'today',
  PENDING: 'pending',
  DONE: 'done',
} as const

export const CaseTaskGroupText: Record<string, string> = {
  [CaseTaskGroup.OVERDUE]: '逾期',
  [CaseTaskGroup.TODAY]: '今日',
  [CaseTaskGroup.PENDING]: '未完成',
  [CaseTaskGroup.DONE]: '已完成',
}
