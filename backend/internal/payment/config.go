package payment

import (
	"ALLinSSL/backend/public"
	"time"
)

// GetPaymentConfig 获取支付配置
func GetPaymentConfig() (*PaymentConfig, error) {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	s.TableName = "payment_config"
	res, err := s.Where("id=?", []interface{}{1}).Select()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return &PaymentConfig{}, nil
	}
	row := res[0]
	cfg := &PaymentConfig{}
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
	return cfg, nil
}

// SavePaymentConfig 保存支付配置
func SavePaymentConfig(cfg *PaymentConfig) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "payment_config"
	now := time.Now().Format("2006-01-02 15:04:05")
	data := map[string]interface{}{
		"wechat_app_id":      cfg.WechatAppID,
		"wechat_mch_id":      cfg.WechatMchID,
		"wechat_api_key":     cfg.WechatAPIKey,
		"wechat_notify_url":  cfg.WechatNotifyURL,
		"alipay_app_id":      cfg.AlipayAppID,
		"alipay_private_key": cfg.AlipayPrivateKey,
		"alipay_public_key":  cfg.AlipayPublicKey,
		"alipay_notify_url":  cfg.AlipayNotifyURL,
		"update_time":        now,
	}
	// 检查是否存在
	res, err := s.Where("id=?", []interface{}{1}).Select()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		data["id"] = 1
		s2, err2 := public.NewSqlite("data/data.db", "")
		if err2 != nil {
			return err2
		}
		defer s2.Close()
		s2.TableName = "payment_config"
		_, err2 = s2.Insert(data)
		return err2
	}
	_, err = s.Where("id=?", []interface{}{1}).Update(data)
	return err
}
