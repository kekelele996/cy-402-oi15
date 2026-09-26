export interface User {
  id: number
  username: string
  real_name: string
  role: string
  license_no: string
  email: string
  phone: string
  avatar: string
  created_at: string
}

export interface Client {
  id: number
  name: string
  id_number: string
  contact: string
  address: string
  remark: string
  created_at: string
}

export interface CaseItem {
  id: number
  case_no: string
  title: string
  case_type: string
  status: string
  accept_date: string | null
  close_date: string | null
  summary: string
  client_id: number
  lead_lawyer_id: number
  co_lawyer_ids: number[]
  created_at: string
}

export interface DocumentItem {
  id: number
  title: string
  file_type: string
  file_url: string
  upload_time: string
  case_id: number
  uploader_id: number
  created_at: string
}

export interface Billing {
  id: number
  bill_no: string
  billing_type: string
  amount: number | string
  status: string
  case_id: number
  client_id: number
  invoice_info: string
  created_at: string
}

export interface AuditLog {
  id: number
  operator_id: number
  operator_name: string
  action: string
  entity_type: string
  entity_id: string
  detail: string
  ip: string
  created_at: string
}

export interface CaseTodo {
  id: number
  case_id: number
  title: string
  due_date: string | null
  assignee_id: number
  assignee_name: string
  status: string
  completed_at: string | null
  created_at: string
}

export interface CaseTodoStats {
  total_pending: number
  overdue_count: number
  today_count: number
  pending_count: number
  done_count: number
}

export interface CaseTodoBoard {
  overdue: CaseTodo[]
  today: CaseTodo[]
  pending: CaseTodo[]
  done: CaseTodo[]
  stats: CaseTodoStats
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
