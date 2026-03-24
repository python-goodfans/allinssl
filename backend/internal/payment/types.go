package payment

// Order 订单
type Order struct {
	ID         string  `json:"id"`
	OrderNo    string  `json:"order_no"`
	UserID     string  `json:"user_id"`
	PlanID     string  `json:"plan_id"`
	PlanName   string  `json:"plan_name"`
	Amount     float64 `json:"amount"`
	PayType    string  `json:"pay_type"`   // wechat / alipay
	Status     string  `json:"status"`     // pending / paid / failed / expired
	PayURL     string  `json:"pay_url"`    // 支付链接或二维码URL
	TradeNo    string  `json:"trade_no"`   // 第三方交易号
	CreateTime string  `json:"create_time"`
	PayTime    string  `json:"pay_time"`
}

// Plan 套餐
type Plan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"`    // 有效天数
	Features    string  `json:"features"`    // JSON格式的功能列表
	Status      int     `json:"status"`      // 0: 禁用, 1: 启用
	SortOrder   int     `json:"sort_order"`
	CreateTime  string  `json:"create_time"`
	UpdateTime  string  `json:"update_time"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	WechatAppID      string `json:"wechat_app_id"`
	WechatMchID      string `json:"wechat_mch_id"`
	WechatAPIKey     string `json:"wechat_api_key"`
	WechatNotifyURL  string `json:"wechat_notify_url"`
	AlipayAppID      string `json:"alipay_app_id"`
	AlipayPrivateKey string `json:"alipay_private_key"`
	AlipayPublicKey  string `json:"alipay_public_key"`
	AlipayNotifyURL  string `json:"alipay_notify_url"`
}
