package api

import (
	"ALLinSSL/backend/internal/user"
	"ALLinSSL/backend/public"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

// Register handles user registration
func Register(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Email    string `form:"email" binding:"required"`
		Code     string `form:"code" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	form.Username = strings.TrimSpace(form.Username)
	form.Email = strings.TrimSpace(form.Email)
	form.Code = strings.TrimSpace(form.Code)

	// Verify captcha code
	session := sessions.Default(c)
	verifyCode := session.Get("_verifyCode")
	if _verifyCode, ok := verifyCode.(string); ok {
		if !strings.EqualFold(form.Code, _verifyCode) {
			public.FailMsg(c, "验证码错误")
			return
		}
	} else {
		public.FailMsg(c, "验证码已过期，请重新获取")
		return
	}

	// Check if username already exists
	exists, err := user.CheckUsernameExists(form.Username)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if exists {
		public.FailMsg(c, "用户名已存在")
		return
	}

	// Create the user
	if err := user.CreateUser(form.Username, form.Password, form.Email); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	// Auto login after registration
	session.Set("login", true)
	session.Set("__login_key", public.LoginKey)
	_ = session.Save()

	public.SuccessMsg(c, "注册成功")
}

// CheckUsername checks whether a username is available
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

	public.SuccessData(c, map[string]interface{}{
		"exists": exists,
	}, 0)
}
