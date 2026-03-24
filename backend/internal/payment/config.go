package payment

import (
	"ALLinSSL/backend/public"
	"time"
)

// GetPaymentConfig 获取支付配置
func GetPaymentConfig() (*PaymentConfig, error) {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()

	s.TableName = "payment_config"
	res, err := s.Where("id=?", []interface{}{1}).Select()
	if err != nil {
		return nil, err
	}

	cfg := &PaymentConfig{}
	if len(res) == 0 {
		return cfg, nil
	}

	row := res[0]
	if v, ok := row["wechat_app_id"].(string); ok {
		cfg.WechatAppID = v
	}
	if v, ok := row["wechat_mch_id"].(string); ok {
		cfg.WechatMchID = v
	}
	if v, ok := row["wechat_api_key"].(string); ok {
		cfg.WechatAPIKey = v
	}
	if v, ok := row["wechat_notify_url"].(string); ok {
		cfg.WechatNotifyURL = v
	}
	if v, ok := row["alipay_app_id"].(string); ok {
		cfg.AlipayAppID = v
	}
	if v, ok := row["alipay_private_key"].(string); ok {
		cfg.AlipayPrivateKey = v
	}
	if v, ok := row["alipay_public_key"].(string); ok {
		cfg.AlipayPublicKey = v
	}
	if v, ok := row["alipay_notify_url"].(string); ok {
		cfg.AlipayNotifyURL = v
	}
	if v, ok := row["alipay_sandbox"].(int64); ok {
		cfg.AlipaySandbox = v == 1
	}
	return cfg, nil
}

// SavePaymentConfig 保存支付配置
func SavePaymentConfig(cfg *PaymentConfig) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	sandbox := 0
	if cfg.AlipaySandbox {
		sandbox = 1
	}

	data := map[string]any{
		"id":                 1,
		"wechat_app_id":      cfg.WechatAppID,
		"wechat_mch_id":      cfg.WechatMchID,
		"wechat_api_key":     cfg.WechatAPIKey,
		"wechat_notify_url":  cfg.WechatNotifyURL,
		"alipay_app_id":      cfg.AlipayAppID,
		"alipay_private_key": cfg.AlipayPrivateKey,
		"alipay_public_key":  cfg.AlipayPublicKey,
		"alipay_notify_url":  cfg.AlipayNotifyURL,
		"alipay_sandbox":     sandbox,
		"update_time":        now,
	}

	s.TableName = "payment_config"
	// 检查是否存在
	res, err := s.Where("id=?", []interface{}{1}).Select()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		_, err = s.Insert(data)
	} else {
		_, err = s.Where("id=?", []interface{}{1}).Update(data)
	}
	return err
}
