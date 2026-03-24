package api

import (
	"ALLinSSL/backend/internal/payment"
	"ALLinSSL/backend/public"
	"encoding/json"
	"io"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CreateOrder creates a new payment order
func CreateOrder(c *gin.Context) {
	var form struct {
		PlanID        string `form:"plan_id" binding:"required"`
		PaymentMethod string `form:"payment_method" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	session := sessions.Default(c)
	userID := ""
	if v, ok := session.Get("__username").(string); ok {
		userID = v
	}

	order, err := payment.CreateOrder(userID, form.PlanID, form.PaymentMethod)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	var payURL string
	switch form.PaymentMethod {
	case "wechat":
		payURL, err = payment.CreateWechatNativeOrder(order, cfg)
		if err != nil {
			public.FailMsg(c, "微信支付创建失败: "+err.Error())
			return
		}
	case "alipay":
		payURL, err = payment.CreateAlipayOrder(order, cfg)
		if err != nil {
			public.FailMsg(c, "支付宝支付创建失败: "+err.Error())
			return
		}
	default:
		public.FailMsg(c, "不支持的支付方式")
		return
	}

	public.SuccessData(c, map[string]interface{}{
		"order_no":       order.OrderNo,
		"pay_url":        payURL,
		"payment_method": form.PaymentMethod,
		"amount":         order.Amount,
	}, 0)
}

// GetOrderStatus returns the status of an order
func GetOrderStatus(c *gin.Context) {
	var form struct {
		OrderNo string `form:"order_no" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	order, err := payment.GetOrderByNo(form.OrderNo)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if order == nil {
		public.FailMsg(c, "订单不存在")
		return
	}

	public.SuccessData(c, order, 0)
}

// WechatNotify handles the WeChat payment callback
func WechatNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(200, "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[读取请求体失败]]></return_msg></xml>")
		return
	}

	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		c.String(200, "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[获取配置失败]]></return_msg></xml>")
		return
	}

	result, err := payment.VerifyWechatNotify(body, cfg)
	if err != nil {
		c.String(200, "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[验证失败]]></return_msg></xml>")
		return
	}

	if err := payment.UpdateOrderStatus(result.OrderNo, result.Status, result.TradeNo); err != nil {
		c.String(200, "<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[更新订单失败]]></return_msg></xml>")
		return
	}

	c.String(200, "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>")
}

// AlipayNotify handles the Alipay payment callback
func AlipayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(200, "fail")
		return
	}

	params := make(map[string]string)
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		c.String(200, "fail")
		return
	}

	if !payment.VerifyAlipayNotify(params, cfg) {
		c.String(200, "fail")
		return
	}

	orderNo := params["out_trade_no"]
	tradeNo := params["trade_no"]
	tradeStatus := params["trade_status"]

	status := "pending"
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		status = "paid"
	} else if tradeStatus == "TRADE_CLOSED" {
		status = "failed"
	}

	if err := payment.UpdateOrderStatus(orderNo, status, tradeNo); err != nil {
		c.String(200, "fail")
		return
	}

	c.String(200, "success")
}

// GetPaymentConfig returns the payment configuration
func GetPaymentConfig(c *gin.Context) {
	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	// Mask sensitive fields for display (do not expose actual keys)
	if cfg.WechatAPIKey != "" {
		cfg.WechatAPIKey = "****"
	}
	if cfg.AlipayPrivateKey != "" {
		cfg.AlipayPrivateKey = "****"
	}
	public.SuccessData(c, cfg, 0)
}

const maskedValue = "****"

// SavePaymentConfig saves the payment configuration
func SavePaymentConfig(c *gin.Context) {
	var cfg payment.PaymentConfig
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if err := json.Unmarshal(body, &cfg); err != nil {
		// Try form binding
		var form struct {
			WechatAppID      string `form:"wechat_app_id"`
			WechatMchID      string `form:"wechat_mch_id"`
			WechatAPIKey     string `form:"wechat_api_key"`
			WechatNotifyURL  string `form:"wechat_notify_url"`
			AlipayAppID      string `form:"alipay_app_id"`
			AlipayPrivateKey string `form:"alipay_private_key"`
			AlipayPublicKey  string `form:"alipay_public_key"`
			AlipayNotifyURL  string `form:"alipay_notify_url"`
			AlipaySandbox    bool   `form:"alipay_sandbox"`
		}
		if err2 := c.ShouldBind(&form); err2 != nil {
			public.FailMsg(c, err2.Error())
			return
		}
		cfg = payment.PaymentConfig{
			WechatAppID:      form.WechatAppID,
			WechatMchID:      form.WechatMchID,
			WechatAPIKey:     form.WechatAPIKey,
			WechatNotifyURL:  form.WechatNotifyURL,
			AlipayAppID:      form.AlipayAppID,
			AlipayPrivateKey: form.AlipayPrivateKey,
			AlipayPublicKey:  form.AlipayPublicKey,
			AlipayNotifyURL:  form.AlipayNotifyURL,
			AlipaySandbox:    form.AlipaySandbox,
		}
	}

	// If masked values were sent back, load existing values to preserve them
	if cfg.WechatAPIKey == maskedValue || cfg.AlipayPrivateKey == maskedValue {
		existing, err := payment.GetPaymentConfig()
		if err != nil {
			public.FailMsg(c, err.Error())
			return
		}
		if cfg.WechatAPIKey == maskedValue {
			cfg.WechatAPIKey = existing.WechatAPIKey
		}
		if cfg.AlipayPrivateKey == maskedValue {
			cfg.AlipayPrivateKey = existing.AlipayPrivateKey
		}
	}

	if err := payment.SavePaymentConfig(&cfg); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	public.SuccessMsg(c, "保存成功")
}

// GetOrderList returns a paginated list of orders
func GetOrderList(c *gin.Context) {
	var form struct {
		P     int    `form:"p"`
		Limit int    `form:"limit"`
		UserID string `form:"user_id"`
	}
	_ = c.Bind(&form)
	if form.P <= 0 {
		form.P = 1
	}
	if form.Limit <= 0 {
		form.Limit = 20
	}

	orders, total, err := payment.GetOrderList(form.UserID, form.P, form.Limit)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	public.SuccessData(c, orders, total)
}

// GetPlanList returns the list of available plans
func GetPlanList(c *gin.Context) {
	plans, err := payment.GetPlanList()
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	public.SuccessData(c, plans, len(plans))
}

// SavePlan creates or updates a plan
func SavePlan(c *gin.Context) {
	var plan payment.Plan
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if err := json.Unmarshal(body, &plan); err != nil {
		var form struct {
			ID          string  `form:"id"`
			Name        string  `form:"name" binding:"required"`
			Description string  `form:"description"`
			Price       float64 `form:"price"`
			Duration    int     `form:"duration"`
			Features    string  `form:"features"`
			Status      int     `form:"status"`
			SortOrder   int     `form:"sort_order"`
		}
		if err2 := c.ShouldBind(&form); err2 != nil {
			public.FailMsg(c, err2.Error())
			return
		}
		plan = payment.Plan{
			ID:          form.ID,
			Name:        form.Name,
			Description: form.Description,
			Price:       form.Price,
			Duration:    form.Duration,
			Features:    form.Features,
			Status:      form.Status,
			SortOrder:   form.SortOrder,
		}
	}

	if err := payment.SavePlan(&plan); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	public.SuccessMsg(c, "保存成功")
}

// DelPlan deletes a plan
func DelPlan(c *gin.Context) {
	var form struct {
		ID string `form:"id" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	if err := payment.DelPlan(form.ID); err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	public.SuccessMsg(c, "删除成功")
}
