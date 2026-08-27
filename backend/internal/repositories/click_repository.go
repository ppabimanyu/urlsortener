package repositories

import (
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DateClickStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type GroupCountStat struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type ClickRepository interface {
	Create(click *models.URLClick) error
	GetClicksByDate(urlID uuid.UUID, days int) ([]DateClickStat, error)
	GetDeviceStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetOSStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetBrowserStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetReferrerStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetCountryStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetCityStats(urlID uuid.UUID) ([]GroupCountStat, error)
	GetAllClicksForExport(urlID uuid.UUID) ([]models.URLClick, error)
	GetOverviewTimeline(userID uuid.UUID, days int) ([]DateClickStat, error)
	GetRecentClicksByUserID(userID uuid.UUID, limit int) ([]models.URLClick, error)
}

type clickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) ClickRepository {
	return &clickRepository{db: db}
}

func (r *clickRepository) Create(click *models.URLClick) error {
	return r.db.Create(click).Error
}

func (r *clickRepository) GetClicksByDate(urlID uuid.UUID, days int) ([]DateClickStat, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)

	var stats []DateClickStat
	err := r.db.Model(&models.URLClick{}).
		Select("TO_CHAR(clicked_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("url_id = ? AND clicked_at >= ?", urlID, since).
		Group("TO_CHAR(clicked_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&stats).Error

	return stats, err
}

func (r *clickRepository) GetDeviceStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(device_type, ''), 'Unknown') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("device_type").
		Order("count DESC").
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetOSStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(os, ''), 'Unknown') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("os").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetBrowserStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(browser, ''), 'Unknown') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("browser").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetReferrerStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(referrer, ''), 'Direct') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("referrer").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetCountryStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(country, ''), 'Unknown') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetCityStats(urlID uuid.UUID) ([]GroupCountStat, error) {
	var stats []GroupCountStat
	err := r.db.Model(&models.URLClick{}).
		Select("COALESCE(NULLIF(city, ''), 'Unknown') as name, COUNT(*) as count").
		Where("url_id = ?", urlID).
		Group("city").
		Order("count DESC").
		Limit(10).
		Scan(&stats).Error
	return stats, err
}

func (r *clickRepository) GetAllClicksForExport(urlID uuid.UUID) ([]models.URLClick, error) {
	var clicks []models.URLClick
	err := r.db.Where("url_id = ?", urlID).Order("clicked_at DESC").Find(&clicks).Error
	return clicks, err
}

func (r *clickRepository) GetOverviewTimeline(userID uuid.UUID, days int) ([]DateClickStat, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)

	var stats []DateClickStat
	err := r.db.Table("url_clicks").
		Joins("JOIN urls ON urls.id = url_clicks.url_id").
		Select("TO_CHAR(url_clicks.clicked_at, 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("urls.user_id = ? AND url_clicks.clicked_at >= ?", userID, since).
		Group("TO_CHAR(url_clicks.clicked_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&stats).Error

	return stats, err
}

func (r *clickRepository) GetRecentClicksByUserID(userID uuid.UUID, limit int) ([]models.URLClick, error) {
	if limit <= 0 {
		limit = 10
	}
	var clicks []models.URLClick
	err := r.db.Table("url_clicks").
		Joins("JOIN urls ON urls.id = url_clicks.url_id").
		Where("urls.user_id = ?", userID).
		Order("url_clicks.clicked_at DESC").
		Limit(limit).
		Find(&clicks).Error
	return clicks, err
}
