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
	"math/big"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	alipayGateway        = "https://openapi.alipay.com/gateway.do"
	alipaySandboxGateway = "https://openapi.alipaydev.com/gateway.do"
)

// CreateAlipayOrder creates an Alipay web payment order and returns the redirect URL
func CreateAlipayOrder(order *Order, config *PaymentConfig) (string, error) {
	if config.AlipayAppID == "" || config.AlipayPrivateKey == "" {
		return "", fmt.Errorf("支付宝配置不完整")
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

	sign, err := alipayRSA2Sign(params, config.AlipayPrivateKey)
	if err != nil {
		return "", err
	}
	params["sign"] = sign

	gateway := alipayGateway
	if config.AlipaySandbox {
		gateway = alipaySandboxGateway
	}

	return gateway + "?" + buildAlipayQuery(params), nil
}

// VerifyAlipayNotify verifies the Alipay payment callback signature
func VerifyAlipayNotify(params map[string]string, config *PaymentConfig) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	// Build the string to verify (exclude sign and sign_type)
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && k != "sign_type" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	pubKey, err := parseAlipayPublicKey(config.AlipayPublicKey)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(sb.String()))
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signBytes)
	return err == nil
}

// alipayRSA2Sign signs the Alipay request parameters using RSA2 (SHA256WithRSA)
func alipayRSA2Sign(params map[string]string, privateKeyPEM string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
	}

	privKey, err := parseAlipayPrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(sb.String()))
	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func parseAlipayPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	pemStr = strings.TrimSpace(pemStr)
	if !strings.HasPrefix(pemStr, "-----") {
		pemStr = "-----BEGIN RSA PRIVATE KEY-----\n" + pemStr + "\n-----END RSA PRIVATE KEY-----"
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("无法解析支付宝私钥")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parseAlipayPublicKey(pemStr string) (*rsa.PublicKey, error) {
	pemStr = strings.TrimSpace(pemStr)
	if !strings.HasPrefix(pemStr, "-----") {
		pemStr = "-----BEGIN PUBLIC KEY-----\n" + pemStr + "\n-----END PUBLIC KEY-----"
	}
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("无法解析支付宝公钥")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("非 RSA 公钥")
	}
	return pub, nil
}

func buildAlipayQuery(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

// generateNonce generates a cryptographically random nonce string of the given length
func generateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to time-based value if crypto/rand fails (should not happen in practice)
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			continue
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
