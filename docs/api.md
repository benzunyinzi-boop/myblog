# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **响应格式**: JSON

## 通用响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10000 | 服务器内部错误 |
| 10001 | 参数错误 |
| 10002 | 未登录 / token 失效 |
| 10003 | 无权限 |
| 20001 | 文章不存在 |
| 20002 | 文章 slug 已存在 |
| 20003 | 文章状态冲突 |
| 30001 | 用户不存在 |
| 30002 | 密码错误 |
| 30005 | token 无效 |
| 40001 | 分类不存在 |
| 40002 | 分类名或 slug 已存在 |
| 40101 | 标签不存在 |
| 40102 | 标签名或 slug 已存在 |

---

## 公开接口

### 健康检查

**GET** `/public/health`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "service": "myblog",
    "status": "up",
    "ts": 1778237461,
    "version": "dev"
  }
}
```

### 分类列表

**GET** `/public/categories`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Golang",
        "slug": "golang",
        "description": "Go 语言相关",
        "sort_order": 1
      }
    ],
    "total": 1
  }
}
```

### 标签列表

**GET** `/public/tags`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "并发",
        "slug": "concurrency"
      }
    ],
    "total": 1
  }
}
```

### 文章列表

**GET** `/public/articles`

查询参数:
- `page` (int, 可选): 页码,默认 1
- `page_size` (int, 可选): 每页数量,默认 10,最大 100
- `category_id` (int, 可选): 分类 ID
- `tag_id` (int, 可选): 标签 ID
- `keyword` (string, 可选): 关键词(搜索 title/summary)

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Go 并发模式",
        "slug": "go-concurrency",
        "summary": "errgroup + context",
        "cover_image": "",
        "category_id": 1,
        "author_id": 2,
        "status": "published",
        "view_count": 10,
        "published_at": 1778237461,
        "created_at": 1778237400,
        "tags": [
          {"id": 1, "name": "并发", "slug": "concurrency"}
        ]
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

### 文章详情

**GET** `/public/articles/:slug`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "Go 并发模式",
    "slug": "go-concurrency",
    "summary": "errgroup + context",
    "content": "正文内容...",
    "cover_image": "",
    "category_id": 1,
    "author_id": 2,
    "status": "published",
    "view_count": 11,
    "published_at": 1778237461,
    "created_at": 1778237400,
    "tags": [
      {"id": 1, "name": "并发", "slug": "concurrency"}
    ]
  }
}
```

注意:
- 只返回已发布文章(status=published)
- 每次访问自动增加浏览计数

### 个人资料

**GET** `/public/profile`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "name": "张三",
    "bio": "全栈开发者",
    "avatar": "/uploads/2026/05/08/avatar.jpg",
    "email": "example@example.com",
    "github": "https://github.com/username",
    "twitter": "https://twitter.com/username",
    "linkedin": "https://linkedin.com/in/username",
    "website": "https://example.com"
  }
}
```

---

## 管理接口(需认证)

所有管理接口需要在请求头中携带 JWT token:
```
Authorization: Bearer <access_token>
```

### 登录

**POST** `/admin/auth/login`

请求:
```json
{
  "username": "admin",
  "password": "Admin@123"
}
```

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "eyJhbGci...",
    "token_type": "Bearer",
    "expires_at": 1778244664,
    "user": {
      "id": 2,
      "username": "admin",
      "nickname": "Admin",
      "avatar": "",
      "role": "admin"
    }
  }
}
```

### 测试认证

**GET** `/admin/ping`

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "message": "pong",
    "role": "admin",
    "user_id": 2,
    "username": "admin"
  }
}
```

---

## 分类管理

### 创建分类

**POST** `/admin/categories`

请求:
```json
{
  "name": "Golang",
  "slug": "golang",
  "description": "Go 语言相关",
  "sort_order": 1
}
```

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "name": "Golang",
    "slug": "golang",
    "description": "Go 语言相关",
    "sort_order": 1
  }
}
```

### 更新分类

**PUT** `/admin/categories/:id`

请求:
```json
{
  "name": "Golang",
  "slug": "golang",
  "description": "Go 语言与并发编程",
  "sort_order": 10
}
```

### 删除分类

**DELETE** `/admin/categories/:id`

### 列出分类

**GET** `/admin/categories`

---

## 标签管理

### 创建标签

**POST** `/admin/tags`

请求:
```json
{
  "name": "并发",
  "slug": "concurrency"
}
```

### 删除标签

**DELETE** `/admin/tags/:id`

### 列出标签

**GET** `/admin/tags`

---

## 文章管理

### 创建文章

**POST** `/admin/articles`

请求:
```json
{
  "title": "Go 并发模式",
  "slug": "go-concurrency",
  "summary": "errgroup + context",
  "content": "正文内容...",
  "cover_image": "/uploads/2026/05/08/cover.jpg",
  "category_id": 1,
  "tag_ids": [1, 2],
  "status": "draft"
}
```

注意:
- `status` 可选值:`draft`(草稿,默认)、`published`(已发布)
- `author_id` 自动从 JWT token 中获取

### 更新文章

**PUT** `/admin/articles/:id`

请求:
```json
{
  "title": "Go 并发模式",
  "slug": "go-concurrency",
  "summary": "新摘要",
  "content": "更新后的正文...",
  "cover_image": "/uploads/2026/05/08/cover.jpg",
  "category_id": 1,
  "tag_ids": [1]
}
```

注意:
- 不能通过此接口修改 `status`,使用发布/下线接口

### 删除文章

**DELETE** `/admin/articles/:id`

### 获取文章详情

**GET** `/admin/articles/:id`

### 列出文章

**GET** `/admin/articles`

查询参数:
- `page` (int, 可选): 页码
- `page_size` (int, 可选): 每页数量
- `status` (string, 可选): `draft` 或 `published`
- `category_id` (int, 可选): 分类 ID
- `tag_id` (int, 可选): 标签 ID
- `keyword` (string, 可选): 关键词

### 发布文章

**POST** `/admin/articles/:id/publish`

将草稿转为已发布状态,设置 `published_at` 时间戳。

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "status": "published",
    "published_at": 1778237461,
    ...
  }
}
```

错误:
- `20003`: 文章已发布(重复发布)

### 下线文章

**POST** `/admin/articles/:id/unpublish`

将已发布文章转为草稿状态。

错误:
- `20003`: 文章已是草稿(重复下线)

---

## 个人资料管理

### 更新个人资料

**PUT** `/admin/profile`

请求:
```json
{
  "name": "张三",
  "bio": "全栈开发者",
  "avatar": "/uploads/2026/05/08/avatar.jpg",
  "email": "example@example.com",
  "github": "https://github.com/username",
  "twitter": "https://twitter.com/username",
  "linkedin": "https://linkedin.com/in/username",
  "website": "https://example.com"
}
```

注意:
- profiles 表只有一条记录(id=1)
- 首次调用会创建,后续调用会更新

---

## 文件上传

### 上传文件

**POST** `/admin/uploads`

请求:
- Content-Type: `multipart/form-data`
- 字段名: `file`
- 限制: 最大 10MB

响应:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "url": "/uploads/2026/05/08/abc12345-filename.jpg"
  }
}
```

注意:
- 文件按日期分目录存储:`uploads/2026/05/08/`
- 文件名格式:`uuid前8位-原始名.ext`
- 访问 URL:`/uploads/2026/05/08/abc12345-filename.jpg`

---

## 静态文件

**GET** `/uploads/*`

访问上传的文件,如图片、附件等。

示例:
```
GET /uploads/2026/05/08/abc12345-cover.jpg
```
