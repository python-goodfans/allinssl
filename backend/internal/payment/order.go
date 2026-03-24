package payment

import (
	"ALLinSSL/backend/public"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateOrder 创建订单
func CreateOrder(userID, planID, paymentMethod string) (*Order, error) {
	// 获取套餐信息
	plans, err := GetPlanList()
	if err != nil {
		return nil, err
	}

	var plan *Plan
	for i := range plans {
		if plans[i].ID == planID {
			plan = &plans[i]
			break
		}
	}
	if plan == nil {
		return nil, fmt.Errorf("套餐不存在")
	}

	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	orderID := uuid.New().String()
	orderNo := fmt.Sprintf("ORD%d", time.Now().UnixNano())

	order := &Order{
		ID:            orderID,
		OrderNo:       orderNo,
		UserID:        userID,
		PlanID:        planID,
		Amount:        plan.Price,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		CreateTime:    now,
		UpdateTime:    now,
	}

	s.TableName = "orders"
	_, err = s.Insert(map[string]any{
		"id":             order.ID,
		"order_no":       order.OrderNo,
		"user_id":        order.UserID,
		"plan_id":        order.PlanID,
		"amount":         order.Amount,
		"payment_method": order.PaymentMethod,
		"status":         order.Status,
		"create_time":    order.CreateTime,
		"update_time":    order.UpdateTime,
	})
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}
	return order, nil
}

// GetOrderByNo 根据订单号查询订单
func GetOrderByNo(orderNo string) (*Order, error) {
	s, err := public.NewSqlite("data/settings.db", "")
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
	return rowToOrder(res[0]), nil
}

// UpdateOrderStatus 更新订单状态
func UpdateOrderStatus(orderNo, status, tradeNo string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	data := map[string]any{
		"status":      status,
		"trade_no":    tradeNo,
		"update_time": now,
	}
	if status == "paid" {
		data["pay_time"] = now
	}

	s.TableName = "orders"
	_, err = s.Where("order_no=?", []interface{}{orderNo}).Update(data)
	return err
}

// GetOrderList 获取订单列表（分页）
func GetOrderList(userID string, page, pageSize int) ([]Order, int, error) {
	s, err := public.NewSqlite("data/settings.db", "")
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

	s.TableName = "orders"
	if userID != "" {
		s = s.Where("user_id=?", []interface{}{userID})
	}

	offset := (page - 1) * pageSize
	s = s.Order("create_time", "DESC").Limit([]int64{int64(pageSize), int64(offset)})

	res, err := s.Select()
	if err != nil {
		return nil, 0, err
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
	if v, ok := row["amount"].(float64); ok {
		o.Amount = v
	}
	if v, ok := row["payment_method"].(string); ok {
		o.PaymentMethod = v
	}
	if v, ok := row["status"].(string); ok {
		o.Status = v
	}
	if v, ok := row["trade_no"].(string); ok {
		o.TradeNo = v
	}
	if v, ok := row["create_time"].(string); ok {
		o.CreateTime = v
	}
	if v, ok := row["update_time"].(string); ok {
		o.UpdateTime = v
	}
	if v, ok := row["pay_time"].(string); ok {
		o.PayTime = v
	}
	return o
}
