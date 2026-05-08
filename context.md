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

## 当前状态(2026-05-08)

- [x] 整体实现计划已成型并用户批准
- [x] 视觉预览 5 张 HTML 已交付并被用户认可
- [x] M1 后端骨架完成
- [x] **M2 完成**:MySQL/Redis 接入、migration(6 张表)、users+JWT 登录、admin ping 受保护、种子工具、端到端 7 点验证通过
- [ ] git 未 init(等用户确认后再建 repo)

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

## 交接下一步

下个会话如果用户说"继续",进入 **M3 · 核心 API / 前台**。两个推进方向可选:

**路线 A(推荐):后端继续深挖 M3 后半段的核心 API**
1. 文章 CRUD(repository + service + admin handler)
2. 分类/标签管理
3. 公开接口:文章列表(分页+分类+标签+关键词)、详情(含 Redis 缓存)
4. Profile(关于我)GET/PUT
5. 图片上传(本地 uploads/)

**路线 B:先起 Vue 3 前端脚手架,把预览页迁到 web/**
1. Vite + Vue 3 + TS,接 Naive UI
2. 把 `docs/preview/` 的 CSS 迁到 `web/src/styles/`
3. 打通 /login 页,接上 M2 的 /admin/auth/login

如果用户不指定,默认走 A(先让后端完整,再做前端对接才有数据)。

> 本项目使用 Claude Code 协作,用户是十年 Go 后端出身女生开发者,视觉品味偏紫色科技感。
