package user

import (
	"ALLinSSL/backend/public"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

const defaultSalt = "_bt_all_in_ssl"

// CreateUser 创建用户
func CreateUser(username, password, email string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	// 加密密码
	salt := defaultSalt
	keyMd5 := md5.Sum([]byte(password + salt))
	passwdMd5 := hex.EncodeToString(keyMd5[:])

	now := time.Now().Format("2006-01-02 15:04:05")
	s.TableName = "users"
	_, err = s.Insert(map[string]any{
		"username":    username,
		"password":    passwdMd5,
		"salt":        salt,
		"email":       email,
		"role":        "user",
		"status":      1,
		"create_time": now,
		"update_time": now,
	})
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

// GetUserByUsername 根据用户名查询用户
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
		return false, err
	}
	return u != nil, nil
}
