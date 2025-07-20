package models

import (
	"time"

	"github.com/google/uuid"
)

type OtherDocument struct {
	Id           uuid.UUID `json:"id"`
	VehicleNo    string    `json:"vehicle_no"`
	OwnerName    string    `json:"owner_name"`
	DocumentType string    `json:"document_type"` // Đăng ký xe, phí đường bộ, giấy phép kinh doanh vận tải, ...
	DocumentNo   string    `json:"document_no"`
	DocumentURL  string    `json:"document_url"`
	IssueDate    string    `json:"issue_date"`
	ExpiryDate   string    `json:"expiry_date"`
	Issuer       string    `json:"issuer"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
