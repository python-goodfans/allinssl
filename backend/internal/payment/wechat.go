package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

const wechatUnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

// CreateWechatPayOrder 创建微信支付订单，返回二维码URL
func CreateWechatPayOrder(order *Order, config *PaymentConfig) (string, error) {
	if config.WechatAppID == "" || config.WechatMchID == "" || config.WechatAPIKey == "" {
		return "", fmt.Errorf("微信支付配置不完整")
	}

	params := map[string]string{
		"appid":            config.WechatAppID,
		"mch_id":           config.WechatMchID,
		"nonce_str":        order.ID,
		"body":             order.PlanName,
		"out_trade_no":     order.OrderNo,
		"total_fee":        fmt.Sprintf("%d", int(order.Amount*100)),
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       config.WechatNotifyURL,
		"trade_type":       "NATIVE",
	}
	params["sign"] = wechatSign(params, config.WechatAPIKey)

	xmlBody := mapToXML(params)
	resp, err := http.Post(wechatUnifiedOrderURL, "application/xml", strings.NewReader(xmlBody))
	if err != nil {
		return "", fmt.Errorf("微信支付请求失败: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取微信响应失败: %w", err)
	}

	result, err := parseXML(body)
	if err != nil {
		return "", fmt.Errorf("解析微信响应失败: %w", err)
	}
	if result["return_code"] != "SUCCESS" {
		return "", fmt.Errorf("微信下单失败: %s", result["return_msg"])
	}
	if result["result_code"] != "SUCCESS" {
		return "", fmt.Errorf("微信业务失败: %s", result["err_code_des"])
	}
	codeURL, ok := result["code_url"]
	if !ok || codeURL == "" {
		return "", fmt.Errorf("未获取到微信支付二维码URL")
	}
	return codeURL, nil
}

// VerifyWechatNotify 验证微信回调签名
func VerifyWechatNotify(body []byte, apiKey string) (map[string]string, error) {
	params, err := parseXML(body)
	if err != nil {
		return nil, err
	}
	sign, ok := params["sign"]
	if !ok {
		return nil, fmt.Errorf("缺少sign参数")
	}
	delete(params, "sign")
	expected := wechatSign(params, apiKey)
	if !strings.EqualFold(sign, expected) {
		return nil, fmt.Errorf("签名验证失败")
	}
	return params, nil
}

func wechatSign(params map[string]string, apiKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	str := strings.Join(parts, "&") + "&key=" + apiKey
	sum := md5.Sum([]byte(str))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}

func mapToXML(params map[string]string) string {
	var sb strings.Builder
	sb.WriteString("<xml>")
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	sb.WriteString("</xml>")
	return sb.String()
}

type xmlKV struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func parseXML(data []byte) (map[string]string, error) {
	result := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	var currentKey string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local != "xml" {
				currentKey = t.Name.Local
			}
		case xml.CharData:
			if currentKey != "" {
				result[currentKey] = strings.TrimSpace(string(t))
			}
		case xml.EndElement:
			currentKey = ""
		}
	}
	return result, nil
}
