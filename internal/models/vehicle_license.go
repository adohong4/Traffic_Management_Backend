package models

import (
	"time"

	"github.com/google/uuid"
)

type VehicleDocument struct {
	ID           uuid.UUID `json:"id"`
	Brand        string    `json:"brand"`         // Nhãn hiệu
	VehicleNo    string    `json:"vehicle_no"`    // Biển số xe
	ColorPlate   string    `json:"color_plate"`   // Màu biển số xe
	ChassisNo    string    `json:"chassis_no"`    // Số khung
	EngineNo     string    `json:"engine_no"`     // Số máy
	ColorVehicle string    `json:"color_vehicle"` // Màu xe
	OwnerName    string    `json:"owner_name"`    // Chủ xe
	DocumentType string    `json:"document_type"` // Loại giấy tờ (đăng ký, đăng kiểm, bảo hiểm, ...)
	DocumentNo   string    `json:"document_no"`   // Số giấy tờ
	IssueDate    string    `json:"issue_date"`    // Ngày cấp
	ExpiryDate   string    `json:"expiry_date"`   // Ngày hết hạn
	Issuer       string    `json:"issuer"`        // Nơi cấp
	Status       string    `json:"status"`        // Tình trạng (còn hiệu lực, hết hạn, bị thu hồi...)
	Version      int       `json:"version"`       // Phiên bản, tự động tăng
	CreatorId    uuid.UUID `json:"creator_id"`    // ID của người tạo
	ModifierId   uuid.UUID `json:"modifier_id"`   // ID của người sửa
	CreatedAt    time.Time `json:"created_at"`    // Thời gian tạo
	UpdatedAt    time.Time `json:"updated_at"`    // Thời gian cập nhật
}
