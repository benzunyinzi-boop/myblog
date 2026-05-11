export type Experience = {
  period: string
  role: string
  org: string
  stack: string[]
  highlights: string[]
}

export const identities = [
  'Go 后端工程师',
  '分布式系统观察者',
  '偶尔摄影师',
  '永远在写下一篇'
]

export const experiences: Experience[] = [
  {
    period: '2023.06 — Now',
    role: 'Senior Backend Engineer',
    org: 'myblog.dev',
    stack: ['Go', 'MySQL', 'Redis', 'Kafka'],
    highlights: [
      '主导高并发 API 网关重构，p99 从 380ms 降到 95ms',
      '设计多活缓存结构，击穿与雪崩事故清零',
      '搭建自动化部署与观测体系，灰度周期从 3 天压到 4 小时'
    ]
  },
  {
    period: '2020.03 — 2023.05',
    role: 'Backend Engineer',
    org: '某互联网公司',
    stack: ['Go', 'PostgreSQL', 'RabbitMQ'],
    highlights: [
      '负责订单与支付核心域，日均成功交易 120w+',
      '推动从单体到服务拆分，领域边界与事务模型标准化',
      '主导一次涉及 200+ 接口的兼容性迁移'
    ]
  },
  {
    period: '2017.07 — 2020.02',
    role: 'Software Engineer',
    org: '入行起点',
    stack: ['Python', 'MySQL', 'Redis'],
    highlights: [
      '从零搭建业务后台，完成首次可规模化落地',
      '接触生产事故与复盘，开始记录工程笔记'
    ]
  }
]
