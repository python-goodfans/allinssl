package payment

import (
	"ALLinSSL/backend/public"
	"time"
)

// CreateOrder 创建订单
func CreateOrder(order *Order) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "orders"
	_, err = s.Insert(map[string]interface{}{
		"id":          order.ID,
		"order_no":    order.OrderNo,
		"user_id":     order.UserID,
		"plan_id":     order.PlanID,
		"plan_name":   order.PlanName,
		"amount":      order.Amount,
		"pay_type":    order.PayType,
		"status":      order.Status,
		"pay_url":     order.PayURL,
		"trade_no":    order.TradeNo,
		"create_time": order.CreateTime,
		"pay_time":    order.PayTime,
	})
	return err
}

// GetOrderByNo 根据订单号查询订单
func GetOrderByNo(orderNo string) (*Order, error) {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	s.TableName = "orders"
	res, err := s.Where("order_no=?", []interface{}{orderNo}).Select()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	row := res[0]
	return rowToOrder(row), nil
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(orderNo, status, tradeNo string) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "orders"
	data := map[string]interface{}{
		"status": status,
	}
	if tradeNo != "" {
		data["trade_no"] = tradeNo
	}
	if status == "paid" {
		data["pay_time"] = time.Now().Format("2006-01-02 15:04:05")
	}
	_, err = s.Where("order_no=?", []interface{}{orderNo}).Update(data)
	return err
}

// GetOrderList 获取订单列表（分页）
func GetOrderList(userID string, page, limit int) ([]Order, int, error) {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, 0, err
	}
	defer s.Close()
	s.TableName = "orders"

	if userID != "" {
		s = s.Where("user_id=?", []interface{}{userID})
	}

	total, err := s.Count()
	if err != nil {
		return nil, 0, err
	}

	// Re-open for pagination query
	s2, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, int(total), err
	}
	defer s2.Close()
	s2.TableName = "orders"

	if userID != "" {
		s2 = s2.Where("user_id=?", []interface{}{userID})
	}

	offset := (page - 1) * limit
	s2 = s2.Order("create_time", "DESC").Limit([]int64{int64(offset), int64(limit)})
	res, err := s2.Select()
	if err != nil {
		return nil, int(total), err
	}

	orders := make([]Order, 0, len(res))
	for _, row := range res {
		orders = append(orders, *rowToOrder(row))
	}
	return orders, int(total), nil
}

func rowToOrder(row map[string]interface{}) *Order {
	o := &Order{}
	if v, ok := row["id"].(string); ok {
		o.ID = v
	}
	if v, ok := row["order_no"].(string); ok {
		o.OrderNo = v
	}
	if v, ok := row["user_id"].(string); ok {
		o.UserID = v
	}
	if v, ok := row["plan_id"].(string); ok {
		o.PlanID = v
	}
	if v, ok := row["plan_name"].(string); ok {
		o.PlanName = v
	}
	if v, ok := row["amount"].(float64); ok {
		o.Amount = v
	}
	if v, ok := row["pay_type"].(string); ok {
		o.PayType = v
	}
	if v, ok := row["status"].(string); ok {
		o.Status = v
	}
	if v, ok := row["pay_url"].(string); ok {
		o.PayURL = v
	}
	if v, ok := row["trade_no"].(string); ok {
		o.TradeNo = v
	}
	if v, ok := row["create_time"].(string); ok {
		o.CreateTime = v
	}
	if v, ok := row["pay_time"].(string); ok {
		o.PayTime = v
	}
	return o
}
