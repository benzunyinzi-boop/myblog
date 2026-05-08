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
- [x] M2 完成:MySQL/Redis 接入、migration(6 张表)、users+JWT 登录、admin ping 受保护、种子工具
- [x] **M3 完成**:Category/Tag/Article/Profile/Upload 全链路,API 文档,验证计划
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

## 交接下一步

下个会话如果用户说"继续",进入 **M4 · 前端开发** 或 **M5 · 部署**。

**M4 · 前端开发(推荐)**
1. Vite + Vue 3 + TS + Naive UI 脚手架
2. 把 `docs/preview/` 的 CSS 迁到 `web/src/styles/`
3. 对接 M3 API:登录 → 文章列表 → 文章详情 → 关于我
4. 管理后台:文章 CRUD + 分类标签管理 + 图片上传

**M5 · 部署**
1. Dockerfile + docker-compose(MySQL/Redis/myblog)
2. Nginx 反向代理 + 静态文件服务
3. 域名 + HTTPS(Let's Encrypt)
4. 阿里云/腾讯云 VPS 部署

如果用户不指定,默认走 M4(前端开发,让整个博客能跑起来)。

> 本项目使用 Claude Code 协作,用户是十年 Go 后端出身女生开发者,视觉品味偏紫色科技感。
