import request from '@/utils/request'

export function listCaseTodos(caseId: number) {
  return request.get(`/cases/${caseId}/todos`)
}

export function listCaseAssignees(caseId: number) {
  return request.get(`/cases/${caseId}/assignees`)
}

export function createCaseTodo(caseId: number, data: { title: string; due_date: string; assignee_id: number }) {
  return request.post(`/cases/${caseId}/todos`, data)
}

export function completeCaseTodo(id: number) {
  return request.post(`/case-todos/${id}/complete`)
}
