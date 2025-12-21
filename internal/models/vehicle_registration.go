package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type VehicleRegistration struct {
	ID                uuid.UUID  `json:"id" db:"id" validate:"required"`
	OwnerID           *uuid.UUID `json:"owner_id" db:"owner_id"`                     // ID chủ sở hữu
	Brand             string     `json:"brand" db:"brand"`                           // Nhãn hiệu
	TypeVehicle       string     `json:"type_vehicle" db:"type_vehicle"`             // Loại phương tiện
	VehiclePlateNo    string     `json:"vehicle_no" db:"vehicle_no"`                 // Biển số xe
	ColorPlate        string     `json:"color_plate" db:"color_plate"`               // Màu biển số xe
	ChassisNo         string     `json:"chassis_no" db:"chassis_no"`                 // Số khung
	EngineNo          string     `json:"engine_no" db:"engine_no"`                   // Số máy
	ColorVehicle      string     `json:"color_vehicle" db:"color_vehicle"`           // Màu xe
	OwnerName         string     `json:"owner_name" db:"owner_name"`                 // Chủ xe
	Seats             *int       `json:"seats,omitempty" db:"seats"`                 // Số chỗ ngồi
	IssueDate         string     `json:"issue_date" db:"issue_date"`                 // Ngày cấp
	Issuer            string     `json:"issuer" db:"issuer"`                         // Nơi cấp
	RegistrationDate  *string    `json:"registration_date" db:"registration_date"`   // Ngay dang kiem
	ExpiryDate        *string    `json:"expiry_date" db:"expiry_date"`               // Ngày hết hạn dang kiem
	RegistrationPlace *string    `json:"registration_place" db:"registration_place"` // nơi đăng kiểm
	OnBlockchain      bool       `json:"on_blockchain" db:"on_blockchain"`           // Trạng thái lưu ở blockchain (lưa/ chưa lưu)
	BlockchainTxHash  string     `json:"blockchain_txhash" db:"blockchain_txhash"`   // Mã lưu ở blockchain
	Status            string     `json:"status" db:"status"`                         // Tình trạng (còn hiệu lực, hết hạn, bị thu hồi...)
	Version           int        `json:"version" db:"version"`                       // Phiên bản, tự động tăng
	CreatorId         uuid.UUID  `json:"creator_id" db:"creator_id" `                // ID của người tạo
	ModifierId        *uuid.UUID `json:"modifier_id" db:"modifier_id"`               // ID của người sửa
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`                 // Thời gian tạo
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`                 // Thời gian cập nhật
	Active            bool       `json:"active" db:"active"`
}

// Prepare the vehicle document for creation
func (v *VehicleRegistration) PrepareCreate() error {
	v.VehiclePlateNo = strings.TrimSpace(v.VehiclePlateNo) // eliminate spaces at the begininge and the end
	v.ChassisNo = strings.TrimSpace(v.ChassisNo)
	v.EngineNo = strings.TrimSpace(v.EngineNo)
	v.Brand = strings.TrimSpace(v.Brand)
	v.ColorPlate = strings.TrimSpace(v.ColorPlate)
	v.ColorVehicle = strings.TrimSpace(v.ColorVehicle)
	v.OwnerName = strings.TrimSpace(v.OwnerName)
	v.Issuer = strings.TrimSpace(v.Issuer)
	v.Status = strings.TrimSpace(v.Status)

	v.OnBlockchain = false
	v.BlockchainTxHash = ""
	v.ID = uuid.New()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	v.Active = true
	v.Version = 1
	return nil
}

// Prepare the vehicle document for updating
func (v *VehicleRegistration) PrepareUpdate() error {
	v.VehiclePlateNo = strings.TrimSpace(v.VehiclePlateNo)
	v.ChassisNo = strings.TrimSpace(v.ChassisNo)
	v.EngineNo = strings.TrimSpace(v.EngineNo)
	v.Brand = strings.TrimSpace(v.Brand)
	v.ColorPlate = strings.TrimSpace(v.ColorPlate)
	v.ColorVehicle = strings.TrimSpace(v.ColorVehicle)
	v.OwnerName = strings.TrimSpace(v.OwnerName)
	v.Issuer = strings.TrimSpace(v.Issuer)
	v.Status = strings.TrimSpace(v.Status)

	v.UpdatedAt = time.Now()
	v.Version++
	return nil
}

// All vehicle document response
type VehicleRegistrationList struct {
	TotalCount      int                    `json:"total_count"`
	TotalPages      int                    `json:"total_pages"`
	Page            int                    `json:"page"`
	Size            int                    `json:"size"`
	HasMore         bool                   `json:"has_more"`
	VehicleDocument []*VehicleRegistration `json:"vehicle_registration"`
}

// Confirm Blockchain Request
type ConfirmBlockchainRequest struct {
	BlockchainTxHash string `json:"blockchain_txhash" validate:"required"`
	OnBlockchain     bool   `json:"on_blockchain"`
}

// CountItem for generic count responses
type CountItem struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// VehicleTypeCounts response for type stats
type VehicleTypeCounts []*CountItem

// StatusCounts response for status stats
type StatusCounts []*CountItem

// BrandCounts response for brand stats
type BrandCounts []*CountItem
