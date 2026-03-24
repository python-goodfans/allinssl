package payment

// Order represents a payment order
type Order struct {
	ID            string  `json:"id"`
	OrderNo       string  `json:"order_no"`
	UserID        string  `json:"user_id"`
	PlanID        string  `json:"plan_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"` // wechat, alipay
	Status        string  `json:"status"`          // pending, paid, failed, refunded
	CreateTime    string  `json:"create_time"`
	UpdateTime    string  `json:"update_time"`
	PayTime       string  `json:"pay_time"`
	TradeNo       string  `json:"trade_no"` // 第三方交易号
}

// Plan represents a subscription plan
type Plan struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Duration    int     `json:"duration"` // 天数
	Features    string  `json:"features"` // JSON 字符串
	Status      int     `json:"status"`
	SortOrder   int     `json:"sort_order"`
	CreateTime  string  `json:"create_time"`
}

// PaymentConfig holds configuration for payment gateways
type PaymentConfig struct {
	WechatAppID      string `json:"wechat_app_id"`
	WechatMchID      string `json:"wechat_mch_id"`
	WechatAPIKey     string `json:"wechat_api_key"`
	WechatNotifyURL  string `json:"wechat_notify_url"`
	AlipayAppID      string `json:"alipay_app_id"`
	AlipayPrivateKey string `json:"alipay_private_key"`
	AlipayPublicKey  string `json:"alipay_public_key"`
	AlipayNotifyURL  string `json:"alipay_notify_url"`
	AlipaySandbox    bool   `json:"alipay_sandbox"`
}

// WechatNotifyResult holds the parsed result of a WeChat payment callback
type WechatNotifyResult struct {
	OrderNo   string
	TradeNo   string
	Status    string
	TotalFee  string
	TimeEnd   string
	ResultCode string
}

// AlipayNotifyResult holds the parsed result of an Alipay payment callback
type AlipayNotifyResult struct {
	OrderNo    string
	TradeNo    string
	TradeStatus string
}
