package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Target    string    `json:"target"` // Đối tượng nhận (tất cả/cá nhân/nhóm)
	CreatorId uuid.UUID `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}
