package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uuid.UUID  `json:"id" db:"id" validate:"required"`
	IdentityNo   string     `json:"identity_no" db:"identity_no" validate:"required,lte=20"`     // CCCD
	Password     string     `json:"password,omitempty" db:"password" validate:"omitempty,gte=6"` // Mật khẩu truyền thống
	HashPassword string     `json:"hash_password" db:"hash_password" validate:"required"`
	Active       bool       `json:"active" db:"active"`
	Role         *string    `json:"role,omitempty" db:"role" validate:"omitempty,lte=20"` // Vai trò (admin, user, etc.)
	Version      int        `json:"version" db:"version"`                                 // Phiên bản, tự động tăng
	CreatorId    *uuid.UUID `json:"creator_id" db:"creator_id"`                           // ID của người tạo
	ModifierId   *uuid.UUID `json:"modifier_id" db:"modifier_id"`                         // ID của người sửa
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`                           // Thời gian tạo
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`                           // Thời gian cập nhật
}

// Hash user password with bcrypt
func (u *User) HashedPassword() error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashPassword = string(hashPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
	u.HashPassword = ""
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Password = strings.TrimSpace(u.Password)
	u.IdentityNo = strings.TrimSpace(u.IdentityNo)

	if err := u.HashedPassword(); err != nil {
		return err
	}

	u.Id = uuid.New()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.Active = true
	u.Version = 1
	return nil
}

// PrepareUpdate prepares a user for update
func (u *User) PrepareUpdate() error {
	u.IdentityNo = strings.TrimSpace(u.IdentityNo)
	u.Password = strings.TrimSpace(u.Password)

	if u.Password != "" {
		if err := u.HashedPassword(); err != nil {
			return err
		}
	}

	u.UpdatedAt = time.Now()
	u.Version++
	return nil
}

// All Users reponse
type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

// Find user query
type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
