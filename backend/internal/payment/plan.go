package payment

import (
	"ALLinSSL/backend/public"
	"time"
)

// GetPlanList 获取套餐列表
func GetPlanList() ([]Plan, error) {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	s.TableName = "plans"
	res, err := s.Order("sort_order", "ASC").Select()
	if err != nil {
		return nil, err
	}
	plans := make([]Plan, 0, len(res))
	for _, row := range res {
		plans = append(plans, *rowToPlan(row))
	}
	return plans, nil
}

// GetPlanByID 根据ID获取套餐
func GetPlanByID(id string) (*Plan, error) {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()
	s.TableName = "plans"
	res, err := s.Where("id=?", []interface{}{id}).Select()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}
	return rowToPlan(res[0]), nil
}

// AddPlan 添加套餐
func AddPlan(plan *Plan) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "plans"
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = s.Insert(map[string]interface{}{
		"id":          plan.ID,
		"name":        plan.Name,
		"description": plan.Description,
		"price":       plan.Price,
		"duration":    plan.Duration,
		"features":    plan.Features,
		"status":      plan.Status,
		"sort_order":  plan.SortOrder,
		"create_time": now,
		"update_time": now,
	})
	return err
}

// UpdatePlan 更新套餐
func UpdatePlan(plan *Plan) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "plans"
	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = s.Where("id=?", []interface{}{plan.ID}).Update(map[string]interface{}{
		"name":        plan.Name,
		"description": plan.Description,
		"price":       plan.Price,
		"duration":    plan.Duration,
		"features":    plan.Features,
		"status":      plan.Status,
		"sort_order":  plan.SortOrder,
		"update_time": now,
	})
	return err
}

// DeletePlan 删除套餐
func DeletePlan(id string) error {
	s, err := public.NewSqlite("data/data.db", "")
	if err != nil {
		return err
	}
	defer s.Close()
	s.TableName = "plans"
	_, err = s.Where("id=?", []interface{}{id}).Delete()
	return err
}

func rowToPlan(row map[string]interface{}) *Plan {
	p := &Plan{}
	if v, ok := row["id"].(string); ok {
		p.ID = v
	}
	if v, ok := row["name"].(string); ok {
		p.Name = v
	}
	if v, ok := row["description"].(string); ok {
		p.Description = v
	}
	if v, ok := row["price"].(float64); ok {
		p.Price = v
	}
	if v, ok := row["duration"].(int64); ok {
		p.Duration = int(v)
	}
	if v, ok := row["features"].(string); ok {
		p.Features = v
	}
	if v, ok := row["status"].(int64); ok {
		p.Status = int(v)
	}
	if v, ok := row["sort_order"].(int64); ok {
		p.SortOrder = int(v)
	}
	if v, ok := row["create_time"].(string); ok {
		p.CreateTime = v
	}
	if v, ok := row["update_time"].(string); ok {
		p.UpdateTime = v
	}
	return p
}
