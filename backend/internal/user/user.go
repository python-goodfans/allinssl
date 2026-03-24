package user

import (
	"ALLinSSL/backend/public"
	"crypto/md5"
	"encoding/hex"
	"time"
)

const defaultSalt = "_bt_all_in_ssl"

// CreateUser creates a new user in the database
func CreateUser(username, password, email string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	s.TableName = "users"
	now := time.Now().Format("2006-01-02 15:04:05")

	keyMd5 := md5.Sum([]byte(password + defaultSalt))
	passwdMd5 := hex.EncodeToString(keyMd5[:])

	_, err = s.Insert(map[string]interface{}{
		"username":    username,
		"password":    passwdMd5,
		"salt":        defaultSalt,
		"email":       email,
		"role":        "user",
		"status":      1,
		"create_time": now,
		"update_time": now,
	})
	return err
}

// GetUserByUsername retrieves a user by username
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

// CheckUsernameExists checks whether a username already exists in the database
func CheckUsernameExists(username string) (bool, error) {
	u, err := GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return u != nil, nil
}
