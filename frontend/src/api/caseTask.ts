import request from '@/utils/request'

export function listCaseTasks(caseId: number, group?: string) {
  return request.get(`/cases/${caseId}/tasks`, { params: group ? { group } : {} })
}

export function createCaseTask(caseId: number, data: { title: string; due_date?: string | null; assignee_id: number }) {
  return request.post(`/cases/${caseId}/tasks`, data)
}

export function completeCaseTask(caseId: number, taskId: number) {
  return request.post(`/cases/${caseId}/tasks/${taskId}/complete`)
}

export function listCaseAssignees(caseId: number) {
  return request.get(`/cases/${caseId}/assignees`)
}
