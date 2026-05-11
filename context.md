# context.md · myblog

> 决策记录与会话交接文档(CLAUDE.md 强制)。

## 项目全景

- 博客微服务,Go 1.22+(本地 1.26)/ Gin / GORM / MySQL / Redis / zap
- 前台 Vue 3 + TS + Naive UI,后台同工程不同路由
- 部署:Docker + docker-compose + Nginx,目标是阿里云/腾讯云 VPS
- 视觉风格:**紫色科技感**(深空紫底 + 霓虹渐变 + 动态网格/流光线条),用户已预览确认(2026-05-08)
- 全局计划文档:`~/.claude/plans/10-curious-flurry.md`

## 关键决策记录

### D1 · 先做视觉预览再写后端(2026-05-08)
- 原因:视觉方向"小清新+科技感"落地效果不确定,用户需要先直观看效果
- 动作:在 `docs/preview/` 交付了 5 张静态 HTML(首页/技术/摄影/关于/文章详情)
- 结果:第一版"莫兰迪小清新"被否,改第二版"紫色科技感",用户确认"挺好的"
- 资产:`docs/preview/styles.css` 里的 CSS 变量和组件将来直接迁移到 `web/src/styles/`

### D2 · 导航栏结构定了四项(2026-05-08)
- 首页 / 技术 / 摄影 / 关于
- 技术栏细分为:Golang / MySQL / Redis / Message Queue / 运维部署
- 意味着后端需要一个"分类(category)"维度,`articles.category_id` 设计成枚举或单独表都可以,决定走 `categories` 独立表,保留灵活性

### D3 · M1 不接 DB(2026-05-08)
- 原因:避免本地没起 MySQL/Redis 阻塞骨架搭建
- 动作:M1 只跑 HTTP 层 + 一个 health 探活
- 下一步(M2):接入 MySQL/Redis、迁移、users + JWT

### D4 · 视觉预览默认暗色主题(2026-05-08)
- 用户反馈"紫色加上科技感",浅色底和紫色的对比不够通透,改为深色优先
- 保留浅色变量作为 `[data-theme="light"]`,右上 `◐` 可切

## 当前状态(2026-05-09)

- [x] 整体实现计划已成型并用户批准
- [x] 视觉预览 5 张 HTML 已交付并被用户认可
- [x] M1 后端骨架完成
- [x] M2 完成:MySQL/Redis 接入、migration(6 张表)、users+JWT 登录、admin ping 受保护、种子工具
- [x] M3 完成:Category/Tag/Article/Profile/Upload 全链路,API 文档,验证计划
- [x] **M3 端到端验证 25/25 全绿**(`test/e2e-verify.sh`)
- [x] git 已 init 并推送到 GitHub(benzunyinzi-boop/myblog)

## 环境约定(本地开发,2026-05-08)

- Redis:`redis-server --daemonize yes`,`/tmp/myblog-redis/redis.pid`
- MySQL:本地 5.7.16,root 密码由用户保管
- myblog 业务账号:`myblog / myblog123`@127.0.0.1(配置走 `configs/config.yaml`,真实密钥不入库)
- 管理员账号:`admin / Admin@123`(仅本地)
- 端口:后端本地跑 `MYBLOG_SERVER_PORT=18080`(8080 被本地 python http server 占用)

## 安全备忘(待处理)

- root 密码曾出现在对话中,**上生产前必须轮换**
- `configs/config.yaml` 已在 .gitignore,不会入库
- JWT secret 目前是 `dev-secret-change-me`,M5 前替换

## 已知小问题(非阻塞)

- Gin 默认不给 GET 路由注册 HEAD handler,`curl -I /api/v1/public/health` 会 404。生产健康检查用 GET 即可,不处理。
- `go.mod` 被 `tidy` 提升到 `go 1.25.0` 以匹配本地 `go1.26.1` 工具链(CLAUDE.md 要求 ≥1.22,兼容)

## 关键决策补充

### D6 · seed 工具幂等处理(2026-05-08)
- `make seed-admin` 查询已有 username,存在则更新密码/角色,不存在则创建
- 原因:开发期密码丢了能直接重置,不用手工删数据

### D7 · GORM 日志桥接到 zap(2026-05-08)
- 自定义 `zapWriter` 实现 `gormlogger.Writer.Printf`,让 GORM 的慢查询/错误进入 zap 流
- 开启 `IgnoreRecordNotFoundError`,避免把"查不到"当错误打

