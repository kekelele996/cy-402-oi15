import { useEffect, useState } from 'react'
import { Badge, Button, DatePicker, Empty, Input, List, Select, Space, Tabs, Tag, message } from 'antd'
import { CheckOutlined, PlusOutlined } from '@ant-design/icons'
import type { Dayjs } from 'dayjs'
import { useCaseTodoStore } from '@/stores/caseTodoStore'
import { CaseTodoGroup, CaseTodoGroupColor, CaseTodoGroupText } from '@/constants/caseTodo'
import { formatDate, formatDateTime } from '@/utils/dateFormat'
import type { CaseTodo } from '@/types'

export default function CaseTodoPanel({ caseId }: { caseId: number }) {
  const store = useCaseTodoStore()
  const [title, setTitle] = useState('')
  const [dueDate, setDueDate] = useState<Dayjs | null>(null)
  const [assigneeId, setAssigneeId] = useState<number>()
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    store.fetchBoard(caseId)
    store.fetchAssignees(caseId)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [caseId])

  async function onCreate() {
    if (!title.trim() || !dueDate || !assigneeId) {
      message.warning('请填写事项、截止日期和负责人')
      return
    }
    setSaving(true)
    try {
      await store.create(caseId, { title: title.trim(), due_date: dueDate.format('YYYY-MM-DD'), assignee_id: assigneeId })
      message.success('待办事项已创建')
      setTitle('')
      setDueDate(null)
      setAssigneeId(undefined)
    } finally {
      setSaving(false)
    }
  }

  async function onComplete(id: number) {
    await store.complete(caseId, id)
    message.success('待办事项已完成')
  }

  const groups = [
    { key: CaseTodoGroup.OVERDUE, list: store.board.overdue, count: store.board.stats.overdue_count },
    { key: CaseTodoGroup.TODAY, list: store.board.today, count: store.board.stats.today_count },
    { key: CaseTodoGroup.PENDING, list: store.board.pending, count: store.board.stats.pending_count },
    { key: CaseTodoGroup.DONE, list: store.board.done, count: store.board.stats.done_count },
  ]

  function renderItem(t: CaseTodo) {
    const done = t.status === 'done'
    return (
      <List.Item
        actions={
          done
            ? []
            : [
                <Button key="done" size="small" type="link" icon={<CheckOutlined />} onClick={() => onComplete(t.id)}>
                  完成
                </Button>,
              ]
        }
      >
        <List.Item.Meta
          title={
            <Space>
              <span style={done ? { textDecoration: 'line-through', color: '#999' } : undefined}>{t.title}</span>
              <Tag>{t.assignee_name || `#${t.assignee_id}`}</Tag>
            </Space>
          }
          description={
            <Space size="middle">
              <span>截止：{formatDate(t.due_date)}</span>
              {done && <span>完成时间：{formatDateTime(t.completed_at)}</span>}
            </Space>
          }
        />
      </List.Item>
    )
  }

  return (
    <div>
      <Space.Compact block style={{ marginBottom: 16, maxWidth: 780 }}>
        <Input
          placeholder="事项（如：开庭、补证、回访）"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          style={{ width: 280 }}
        />
        <DatePicker placeholder="截止日期" value={dueDate} onChange={setDueDate} />
        <Select
          placeholder="负责人"
          style={{ width: 160 }}
          value={assigneeId}
          onChange={setAssigneeId}
          options={store.assignees.map((u) => ({ label: u.real_name || u.username, value: u.id }))}
        />
        <Button type="primary" icon={<PlusOutlined />} loading={saving} onClick={onCreate}>
          添加
        </Button>
      </Space.Compact>
      <Tabs
        items={groups.map((g) => ({
          key: g.key,
          label: (
            <Badge count={g.count} color={CaseTodoGroupColor[g.key]} overflowCount={99}>
              <span style={{ paddingRight: 10 }}>{CaseTodoGroupText[g.key]}</span>
            </Badge>
          ),
          children: g.list.length ? (
            <List<CaseTodo> size="small" dataSource={g.list} renderItem={renderItem} rowKey="id" />
          ) : (
            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={`暂无${CaseTodoGroupText[g.key]}事项`} />
          ),
        }))}
      />
    </div>
  )
}
