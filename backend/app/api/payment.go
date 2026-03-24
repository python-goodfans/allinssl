package api

import (
	"ALLinSSL/backend/internal/payment"
	"ALLinSSL/backend/public"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建支付订单
func CreateOrder(c *gin.Context) {
	var form struct {
		PlanID        string `form:"plan_id" binding:"required"`
		PaymentMethod string `form:"payment_method" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	if form.PaymentMethod != "wechat" && form.PaymentMethod != "alipay" {
		public.FailMsg(c, "不支持的支付方式")
		return
	}

	// 获取当前用户ID（简化处理，使用空字符串表示游客）
	userID := ""

	order, err := payment.CreateOrder(userID, form.PlanID, form.PaymentMethod)
	if err != nil {
		public.FailMsg(c, "创建订单失败: "+err.Error())
		return
	}

	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		public.FailMsg(c, "获取支付配置失败: "+err.Error())
		return
	}

	result := map[string]interface{}{
		"order_no":       order.OrderNo,
		"amount":         order.Amount,
		"payment_method": order.PaymentMethod,
	}

	switch form.PaymentMethod {
	case "wechat":
		codeURL, err := payment.CreateWechatNativeOrder(order, cfg)
		if err != nil {
			public.FailMsg(c, "创建微信支付失败: "+err.Error())
			return
		}
		result["code_url"] = codeURL
	case "alipay":
		payURL, err := payment.CreateAlipayOrder(order, cfg)
		if err != nil {
			public.FailMsg(c, "创建支付宝支付失败: "+err.Error())
			return
		}
		result["pay_url"] = payURL
	}

	public.SuccessData(c, result, 0)
}

// GetOrderStatus 查询订单状态
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
		public.FailMsg(c, "查询订单失败: "+err.Error())
		return
	}
	if order == nil {
		public.FailMsg(c, "订单不存在")
		return
	}

	public.SuccessData(c, order, 0)
}

// WechatNotify 微信支付回调
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

// AlipayNotify 支付宝回调
func AlipayNotify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(200, "fail")
		return
	}

	params := make(map[string]string)
	for k, vs := range c.Request.PostForm {
		if len(vs) > 0 {
			params[k] = vs[0]
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

	tradeStatus := params["trade_status"]
	status := "failed"
	if tradeStatus == "TRADE_SUCCESS" || tradeStatus == "TRADE_FINISHED" {
		status = "paid"
	}

	orderNo := params["out_trade_no"]
	tradeNo := params["trade_no"]

	if err := payment.UpdateOrderStatus(orderNo, status, tradeNo); err != nil {
		c.String(200, "fail")
		return
	}

	c.String(200, "success")
}

// GetPaymentConfig 获取支付配置
func GetPaymentConfig(c *gin.Context) {
	cfg, err := payment.GetPaymentConfig()
	if err != nil {
		public.FailMsg(c, "获取支付配置失败: "+err.Error())
		return
	}
	// 隐藏敏感信息
	safeCfg := *cfg
	if safeCfg.WechatAPIKey != "" {
		safeCfg.WechatAPIKey = "******"
	}
	if safeCfg.AlipayPrivateKey != "" {
		safeCfg.AlipayPrivateKey = "******"
	}
	public.SuccessData(c, safeCfg, 0)
}

// SavePaymentConfig 保存支付配置
func SavePaymentConfig(c *gin.Context) {
	var cfg payment.PaymentConfig
	if err := c.Bind(&cfg); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	// 如果关键字段是"******"，则保留原有值
	if cfg.WechatAPIKey == "******" || cfg.AlipayPrivateKey == "******" {
		existing, err := payment.GetPaymentConfig()
		if err == nil {
			if cfg.WechatAPIKey == "******" {
				cfg.WechatAPIKey = existing.WechatAPIKey
			}
			if cfg.AlipayPrivateKey == "******" {
				cfg.AlipayPrivateKey = existing.AlipayPrivateKey
			}
		}
	}

	if err := payment.SavePaymentConfig(&cfg); err != nil {
		public.FailMsg(c, "保存支付配置失败: "+err.Error())
		return
	}
	public.SuccessMsg(c, "保存成功")
}

// GetOrderList 获取订单列表
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
		public.FailMsg(c, "获取订单列表失败: "+err.Error())
		return
	}
	public.SuccessData(c, orders, total)
}

// GetPlanList 获取套餐列表
func GetPlanList(c *gin.Context) {
	plans, err := payment.GetPlanList()
	if err != nil {
		public.FailMsg(c, "获取套餐列表失败: "+err.Error())
		return
	}
	public.SuccessData(c, plans, len(plans))
}

// SavePlan 保存/更新套餐
func SavePlan(c *gin.Context) {
	var plan payment.Plan
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		public.FailMsg(c, err.Error())
		return
	}
	if err := json.Unmarshal(body, &plan); err != nil {
		// 尝试表单解析
		if bindErr := c.ShouldBind(&plan); bindErr != nil {
			public.FailMsg(c, "参数解析失败: "+err.Error())
			return
		}
	}

	if plan.Name == "" {
		public.FailMsg(c, "套餐名称不能为空")
		return
	}

	if err := payment.SavePlan(&plan); err != nil {
		public.FailMsg(c, "保存套餐失败: "+err.Error())
		return
	}
	public.SuccessMsg(c, "保存成功")
}

// DelPlan 删除套餐
func DelPlan(c *gin.Context) {
	var form struct {
		ID string `form:"id" binding:"required"`
	}
	if err := c.Bind(&form); err != nil {
		public.FailMsg(c, err.Error())
		return
	}

	if err := payment.DelPlan(form.ID); err != nil {
		public.FailMsg(c, "删除套餐失败: "+err.Error())
		return
	}
	public.SuccessMsg(c, "删除成功")
}
