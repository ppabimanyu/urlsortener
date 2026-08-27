package models

import (
	"time"

	"github.com/google/uuid"
)

type URLClick struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	URLID       uuid.UUID `gorm:"type:uuid;not null;index" json:"url_id"`
	ClickedAt   time.Time `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP;index" json:"clicked_at"`
	IPAddress   string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	DeviceType  string    `gorm:"type:varchar(50)" json:"device_type"`
	OS          string    `gorm:"type:varchar(50)" json:"os"`
	Browser     string    `gorm:"type:varchar(50)" json:"browser"`
	Referrer    string    `gorm:"type:varchar(255)" json:"referrer"`
	Country     string    `gorm:"type:varchar(100)" json:"country"`
	CountryCode string    `gorm:"type:varchar(10)" json:"country_code"`
	City        string    `gorm:"type:varchar(100)" json:"city"`
}
