package user

import (
	"ALLinSSL/backend/public"
	"fmt"
	"time"
)

// CreateUser 创建用户记录
func CreateUser(username, passwordMd5, email, phone string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "users"
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = s.Insert(map[string]interface{}{
		"username":    username,
		"password":    passwordMd5,
		"salt":        "_bt_all_in_ssl",
		"email":       email,
		"phone":       phone,
		"create_time": now,
	})
	return err
}

// GetUserByUsername 查询用户
func GetUserByUsername(username string) (map[string]interface{}, error) {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	s.TableName = "users"
	res, err := s.Where("username=?", []interface{}{username}).Select()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res[0], nil
}

// CheckUsernameExists 检查用户名是否已存在
func CheckUsernameExists(username string) (bool, error) {
	u, err := GetUserByUsername(username)
	if err != nil {
		return false, fmt.Errorf("查询用户失败: %w", err)
	}
	return u != nil, nil
}
