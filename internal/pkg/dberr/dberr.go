// Package dberr 数据库错误辅助判断。
package dberr

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// IsDuplicateKey 判断是否唯一键冲突(MySQL 1062)。
func IsDuplicateKey(err error) bool {
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}
	return false
}
