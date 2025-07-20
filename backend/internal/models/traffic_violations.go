package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type TrafficViolation struct {
	Id             uuid.UUID  `json:"id" db:"id"`
	VehiclePlateNo string     `json:"vehicle_no" db:"vehicle_no"`   // Biển số xe
	Date           time.Time  `json:"date" db:"date"`               // Ngày vi phạm
	Type           string     `json:"type" db:"type"`               // Loại vi phạm
	Description    string     `json:"description" db:"description"` // Mô tả vi phạm
	Points         int        `json:"points" db:"points"`           // Số điểm bị trừ
	FineAmount     int64      `json:"fine_amount" db:"fine_amount"` // Số tiền phạt (VND)
	Status         string     `json:"status" db:"status"`           // Trạng thái (đã xử lý/chưa xử lý/hủy vi phạm)
	Version        int        `json:"version" db:"version"`         // Phiên bản, tự động tăng
	CreatorId      uuid.UUID  `json:"creator_id" db:"creator_id"`   // ID của người tạo
	ModifierId     *uuid.UUID `json:"modifier_id" db:"modifier_id"` // ID của người sửa
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`   // Thời gian tạo
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`   // Thời gian cập nhật
	Active         bool       `json:"active" db:"active"`
}

// Prepare the traffic violation for creation
func (t *TrafficViolation) PrepareCreate() error {
	t.VehiclePlateNo = strings.TrimSpace(t.VehiclePlateNo)
	t.Type = strings.TrimSpace(t.Type)
	t.Description = strings.TrimSpace(t.Description)
	t.Status = strings.TrimSpace(t.Status)

	t.Id = uuid.New()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.Active = true
	t.Version = 1
	return nil
}

// Prepare the traffic violation for updating
func (t *TrafficViolation) PrepareUpdate() error {
	t.VehiclePlateNo = strings.TrimSpace(t.VehiclePlateNo)
	t.Type = strings.TrimSpace(t.Type)
	t.Description = strings.TrimSpace(t.Description)
	t.Status = strings.TrimSpace(t.Status)

	t.UpdatedAt = time.Now()
	return nil
}

// All traffic violation response
type TrafficViolationList struct {
	TotalCount       int                 `json:"total_count"`
	TotalPages       int                 `json:"total_pages"`
	Page             int                 `json:"page"`
	Size             int                 `json:"size"`
	HasMore          bool                `json:"has_more"`
	TrafficViolation []*TrafficViolation `json:"traffic_violation"`
}
