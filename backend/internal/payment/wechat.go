package payment

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
)

const wechatUnifiedOrderURL = "https://api.mch.weixin.qq.com/pay/unifiedorder"

type wechatUnifiedOrderRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Body           string `xml:"body"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
}

type wechatUnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	CodeURL    string `xml:"code_url"`
}

type wechatNotifyXML struct {
	ReturnCode  string `xml:"return_code"`
	ResultCode  string `xml:"result_code"`
	OutTradeNo  string `xml:"out_trade_no"`
	TransactionID string `xml:"transaction_id"`
	TotalFee    string `xml:"total_fee"`
	TimeEnd     string `xml:"time_end"`
}

// CreateWechatNativeOrder creates a WeChat Native (QR code) payment order and returns the code_url
func CreateWechatNativeOrder(order *Order, config *PaymentConfig) (string, error) {
	if config.WechatAppID == "" || config.WechatMchID == "" || config.WechatAPIKey == "" {
		return "", fmt.Errorf("微信支付配置不完整")
	}

	nonceStr := generateNonce(32)
	totalFee := int(order.Amount * 100) // Convert to cents

	params := map[string]string{
		"appid":            config.WechatAppID,
		"mch_id":           config.WechatMchID,
		"nonce_str":        nonceStr,
		"body":             "ALLinSSL套餐购买",
		"out_trade_no":     order.OrderNo,
		"total_fee":        fmt.Sprintf("%d", totalFee),
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       config.WechatNotifyURL,
		"trade_type":       "NATIVE",
	}

	sign := wechatSign(params, config.WechatAPIKey)
	params["sign"] = sign

	xmlBody := buildWechatXML(params)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetBody(xmlBody).
		Post(wechatUnifiedOrderURL)
	if err != nil {
		return "", err
	}

	var result wechatUnifiedOrderResponse
	if err := xml.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	if result.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("微信支付请求失败: %s", result.ReturnMsg)
	}
	if result.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("微信支付业务失败: %s - %s", result.ErrCode, result.ErrCodeDes)
	}

	return result.CodeURL, nil
}

// VerifyWechatNotify verifies the WeChat payment callback and returns the parsed result
func VerifyWechatNotify(body []byte, config *PaymentConfig) (*WechatNotifyResult, error) {
	var notify wechatNotifyXML
	if err := xml.Unmarshal(body, &notify); err != nil {
		return nil, err
	}

	if notify.ReturnCode != "SUCCESS" || notify.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("微信回调非成功状态")
	}

	return &WechatNotifyResult{
		OrderNo:    notify.OutTradeNo,
		TradeNo:    notify.TransactionID,
		Status:     "paid",
		TotalFee:   notify.TotalFee,
		TimeEnd:    notify.TimeEnd,
		ResultCode: notify.ResultCode,
	}, nil
}

// wechatSign generates the WeChat payment signature
func wechatSign(params map[string]string, apiKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" && params[k] != "" {
			keys = append(keys, k)
		}
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

	h := md5.Sum([]byte(sb.String()))
	return strings.ToUpper(hex.EncodeToString(h[:]))
}

// buildWechatXML converts a map to a WeChat XML request body
func buildWechatXML(params map[string]string) string {
	var sb strings.Builder
	sb.WriteString("<xml>")
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	sb.WriteString("</xml>")
	return sb.String()
}
