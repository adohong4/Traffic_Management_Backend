package models

import (
	"time"

	"github.com/google/uuid"
)

type TrafficViolation struct {
	Id          uuid.UUID `json:"id"`
	LicenseId   uuid.UUID `json:"license_id"` // Liên kết đến License
	Date        string    `json:"date"`       // Ngày vi phạm
	Type        string    `json:"type"`       // Loại vi phạm
	Description string    `json:"description"`
	Points      int       `json:"points"`      // Số điểm bị trừ
	FineAmount  int64     `json:"fine_amount"` // Số tiền phạt (VND)
	Status      string    `json:"status"`      // Trạng thái (đã xử lý/chưa xử lý/hủy vi phạm)
	Version     int       `json:"version"`     // Phiên bản, tự động tăng
	CreatorId   uuid.UUID `json:"creator_id"`  // ID của người tạo
	ModifierId  uuid.UUID `json:"modifier_id"` // ID của người sửa
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
