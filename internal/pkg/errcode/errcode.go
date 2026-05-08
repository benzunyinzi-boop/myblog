// Package errcode 业务错误码表。
// 约定:
//   - 0        成功
//   - 1xxxx    通用错误(参数/鉴权/限流)
//   - 2xxxx    文章相关
//   - 3xxxx    用户相关
//   - 4xxxx    分类/标签
//   - 5xxxx    文件上传
package errcode

// Code 业务错误码
type Code int

const (
	OK Code = 0

	// 通用 1xxxx
	ErrInternal      Code = 10000 // 服务器内部错误
	ErrInvalidParam  Code = 10001 // 参数错误
	ErrUnauthorized  Code = 10002 // 未登录 / token 失效
	ErrForbidden     Code = 10003 // 无权限
	ErrNotFound      Code = 10004 // 资源不存在
	ErrRateLimited   Code = 10005 // 频率限制
	ErrTooManyConn   Code = 10006 // 连接过多
	ErrMethodNotFit  Code = 10007 // 方法不允许

	// 用户 3xxxx
	ErrUserNotFound    Code = 30001 // 用户不存在
	ErrInvalidPassword Code = 30002 // 密码错误
	ErrUserDisabled    Code = 30003 // 账号被禁用
	ErrTokenExpired    Code = 30004 // token 过期
	ErrTokenInvalid    Code = 30005 // token 无效
)

// Message 默认错误信息(handler 可以覆盖)
var Message = map[Code]string{
	OK:                 "ok",
	ErrInternal:        "服务器内部错误",
	ErrInvalidParam:    "参数错误",
	ErrUnauthorized:    "请先登录",
	ErrForbidden:       "没有权限",
	ErrNotFound:        "资源不存在",
	ErrRateLimited:     "请求过于频繁",
	ErrTooManyConn:     "连接数过多",
	ErrMethodNotFit:    "不支持的请求方法",
	ErrUserNotFound:    "用户不存在",
	ErrInvalidPassword: "用户名或密码错误",
	ErrUserDisabled:    "账号已被禁用",
	ErrTokenExpired:    "登录已过期,请重新登录",
	ErrTokenInvalid:    "登录凭证无效",
}

// MsgOf 取默认 message,未登记时返回 "unknown".
func MsgOf(c Code) string {
	if m, ok := Message[c]; ok {
		return m
	}
	return "unknown"
}
