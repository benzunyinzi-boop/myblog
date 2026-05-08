# M3 端到端验证执行报告

## 执行时间

2026-05-08 19:55

## 执行方式

```bash
# 清理数据
MYSQL_PWD='myblog123' mysql -umyblog myblog -e "DELETE FROM articles; DELETE FROM article_tags; DELETE FROM categories; DELETE FROM tags; DELETE FROM profiles;"

# 启动服务
make build && MYBLOG_SERVER_PORT=18080 ./bin/myblog-server &

# 跑验证
./test/e2e-verify.sh
```

## 最终结果

**25 通过 / 0 失败 ✓**

| 模块 | 场景 | 结果 |
|------|------|------|
| 认证 | 登录成功 / 错误密码 / 无 token | 3/3 ✓ |
| 分类 | 创建 / 重复 slug / public 列表 | 3/3 ✓ |
| 标签 | 创建两个标签 | 2/2 ✓ |
| 文章管理 | 创建草稿 / 无效分类 / 重复 slug / 发布 / 重复发布 | 5/5 ✓ |
| 公开文章 | 列表 / 分类过滤 / 标签过滤 / 关键词 / 详情 / 浏览计数 / 404 | 7/7 ✓ |
| Profile | 更新 / public 读取 | 2/2 ✓ |
| 生命周期 | 下线 / 删除 / 删除后 404 | 3/3 ✓ |

## 验证过程中发现并修复的 Bug

### Bug 1:浏览计数不递增

**根因**:`go func() { _ = h.svc.IncrView(c.Request.Context(), ...) }()`
使用 gin 请求的 context,HTTP 响应发出后 ctx 被 cancel,
goroutine 里的 GORM 操作被中断。

**修复**(`25eb86b`):
- 用 `context.Background()` 派生独立 ctx + 3s 超时
- 保留 trace_id 传递
- 失败记 Warn 日志

### Bug 2:Profile schema 与 model 不匹配

**根因**:`000001_init_schema` 里 profiles 表是 M0 规划的复合结构
(user_id / skills JSON / experiences JSON / contacts JSON),
但 M3 Profile model 重新设计成扁平字段
(name / email / github / twitter / linkedin / website),
**缺少配套 migration**,导致 GORM 找不到列。

**修复**(`9d6ef1d`):
- 新增 `000002_profile_schema_rewrite`
- UP:清老列 → 加新列 → bio NOT NULL
- DOWN:完整回滚到 M0 结构

### Bug 3:GORM 命名规约把 GitHub 转成 git_hub

**根因**:GORM 默认 `snake_case` 规约把 `GitHub` → `git_hub`、
`LinkedIn` → `linked_in`,但 migration 列名是 `github` / `linkedin`。

**修复**(`55bf31e`):
- model 字段显式写 `gorm:"column:xxx"` 标签

### Bug 4:Profile Upsert 零值时间戳

**根因**:`db.Save()` 对新记录会生成 INSERT,但 `CreatedAt` 是零值,
MySQL 严格模式拒绝 `'0000-00-00 00:00:00'`。

**修复**(`55bf31e`):
- 改用 `clause.OnConflict` + `AssignmentColumns` 白名单
- 只更新业务字段 + `updated_at`,`created_at` 零值时手动填

## 修复 commit 链

```
55bf31e fix(profile): correct GORM column names and upsert strategy
9d6ef1d fix(migration): rewrite profiles schema to match M3 model
25eb86b fix(handler): use detached context for async view count incr
```

## 关键教训

1. **model 变更必须配套 migration**——这次 Profile model 重写时没同步 migration,
   留下 schema 脱钩。后续流程应该是:改 model 同时写 migration,再跑验证。

2. **GORM 默认命名规约是 snake_case,且按词切分**——`GitHub` 不是 `github`,
   而是 `git_hub`。遇到缩写词(GitHub / LinkedIn / URL / ID)必须显式 `column` 标签。

3. **`Save()` 不适合 Upsert 场景**——新记录会带零值时间戳,在严格模式下会爆。
   应该用 `clause.OnConflict + AssignmentColumns` 明确白名单字段。

4. **异步 goroutine 不能用请求 ctx**——请求结束 ctx 就 cancel,
   背景任务要么用 `context.Background()`,要么用专门的工作池。

## 验证脚本

`test/e2e-verify.sh` 包含:
- 7 个验证模块
- 25 个独立断言
- JSON 路径取值 helper
- 彩色输出 + 统计总结
- 退出码:0 全绿 / 1 有失败
