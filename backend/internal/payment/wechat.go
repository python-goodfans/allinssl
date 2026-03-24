package payment

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
)

const wechatUnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

// wechatXMLMap 用于解析微信 XML 响应
type wechatXMLMap struct {
	XMLName xml.Name `xml:"xml"`
	Items   []struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	} `xml:",any"`
}

// CreateWechatNativeOrder 创建微信 Native 扫码支付订单，返回二维码 URL
func CreateWechatNativeOrder(order *Order, config *PaymentConfig) (string, error) {
	if config.WechatAppID == "" || config.WechatMchID == "" || config.WechatAPIKey == "" {
		return "", fmt.Errorf("微信支付配置不完整")
	}

	nonceStr := generateNonce()
	params := map[string]string{
		"appid":            config.WechatAppID,
		"mch_id":           config.WechatMchID,
		"nonce_str":        nonceStr,
		"body":             "ALLinSSL套餐购买",
		"out_trade_no":     order.OrderNo,
		"total_fee":        fmt.Sprintf("%d", int(order.Amount*100)),
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       config.WechatNotifyURL,
		"trade_type":       "NATIVE",
	}

	params["sign"] = wechatSign(params, config.WechatAPIKey)
	xmlBody := mapToXML(params)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/xml").
		SetBody(xmlBody).
		Post(wechatUnifiedOrderURL)
	if err != nil {
		return "", fmt.Errorf("请求微信支付失败: %w", err)
	}

	result := parseWechatXML(resp.Body())
	if result["return_code"] != "SUCCESS" {
		return "", fmt.Errorf("微信支付返回错误: %s", result["return_msg"])
	}
	if result["result_code"] != "SUCCESS" {
		return "", fmt.Errorf("微信支付业务错误: %s", result["err_code_des"])
	}

	codeURL, ok := result["code_url"]
	if !ok || codeURL == "" {
		return "", fmt.Errorf("未获取到微信支付二维码URL")
	}
	return codeURL, nil
}

// VerifyWechatNotify 验证微信支付回调
func VerifyWechatNotify(body []byte, config *PaymentConfig) (*WechatNotifyResult, error) {
	result := parseWechatXML(body)

	if result["return_code"] != "SUCCESS" {
		return nil, fmt.Errorf("回调状态错误")
	}

	// 验证签名
	sign := result["sign"]
	delete(result, "sign")
	expectedSign := wechatSign(result, config.WechatAPIKey)
	if !strings.EqualFold(sign, expectedSign) {
		return nil, fmt.Errorf("签名验证失败")
	}

	status := "failed"
	if result["result_code"] == "SUCCESS" {
		status = "paid"
	}

	return &WechatNotifyResult{
		OrderNo: result["out_trade_no"],
		TradeNo: result["transaction_id"],
		Status:  status,
	}, nil
}

// wechatSign 计算微信支付签名 (MD5)
func wechatSign(params map[string]string, apiKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || params[k] == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(params[k])
		sb.WriteString("&")
	}
	sb.WriteString("key=")
	sb.WriteString(apiKey)

	h := md5.New()
	h.Write([]byte(sb.String()))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// mapToXML 将 map 转换为微信XML格式
func mapToXML(params map[string]string) string {
	var sb strings.Builder
	sb.WriteString("<xml>")
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	sb.WriteString("</xml>")
	return sb.String()
}

// parseWechatXML 解析微信XML响应
func parseWechatXML(data []byte) map[string]string {
	result := make(map[string]string)
	var m wechatXMLMap
	if err := xml.Unmarshal(data, &m); err != nil {
		return result
	}
	for _, item := range m.Items {
		result[item.XMLName.Local] = item.Value
	}
	return result
}

// generateNonce 生成随机字符串
func generateNonce() string {
	h := hmac.New(md5.New, []byte("nonce"))
	h.Write([]byte(fmt.Sprintf("%d", timeNow())))
	return hex.EncodeToString(h.Sum(nil))
}
