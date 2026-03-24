package payment

import (
	"ALLinSSL/backend/public"
	"time"

	"github.com/google/uuid"
)

// CreateOrder creates a new payment order in the database
func CreateOrder(userID, planID, paymentMethod string) (*Order, error) {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()

	// Get plan amount
	s.TableName = "plans"
	plans, err := s.Where("id=?", []interface{}{planID}).Select()
	if err != nil {
		return nil, err
	}

	var amount float64
	if len(plans) > 0 {
		if v, ok := plans[0]["price"]; ok {
			switch val := v.(type) {
			case float64:
				amount = val
			case int64:
				amount = float64(val)
			}
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	orderID := uuid.New().String()
	orderNo := "ORDER" + time.Now().Format("20060102150405") + public.RandomString(6)

	order := &Order{
		ID:            orderID,
		OrderNo:       orderNo,
		UserID:        userID,
		PlanID:        planID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		CreateTime:    now,
		UpdateTime:    now,
	}

	s2, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s2.Close()
	s2.TableName = "orders"

	_, err = s2.Insert(map[string]interface{}{
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
		return nil, err
	}

	return order, nil
}

// GetOrderByNo retrieves an order by its order number
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

	return mapToOrder(res[0]), nil
}

// UpdateOrderStatus updates the status and trade number of an order
func UpdateOrderStatus(orderNo, status, tradeNo string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	s.TableName = "orders"
	now := time.Now().Format("2006-01-02 15:04:05")
	data := map[string]interface{}{
		"status":      status,
		"trade_no":    tradeNo,
		"update_time": now,
	}
	if status == "paid" {
		data["pay_time"] = now
	}
	_, err = s.Where("order_no=?", []interface{}{orderNo}).Update(data)
	return err
}

// GetOrderList returns a paginated list of orders for a user
func GetOrderList(userID string, page, pageSize int) ([]Order, int, error) {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, 0, err
	}
	defer s.Close()

	s.TableName = "orders"
	where := "1=1"
	params := []interface{}{}
	if userID != "" {
		where += " AND user_id=?"
		params = append(params, userID)
	}

	countRes, err := s.Where(where, params).Count()
	if err != nil {
		return nil, 0, err
	}

	offset := int64((page - 1) * pageSize)
	res, err := s.Where(where, params).Order("create_time", "DESC").Limit([]int64{offset, int64(pageSize)}).Select()
	if err != nil {
		return nil, 0, err
	}

	orders := make([]Order, 0, len(res))
	for _, row := range res {
		orders = append(orders, *mapToOrder(row))
	}
	return orders, int(countRes), nil
}

func mapToOrder(row map[string]interface{}) *Order {
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
