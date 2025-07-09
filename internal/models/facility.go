package models

import (
	"time"

	"github.com/google/uuid"
)

type Facility struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Type      string    `json:"type"` // Đào tạo/đăng kiểm/cấp giấy tờ/...
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
