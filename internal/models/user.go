package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	IdentityNo string    `json:"identity_no"` // Căn cước công dân
	Password   string    `json:"password"`
	Active     bool      `json:"active"`
	Version    int       `json:"version"`     // Phiên bản, tự động tăng
	CreatorId  uuid.UUID `json:"creator_id"`  // ID của người tạo
	ModifierId uuid.UUID `json:"modifier_id"` // ID của người sửa
	CreatedAt  time.Time `json:"created_at"`  // Thời gian tạo
	UpdatedAt  time.Time `json:"updated_at"`  // Thời gian cập nhật
}
