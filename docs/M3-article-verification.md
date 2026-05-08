# M3 Article 管理端 CRUD 端到端验证

## 验证时间
2026-05-08 18:41

## 验证范围
- 文章 CRUD 全链路
- 发布/下线状态机
- 事务一致性(article + article_tags)
- 错误码映射

## 验证结果

| # | 操作 | 预期 | 实际 | 状态 |
|---|---|---|---|---|
| 1 | 创建分类 | 200 + id | id=7 | ✓ |
| 2 | 创建标签 | 200 + id | id=6 | ✓ |
| 3 | 创建草稿文章 | 200 + id + status=draft | id=3, status=draft | ✓ |
| 4 | 发布文章 | status=published | status=published | ✓ |
| 5 | 验证 published_at | unix 秒时间戳 | 1778237461 | ✓ |
| 6 | 重复发布 | code=20003 冲突 | code=20003 | ✓ |
| 7 | 下线文章 | status=draft | status=draft | ✓ |
| 8 | 删除文章 | code=0 | code=0 | ✓ |
| 9 | 再次查询 | code=20001 不存在 | code=20001 | ✓ |

## 关键修复

### published_at 类型不匹配 (commit 71ed963)
- **问题**:migration 定义 `DATETIME(3)`,Go 用 `*int64`,导致 MySQL 1292 错误
- **修复**:改用 `*time.Time`,GORM 自动转换;dto 保持 `*int64` 给前端

## 已实现接口

### 管理端(需 JWT)
- `GET    /api/v1/admin/articles` - 列表(分页+过滤)
- `POST   /api/v1/admin/articles` - 创建
- `GET    /api/v1/admin/articles/:id` - 详情
- `PUT    /api/v1/admin/articles/:id` - 更新
- `DELETE /api/v1/admin/articles/:id` - 删除
- `POST   /api/v1/admin/articles/:id/publish` - 发布
- `POST   /api/v1/admin/articles/:id/unpublish` - 下线

### 过滤能力
- status(draft/published)
- category_id
- tag_id(子查询)
- keyword(title/summary 模糊匹配)
- 分页(page/page_size)

## 事务保证
- Create/Update:article + article_tags 联写
- Delete:先清 article_tags,再软删 article
- ReplaceTags:先删后插 + 去重

## 错误码映射
- 20001:文章不存在
- 20002:slug 重复
- 20003:状态冲突(重复发布/下线草稿)
- 40001:关联分类不存在
- 40101:关联标签不存在

## 提交记录
- 510b5a5: feat(model): add Article model with status machine
- edcadc8: feat(repo): add ArticleRepo with tx and filter support
- 9561fd0: feat(service): add ArticleService with tx and publish lifecycle
- 03569a4: feat(handler): wire article admin endpoints
- 71ed963: fix(model): change Article.PublishedAt from *int64 to *time.Time
