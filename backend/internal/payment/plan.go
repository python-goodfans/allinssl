package payment

import (
	"ALLinSSL/backend/public"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GetPlanList 获取套餐列表
func GetPlanList() ([]Plan, error) {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return nil, err
	}
	defer s.Close()

	s.TableName = "plans"
	res, err := s.Where("status=?", []interface{}{1}).Order("sort_order", "ASC").Select()
	if err != nil {
		return nil, err
	}

	plans := make([]Plan, 0, len(res))
	for _, row := range res {
		plans = append(plans, rowToPlan(row))
	}
	return plans, nil
}

// SavePlan 保存/更新套餐
func SavePlan(plan *Plan) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	s.TableName = "plans"

	if plan.ID == "" {
		plan.ID = uuid.New().String()
		plan.CreateTime = now
		_, err = s.Insert(map[string]any{
			"id":          plan.ID,
			"name":        plan.Name,
			"description": plan.Description,
			"price":       plan.Price,
			"duration":    plan.Duration,
			"features":    plan.Features,
			"status":      plan.Status,
			"sort_order":  plan.SortOrder,
			"create_time": plan.CreateTime,
		})
	} else {
		_, err = s.Where("id=?", []interface{}{plan.ID}).Update(map[string]any{
			"name":        plan.Name,
			"description": plan.Description,
			"price":       plan.Price,
			"duration":    plan.Duration,
			"features":    plan.Features,
			"status":      plan.Status,
			"sort_order":  plan.SortOrder,
		})
	}
	return err
}

// DelPlan 删除套餐
func DelPlan(id string) error {
	s, err := public.NewSqlite("data/settings.db", "")
	if err != nil {
		return err
	}
	defer s.Close()

	s.TableName = "plans"
	_, err = s.Where("id=?", []interface{}{id}).Delete()
	if err != nil {
		return fmt.Errorf("删除套餐失败: %w", err)
	}
	return nil
}

func rowToPlan(row map[string]interface{}) Plan {
	p := Plan{}
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
	return p
}
