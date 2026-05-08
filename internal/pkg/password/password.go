// Package password 密码哈希工具:bcrypt cost=10
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// DefaultCost bcrypt 默认代价
const DefaultCost = 10

// Hash 生成密码哈希
func Hash(plain string) (string, error) {
	if plain == "" {
		return "", fmt.Errorf("password: empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password: hash: %w", err)
	}
	return string(h), nil
}

// Verify 校验明文与哈希是否匹配
func Verify(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
