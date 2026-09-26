import { create } from 'zustand'
import { listCaseTodos, listCaseAssignees, createCaseTodo, completeCaseTodo } from '@/api/caseTodo'
import type { CaseTodoBoard, User } from '@/types'

const emptyBoard: CaseTodoBoard = {
  overdue: [],
  today: [],
  pending: [],
  done: [],
  stats: { total_pending: 0, overdue_count: 0, today_count: 0, pending_count: 0, done_count: 0 },
}

interface CaseTodoState {
  board: CaseTodoBoard
  assignees: User[]
  fetchBoard: (caseId: number) => Promise<void>
  fetchAssignees: (caseId: number) => Promise<void>
  create: (caseId: number, data: { title: string; due_date: string; assignee_id: number }) => Promise<void>
  complete: (caseId: number, id: number) => Promise<void>
}

export const useCaseTodoStore = create<CaseTodoState>((set, get) => ({
  board: emptyBoard,
  assignees: [],
  async fetchBoard(caseId) {
    const res: any = await listCaseTodos(caseId)
    set({ board: res.data })
  },
  async fetchAssignees(caseId) {
    const res: any = await listCaseAssignees(caseId)
    set({ assignees: res.data })
  },
  async create(caseId, data) {
    await createCaseTodo(caseId, data)
    await get().fetchBoard(caseId)
  },
  async complete(caseId, id) {
    await completeCaseTodo(id)
    await get().fetchBoard(caseId)
  },
}))
