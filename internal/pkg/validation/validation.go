// Package validation 项目级自定义 validator。
// 只注册一次(sync.Once),在 bootstrap 初始化阶段调用。
package validation

import (
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	registerOnce sync.Once
	slugRe       = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

// Register 向 gin 使用的 validator 实例注册自定义 tag
func Register() error {
	var err error
	registerOnce.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}
		err = v.RegisterValidation("slug", validateSlug)
	})
	return err
}

// validateSlug 只允许小写字母、数字和连字符(`-`),不能以 `-` 开头或结尾。
func validateSlug(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if s == "" {
		return false
	}
	return slugRe.MatchString(s)
}