### D8 · service 层错误统一管理(2026-05-08)
- 将散落在各 service 文件中的错误定义统一收拢到 `internal/service/errcode.go`
- 便于查看和维护,为后续扩展预留空间
- handler 层按类型映射到 `pkg/errcode` 的 HTTP 错误码

### D9 · Article.PublishedAt 类型修复(2026-05-08)
- 问题:migration 定义 `DATETIME(3)`,Go 用 `*int64`,导致 MySQL 1292 错误
- 修复:改用 `*time.Time`,GORM 自动转换;dto 保持 `*int64` 给前端
- 教训:model 字段类型要和 migration 定义匹配

### D10 · 小步提交节奏(2026-05-08)
- 用户要求"每一步都要提交",保持小步快跑
- 每个功能模块独立提交:model → repo → service → handler → 验证
- 好处:可追溯、可回滚、commit 历史清晰

### D11 · public 浏览计数必须脱离请求 ctx(2026-05-09)
- 问题:`go func(){ ... c.Request.Context() ... }()` 在响应结束后 ctx 会被 cancel
- 结果:GORM 的 `IncrView` 被中断,view_count 不递增
- 修复:`context.Background()` + `WithTimeout(3s)` + 透传 trace_id

### D12 · Profile schema 与 model 必须同步迁移(2026-05-09)
- 问题:profiles 表仍是 M0 规划的复合结构,但 M3 改成扁平字段(name/email/github/...)
- 结果:GORM 1054 unknown column
- 修复:`000002_profile_schema_rewrite` 重写 schema,并记录 up/down

### D13 · GORM 对缩写词的 snake_case 命名不符合直觉(2026-05-09)
- `GitHub` → `git_hub`,`LinkedIn` → `linked_in`
- 如果 migration 列名想保持 `github` / `linkedin`,必须显式 `gorm:"column:..."`

### D14 · Save 不适合严格模式下的 Upsert(2026-05-09)
- 问题:新记录 `Save()` 可能带零值时间戳(`0000-00-00`),MySQL strict 模式拒绝
- 修复:用 `clause.OnConflict + AssignmentColumns` 明确白名单字段

### D15 · 导航摄影链接必须用真实路由(2026-05-11)
- 问题:`/photo` 之前是 `<a href="#photography">` 锚点,根本不跳页
- 修复:加 `/photo` 路由 + `PhotoView.vue`,摄影导航改 `RouterLink`
- 附带:所有公开导航项(首页除外)hover 出现 `NavHoverPreview` 浮层,强化"可点击"的视觉反馈
- 触屏 / ≤960px 自动隐藏浮层,避免 hover 卡死

### D16 · 有机曲线卡片用 id 作种子(2026-05-11)
- TechView 每篇文章用 `OrganicCard`,SVG path 通过 `Math.sin(id * 9301 + i * 49297)` 作种子生成
- 好处:同一篇文章形状稳定,不同文章形状各异,无需存储额外数据
- `preserveAspectRatio="none"` + `overflow="visible"` 让曲线能跟随卡片拉伸且外扩部分不被裁

### D17 · About/Photo 数据先走前端静态 JSON(2026-05-11)
- 经历时间线 → `web/src/data/experience.ts`
- 摄影作品 → `web/src/data/photos.ts`(Unsplash 占位)
- 原因:本轮优先视觉落地,后期再评估是否做后端化(加 `experiences` JSON 字段或新建 `photos` 表)

## 交接下一步

下个会话如果用户说"继续",默认进入 **M4 · 前端开发**。

**备选**
- 如果用户说"跳到 M5",直接做 Docker / Nginx / 部署
- 如果用户说"先补测试",继续把 e2e 脚本纳入 Makefile/CI

**新窗口恢复时优先看的文件**
1. `tasks.md` — 当前总进度与待办
2. `context.md` — 决策、环境和 bug 教训
3. `docs/api.md` — 后端接口契约
4. `docs/preview/` — 前端视觉参考
5. `docs/M3-e2e-verification-result.md` — 已验证结论

> 本项目使用 Claude Code 协作,用户是十年 Go 后端出身女生开发者,视觉品味偏紫色科技感。
