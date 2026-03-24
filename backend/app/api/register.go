package api

import (
	"ALLinSSL/backend/internal/user"
	"ALLinSSL/backend/public"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

func Register(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Email    string `form:"email"`
		Phone    string `form:"phone"`
		Code     string `form:"code"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	form.Username = strings.TrimSpace(form.Username)

	// 检查用户名是否已存在
	exists, err := user.CheckUsernameExists(form.Username)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if exists {
		public.FailMsg(c, "用户名已存在")
		return
	}

	// 密码 MD5 + salt 加密
	salt := "_bt_all_in_ssl"
	keyMd5 := md5.Sum([]byte(form.Password + salt))
	passwordMd5 := hex.EncodeToString(keyMd5[:])

	// 创建用户
	if err = user.CreateUser(form.Username, passwordMd5, form.Email, form.Phone); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	// 注册成功后自动设置 session 登录状态
	session := sessions.Default(c)
	session.Set("login", true)
	session.Set("__login_key", public.LoginKey)
	_ = session.Save()

	public.SuccessMsg(c, "注册成功")
}

func CheckUsername(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	form.Username = strings.TrimSpace(form.Username)

	exists, err := user.CheckUsernameExists(form.Username)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	// available = !exists
	public.SuccessData(c, !exists, 0)
}
