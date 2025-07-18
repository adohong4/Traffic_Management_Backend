package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Gov Agency model
type GovAgency struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Type      string    `json:"type"` // Đào tạo/đăng kiểm/cấp giấy tờ/...
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Version   int       `json:"version" db:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active" db:"active"`
}

// Prepare the Gov Agency for creation
func (f *GovAgency) PrepareCreate() error {
	f.Phone = strings.TrimSpace(f.Phone)
	f.Email = strings.TrimSpace(f.Email)

	f.Id = uuid.New()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.Active = true
	f.Version = 1
	return nil
}

// Prepare the Gov Agency for updating
func (f *GovAgency) PrepareUpdate() error {
	f.Phone = strings.TrimSpace(f.Phone)
	f.Email = strings.TrimSpace(f.Email)

	f.UpdatedAt = time.Now()
	f.Version++
	return nil
}

// All Gov Agency response
type GovAgencyList struct {
	TotalCount int          `json:"total_count"`
	TotalPages int          `json:"total_pages"`
	Page       int          `json:"page"`
	Size       int          `json:"size"`
	HasMore    bool         `json:"has_more"`
	GovAgency  []*GovAgency `json:"gov_agency"`
}
