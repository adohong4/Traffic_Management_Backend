package models

import (
	"time"

	"github.com/google/uuid"
)

type DrivingLicense struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	DOB              string    `json:"dob"`               // Ngày sinh
	IdentityNo       string    `json:"identity_no"`       // Căn cước công dân
	LicenseNo        string    `json:"license_no"`        // Số bằng lái
	IssueDate        string    `json:"issue_date"`        // Ngày cấp
	ExpiryDate       string    `json:"expiry_date"`       // Ngày hết hạn
	Status           string    `json:"status"`            // Trạng thái
	LicenseType      string    `json:"license_type"`      // Loại bằng lái
	IssuingAuthority string    `json:"issuing_authority"` // Nơi cấp
	Nationality      string    `json:"nationality"`       // Quốc tịch
	Points           int       `json:"points"`            // Điểm bằng lái xe
	Active           bool      `json:"active"`
	Version          int       `json:"version"`     // Phiên bản, tự động tăng
	CreatorId        uuid.UUID `json:"creator_id"`  // ID của người tạo
	ModifierId       uuid.UUID `json:"modifier_id"` // ID của người sửa
	CreatedAt        time.Time `json:"created_at"`  // Thời gian tạo
	UpdatedAt        time.Time `json:"updated_at"`  // Thời gian cập nhật
}
