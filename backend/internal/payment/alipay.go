package payment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	alipayGateway        = "https://openapi.alipay.com/gateway.do"
	alipayGatewaySandbox = "https://openapi.alipaydev.com/gateway.do"
)

// CreateAlipayOrder 创建支付宝电脑网站支付订单，返回支付页面URL
func CreateAlipayOrder(order *Order, config *PaymentConfig) (string, error) {
	if config.AlipayAppID == "" || config.AlipayPrivateKey == "" {
		return "", fmt.Errorf("支付宝支付配置不完整")
	}

	bizContent := fmt.Sprintf(`{"out_trade_no":"%s","total_amount":"%.2f","subject":"ALLinSSL套餐购买","product_code":"FAST_INSTANT_TRADE_PAY"}`,
		order.OrderNo, order.Amount)

	params := map[string]string{
		"app_id":      config.AlipayAppID,
		"method":      "alipay.trade.page.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  config.AlipayNotifyURL,
		"biz_content": bizContent,
	}

	// 生成签名
	sign, err := alipaySign(params, config.AlipayPrivateKey)
	if err != nil {
		return "", fmt.Errorf("支付宝签名失败: %w", err)
	}
	params["sign"] = sign

	// 构建支付URL
	gateway := alipayGateway
	if config.AlipaySandbox {
		gateway = alipayGatewaySandbox
	}

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return gateway + "?" + values.Encode(), nil
}

// VerifyAlipayNotify 验证支付宝回调签名
func VerifyAlipayNotify(params map[string]string, config *PaymentConfig) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	// 构建待验签字符串
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		if params[k] != "" {
			parts = append(parts, k+"="+params[k])
		}
	}
	content := strings.Join(parts, "&")

	// RSA2 验签
	return verifyRSA2(content, sign, config.AlipayPublicKey)
}

// alipaySign 使用RSA2对参数进行签名
func alipaySign(params map[string]string, privateKeyStr string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		if params[k] != "" {
			parts = append(parts, k+"="+params[k])
		}
	}
	content := strings.Join(parts, "&")

	privateKey, err := parseRSAPrivateKey(privateKeyStr)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest)
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// verifyRSA2 验证RSA2签名
func verifyRSA2(content, sign, publicKeyStr string) bool {
	publicKey, err := parseRSAPublicKey(publicKeyStr)
	if err != nil {
		return false
	}

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	h := sha256.New()
	h.Write([]byte(content))
	digest := h.Sum(nil)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, signBytes)
	return err == nil
}

// parseRSAPrivateKey 解析RSA私钥
func parseRSAPrivateKey(keyStr string) (*rsa.PrivateKey, error) {
	keyStr = strings.TrimSpace(keyStr)
	if !strings.Contains(keyStr, "-----BEGIN") {
		keyStr = "-----BEGIN RSA PRIVATE KEY-----\n" + keyStr + "\n-----END RSA PRIVATE KEY-----"
	}

	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, fmt.Errorf("解析私钥PEM失败")
	}

	// 尝试PKCS1
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return key, nil
	}

	// 尝试PKCS8
	keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}
	rsaKey, ok := keyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA私钥")
	}
	return rsaKey, nil
}

// parseRSAPublicKey 解析RSA公钥
func parseRSAPublicKey(keyStr string) (*rsa.PublicKey, error) {
	keyStr = strings.TrimSpace(keyStr)
	if !strings.Contains(keyStr, "-----BEGIN") {
		keyStr = "-----BEGIN PUBLIC KEY-----\n" + keyStr + "\n-----END PUBLIC KEY-----"
	}

	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, fmt.Errorf("解析公钥PEM失败")
	}

	keyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	rsaKey, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA公钥")
	}
	return rsaKey, nil
}

// timeNow 返回当前时间戳（允许在测试中替换）
func timeNow() int64 {
	return time.Now().UnixNano()
}
