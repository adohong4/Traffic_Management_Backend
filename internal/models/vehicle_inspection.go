package models

import (
	"time"

	"github.com/google/uuid"
)

type VehicleInspection struct {
	Id             uuid.UUID `json:"id"`
	VehicleNo      string    `json:"vehicle_no"`  // Biển số xe
	ChassisNo      string    `json:"color_plate"` // Số khung
	OwnerName      string    `json:"owner_name"`
	InspectionDate string    `json:"inspection_date"` // Ngày đăng kiểm
	ExpiryDate     string    `json:"expiry_date"`     // Ngày hết hạn đăng kiểm
	Result         string    `json:"result"`          // Kết quả đăng kiểm
	Center         string    `json:"center"`          // Trung tâm đăng kiểm
	CenterId       uuid.UUID `json:"center_id"`       // Mã trung tâm đăng kiểm
	Status         string    `json:"status"`          // Trạng thái (còn hiệu lực, hết hạn...)
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
