// Package response 统一 HTTP 响应格式:{code, message, data}
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yinyin/myblog/internal/pkg/errcode"
)

// Body 统一响应结构
type Body struct {
	Code    errcode.Code `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
}

// OK 返回 200 + code=0
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    errcode.OK,
		Message: errcode.MsgOf(errcode.OK),
		Data:    data,
	})
}

// Fail 根据业务错误码返回
// httpStatus 为 0 时使用默认映射(Unauthorized→401/Forbidden→403/NotFound→404/其他→200 业务错误)
func Fail(c *gin.Context, code errcode.Code, msg string, httpStatus ...int) {
	status := http.StatusOK
	if len(httpStatus) > 0 && httpStatus[0] != 0 {
		status = httpStatus[0]
	} else {
		status = mapHTTPStatus(code)
	}
	if msg == "" {
		msg = errcode.MsgOf(code)
	}
	c.JSON(status, Body{Code: code, Message: msg})
}

// FailInvalid 参数错误快捷方法
func FailInvalid(c *gin.Context, err error) {
	msg := errcode.MsgOf(errcode.ErrInvalidParam)
	if err != nil {
		msg = err.Error()
	}
	c.JSON(http.StatusBadRequest, Body{
		Code:    errcode.ErrInvalidParam,
		Message: msg,
	})
}

func mapHTTPStatus(code errcode.Code) int {
	switch code {
	case errcode.ErrUnauthorized:
		return http.StatusUnauthorized
	case errcode.ErrForbidden:
		return http.StatusForbidden
	case errcode.ErrNotFound:
		return http.StatusNotFound
	case errcode.ErrInternal:
		return http.StatusInternalServerError
	case errcode.ErrRateLimited, errcode.ErrTooManyConn:
		return http.StatusTooManyRequests
	case errcode.ErrMethodNotFit:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusOK
	}
}
