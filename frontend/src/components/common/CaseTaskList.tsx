import { useEffect, useState } from 'react'
import { Badge, Button, DatePicker, Form, Input, Modal, Popconfirm, Select, Space, Table, Tabs, Tag, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { completeCaseTask, createCaseTask } from '@/api/caseTask'
import { useCaseTaskStore } from '@/stores/caseTaskStore'
import { CaseTaskGroup, CaseTaskGroupText, CaseTaskStatus } from '@/constants/case'
import { formatDate, formatDateTime } from '@/utils/dateFormat'
import type { CaseTask, User } from '@/types'

interface Props {
  caseId: number
  candidates: User[]
}

export default function CaseTaskList({ caseId, candidates }: Props) {
  const taskStore = useCaseTaskStore()
  const [group, setGroup] = useState<string>(CaseTaskGroup.PENDING)
  const [modalOpen, setModalOpen] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    taskStore.fetchByCase(caseId, group)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [caseId, group])

  async function refresh() {
    await taskStore.fetchByCase(caseId, group)
  }

  async function onCreate() {
    const values = await form.validateFields()
    setSubmitting(true)
    try {
      await createCaseTask(caseId, {
        title: values.title,
        due_date: values.due_date ? dayjs(values.due_date).format('YYYY-MM-DD') : null,
        assignee_id: values.assignee_id,
      })
      message.success('待办事项已创建')
      setModalOpen(false)
      form.resetFields()
      refresh()
    } finally {
      setSubmitting(false)
    }
  }

  async function onComplete(task: CaseTask) {
    await completeCaseTask(caseId, task.id)
    message.success('待办事项已完成')
    refresh()
  }

  function statusTag(task: CaseTask) {
    if (task.status === CaseTaskStatus.DONE) return <Tag color="green">已完成</Tag>
    if (task.due_date) {
      const due = dayjs(task.due_date)
      if (due.isBefore(dayjs(), 'day')) return <Tag color="red">已逾期</Tag>
      if (due.isSame(dayjs(), 'day')) return <Tag color="orange">今日截止</Tag>
    }
    return <Tag color="blue">待办</Tag>
  }

  const summary = taskStore.summary
  const tabLabel = (key: string, count: number, color?: string) => (
    <Badge count={count} size="small" color={color} offset={[8, -2]}>
      {CaseTaskGroupText[key]}
    </Badge>
  )

  return (
    <>
      <Space style={{ marginBottom: 12, width: '100%', justifyContent: 'space-between' }}>
        <Tabs
          activeKey={group}
          onChange={setGroup}
          items={[
            { key: CaseTaskGroup.OVERDUE, label: tabLabel(CaseTaskGroup.OVERDUE, summary.overdue, 'red') },
            { key: CaseTaskGroup.TODAY, label: tabLabel(CaseTaskGroup.TODAY, summary.today, 'orange') },
            { key: CaseTaskGroup.PENDING, label: tabLabel(CaseTaskGroup.PENDING, summary.pending) },
            { key: CaseTaskGroup.DONE, label: tabLabel(CaseTaskGroup.DONE, summary.done, 'green') },
          ]}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
          新增待办
        </Button>
      </Space>
      <Table<CaseTask>
        rowKey="id"
        size="small"
        dataSource={taskStore.byCase}
        pagination={false}
        locale={{ emptyText: '暂无待办事项' }}
        columns={[
          { title: '事项', dataIndex: 'title' },
          {
            title: '截止日期',
            dataIndex: 'due_date',
            width: 120,
            render: (v: string | null) => formatDate(v),
          },
          {
            title: '负责人',
            dataIndex: 'assignee_name',
            width: 110,
            render: (v: string, row) => v || `#${row.assignee_id}`,
          },
          { title: '状态', key: 'status', width: 100, render: (_, row) => statusTag(row) },
          {
            title: '完成时间',
            dataIndex: 'completed_at',
            width: 160,
            render: (v: string | null) => (v ? formatDateTime(v) : '-'),
          },
          {
            title: '操作',
            key: 'action',
            width: 100,
            render: (_, row) =>
              row.status === CaseTaskStatus.DONE ? null : (
                <Popconfirm title="确认完成该事项？" onConfirm={() => onComplete(row)}>
                  <Button type="link" size="small">
                    完成
                  </Button>
                </Popconfirm>
              ),
          },
        ]}
      />
      <Modal
        title="新增待办"
        open={modalOpen}
        onOk={onCreate}
        onCancel={() => setModalOpen(false)}
        confirmLoading={submitting}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="title" label="事项" rules={[{ required: true, message: '请输入事项内容' }, { max: 200 }]}>
            <Input placeholder="如：开庭、补充证据、客户回访" />
          </Form.Item>
          <Form.Item name="due_date" label="截止日期">
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="assignee_id" label="负责人" rules={[{ required: true, message: '请选择负责人' }]}>
            <Select
              placeholder="仅可选择本案主办/协办人员"
              options={candidates.map((u) => ({ label: u.real_name || u.username, value: u.id }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
