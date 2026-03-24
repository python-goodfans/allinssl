package api

import (
	"ALLinSSL/backend/internal/user"
	"ALLinSSL/backend/public"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Register 用户注册
func Register(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Email    string `form:"email"`
		Code     string `form:"code"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	form.Username = strings.TrimSpace(form.Username)
	form.Email = strings.TrimSpace(form.Email)

	// 验证用户名格式
	if len(form.Username) < 3 || len(form.Username) > 32 {
		public.FailMsg(c, "用户名长度必须在3-32个字符之间")
		return
	}

	// 验证验证码
	session := sessions.Default(c)
	verifyCode := session.Get("_verifyCode")
	if _verifyCode, ok := verifyCode.(string); ok {
		if !strings.EqualFold(form.Code, _verifyCode) {
			public.FailMsg(c, "验证码错误")
			return
		}
	} else if form.Code != "" {
		public.FailMsg(c, "验证码已失效，请刷新重试")
		return
	}

	// 检查用户名是否已存在
	exists, err := user.CheckUsernameExists(form.Username)
	if err != nil {
		public.FailMsg(c, "检查用户名失败: "+err.Error())
		return
	}
	if exists {
		public.FailMsg(c, "用户名已存在")
		return
	}

	// 创建用户（密码已在前端用MD5加密，后端再加salt二次加密）
	if err := user.CreateUser(form.Username, form.Password, form.Email); err != nil {
		public.FailMsg(c, "注册失败: "+err.Error())
		return
	}

	// 注册成功，自动设置登录状态
	session.Set("login", true)
	session.Set("__login_key", public.LoginKey)
	_ = session.Save()

	public.SuccessMsg(c, "注册成功")
}

// CheckUsername 检查用户名是否可用
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
		public.FailMsg(c, "查询失败: "+err.Error())
		return
	}

	public.SuccessData(c, map[string]bool{"exists": exists}, 0)
}
