package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Type       string    `json:"string"`
	Target     string    `json:"target"`      // Đối tượng nhận (tất cả/cá nhân/nhóm)
	TargetUser string    `json:"target_user"` // CCCD
	Status     string    `json:"status"`

	CreatorId  uuid.UUID `json:"creator_id"`
	ModifierID uuid.UUID `json:"modifier_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Active     bool      `json:"active"`
}

func (n *Notification) PrepareCreate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.Target = strings.TrimSpace(n.Target)

	n.Id = uuid.New()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	n.Active = true
	return nil
}

func (n *Notification) PrepareUpdate() error {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.Target = strings.TrimSpace(n.Target)

	n.UpdatedAt = time.Now()
	return nil
}

type NotificationList struct {
	TotalCount   int             `json:"total_count"`
	TotalPages   int             `json:"total_pages"`
	Page         int             `json:"page"`
	Size         int             `json:"size"`
	HasMore      bool            `json:"has_more"`
	Notification []*Notification `json:"notifications"`
}
