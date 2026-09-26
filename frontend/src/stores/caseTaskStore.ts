import { create } from 'zustand'
import { listCaseTasks } from '@/api/caseTask'
import type { CaseTask, CaseTaskSummary } from '@/types'

interface CaseTaskState {
  byCase: CaseTask[]
  summary: CaseTaskSummary
  fetchByCase: (caseId: number, group?: string) => Promise<void>
}

const emptySummary: CaseTaskSummary = { overdue: 0, today: 0, pending: 0, done: 0 }

export const useCaseTaskStore = create<CaseTaskState>((set) => ({
  byCase: [],
  summary: emptySummary,
  async fetchByCase(caseId: number, group?: string) {
    const res: any = await listCaseTasks(caseId, group)
    set({ byCase: res.data.list || [], summary: res.data.summary || emptySummary })
  },
}))
