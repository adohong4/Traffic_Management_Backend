package models

import (
	"time"

	"github.com/google/uuid"
)

type VehicleInsurance struct {
	Id          uuid.UUID `json:"id"`
	VehicleNo   string    `json:"vehicle_no"`
	OwnerName   string    `json:"owner_name"`
	InsuranceNo string    `json:"insurance_no"`
	Provider    string    `json:"provider"`
	IssueDate   string    `json:"issue_date"`
	ExpiryDate  string    `json:"expiry_date"`
	Type        string    `json:"type"`        // Loại bảo hiểm
	Status      string    `json:"status"`      // Còn hiệu lực, hết hạn, ...
	Version     int       `json:"version"`     // Phiên bản, tự động tăng
	CreatorId   uuid.UUID `json:"creator_id"`  // ID của người tạo
	ModifierId  uuid.UUID `json:"modifier_id"` // ID của người sửa
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
