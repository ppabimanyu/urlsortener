package repositories

import (
	"strings"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(url *models.URL) error
	FindByID(id uuid.UUID) (*models.URL, error)
	FindByIDAndUserID(id, userID uuid.UUID) (*models.URL, error)
	FindByShortCode(code string) (*models.URL, error)
	FindByUserID(userID uuid.UUID, search string, activeFilter *bool, page, limit int) ([]models.URL, int64, error)
	Update(url *models.URL) error
	Delete(id, userID uuid.UUID) error
	IncrementClickCount(id uuid.UUID) error
	ExistsByShortCode(code string) (bool, error)
	CountByUserID(userID uuid.UUID) (int64, error)
	TotalClicksByUserID(userID uuid.UUID) (int64, error)
}

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) FindByID(id uuid.UUID) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("id = ?", id).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) FindByIDAndUserID(id, userID uuid.UUID) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) FindByShortCode(code string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_code = ?", code).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) FindByUserID(userID uuid.UUID, search string, activeFilter *bool, page, limit int) ([]models.URL, int64, error) {
	var urls []models.URL
	var total int64

	query := r.db.Model(&models.URL{}).Where("user_id = ?", userID)

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(short_code) LIKE ? OR LOWER(original_url) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if activeFilter != nil {
		query = query.Where("is_active = ?", *activeFilter)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&urls).Error
	if err != nil {
		return nil, 0, err
	}

	return urls, total, nil
}

func (r *urlRepository) Update(url *models.URL) error {
	return r.db.Save(url).Error
}

func (r *urlRepository) Delete(id, userID uuid.UUID) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.URL{}).Error
}

func (r *urlRepository) IncrementClickCount(id uuid.UUID) error {
	return r.db.Model(&models.URL{}).Where("id = ?", id).UpdateColumn("click_count", gorm.Expr("click_count + 1")).Error
}

func (r *urlRepository) ExistsByShortCode(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.URL{}).Where("short_code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *urlRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.URL{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *urlRepository) TotalClicksByUserID(userID uuid.UUID) (int64, error) {
	var result struct {
		Total int64
	}
	err := r.db.Model(&models.URL{}).
		Select("COALESCE(SUM(click_count), 0) as total").
		Where("user_id = ?", userID).
		Scan(&result).Error
	return result.Total, err
}
