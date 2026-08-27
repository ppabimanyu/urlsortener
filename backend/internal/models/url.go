package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URL struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	OriginalURL  string         `gorm:"type:text;not null" json:"original_url"`
	ShortCode    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"short_code"`
	Title        string         `gorm:"type:varchar(255)" json:"title"`
	PasswordHash *string        `gorm:"type:varchar(255)" json:"-"`
	HasPassword  bool           `gorm:"-" json:"has_password"`
	ExpiresAt    *time.Time     `gorm:"type:timestamptz" json:"expires_at"`
	ClickLimit   *int           `gorm:"default:null" json:"click_limit"`
	ClickCount   int64          `gorm:"default:0;not null" json:"click_count"`
	IsActive     bool           `gorm:"default:true;not null" json:"is_active"`
	Clicks       []URLClick     `gorm:"foreignKey:URLID;constraint:OnDelete:CASCADE" json:"clicks,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *URL) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u *URL) AfterFind(tx *gorm.DB) (err error) {
	u.HasPassword = u.PasswordHash != nil && *u.PasswordHash != ""
	return
}
