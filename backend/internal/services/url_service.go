package services

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidURL          = errors.New("invalid target URL format")
	ErrSlugAlreadyExists   = errors.New("this custom slug is already in use")
	ErrInvalidSlug         = errors.New("custom slug must be 3-50 alphanumeric characters (or '-' and '_')")
	ErrURLNotFound         = errors.New("URL not found")
	ErrURLExpired          = errors.New("this URL has expired")
	ErrURLInactive         = errors.New("this URL is currently disabled")
	ErrClickLimitReached   = errors.New("this URL has reached its maximum click limit")
	ErrInvalidURLPassword  = errors.New("incorrect password for this URL")
	ErrPasswordRequired    = errors.New("password required to access this URL")
)

type URLService interface {
	CreateURL(userID uuid.UUID, originalURL, customSlug, title, password string, expiresAt *time.Time, clickLimit *int) (*models.URL, error)
	GetURLs(userID uuid.UUID, search string, activeFilter *bool, page, limit int) ([]models.URL, int64, error)
	GetURLByID(id, userID uuid.UUID) (*models.URL, error)
	UpdateURL(id, userID uuid.UUID, originalURL, customSlug, title string, password *string, expiresAt *time.Time, clickLimit *int, isActive *bool) (*models.URL, error)
	DeleteURL(id, userID uuid.UUID) error
	ToggleStatus(id, userID uuid.UUID) (*models.URL, error)
	ResolveURL(code string) (*models.URL, bool, error)
	VerifyURLPassword(code, password string) (*models.URL, error)
	RecordClick(urlID uuid.UUID, click *models.URLClick)
}

type urlService struct {
	urlRepo   repositories.URLRepository
	clickRepo repositories.ClickRepository
}

func NewURLService(urlRepo repositories.URLRepository, clickRepo repositories.ClickRepository) URLService {
	return &urlService{
		urlRepo:   urlRepo,
		clickRepo: clickRepo,
	}
}

func (s *urlService) CreateURL(userID uuid.UUID, originalURL, customSlug, title, password string, expiresAt *time.Time, clickLimit *int) (*models.URL, error) {
	originalURL = strings.TrimSpace(originalURL)
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = "https://" + originalURL
	}

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || parsedURL.Host == "" {
		return nil, ErrInvalidURL
	}

	var slug string
	customSlug = strings.TrimSpace(customSlug)
	if customSlug != "" {
		if !utils.IsValidSlug(customSlug) {
			return nil, ErrInvalidSlug
		}
		exists, err := s.urlRepo.ExistsByShortCode(customSlug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSlugAlreadyExists
		}
		slug = customSlug
	} else {
		// Generate unique short code with retry
		for i := 0; i < 5; i++ {
			code, err := utils.GenerateShortCode(6)
			if err != nil {
				return nil, err
			}
			exists, err := s.urlRepo.ExistsByShortCode(code)
			if err != nil {
				return nil, err
			}
			if !exists {
				slug = code
				break
			}
		}
		if slug == "" {
			return nil, errors.New("failed to generate unique slug, please try again")
		}
	}

	if title == "" {
		title = parsedURL.Host
	}

	var passwordHash *string
	if password != "" {
		hashed, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
		passwordHash = &hashed
	}

	newURL := &models.URL{
		UserID:       userID,
		OriginalURL:  originalURL,
		ShortCode:    slug,
		Title:        title,
		PasswordHash: passwordHash,
		ExpiresAt:    expiresAt,
		ClickLimit:   clickLimit,
		IsActive:     true,
	}

	if err := s.urlRepo.Create(newURL); err != nil {
		return nil, err
	}

	newURL.HasPassword = passwordHash != nil
	return newURL, nil
}

func (s *urlService) GetURLs(userID uuid.UUID, search string, activeFilter *bool, page, limit int) ([]models.URL, int64, error) {
	return s.urlRepo.FindByUserID(userID, search, activeFilter, page, limit)
}

func (s *urlService) GetURLByID(id, userID uuid.UUID) (*models.URL, error) {
	urlItem, err := s.urlRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}
	return urlItem, nil
}

func (s *urlService) UpdateURL(id, userID uuid.UUID, originalURL, customSlug, title string, password *string, expiresAt *time.Time, clickLimit *int, isActive *bool) (*models.URL, error) {
	existingURL, err := s.urlRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}

	if originalURL != "" {
		originalURL = strings.TrimSpace(originalURL)
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "https://" + originalURL
		}
		parsedURL, err := url.ParseRequestURI(originalURL)
		if err != nil || parsedURL.Host == "" {
			return nil, ErrInvalidURL
		}
		existingURL.OriginalURL = originalURL
	}

	if customSlug != "" && customSlug != existingURL.ShortCode {
		if !utils.IsValidSlug(customSlug) {
			return nil, ErrInvalidSlug
		}
		exists, err := s.urlRepo.ExistsByShortCode(customSlug)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSlugAlreadyExists
		}
		existingURL.ShortCode = customSlug
	}

	if title != "" {
		existingURL.Title = title
	}

	if password != nil {
		if *password == "" {
			existingURL.PasswordHash = nil
		} else {
			hashed, err := utils.HashPassword(*password)
			if err != nil {
				return nil, err
			}
			existingURL.PasswordHash = &hashed
		}
	}

	existingURL.ExpiresAt = expiresAt
	existingURL.ClickLimit = clickLimit

	if isActive != nil {
		existingURL.IsActive = *isActive
	}

	if err := s.urlRepo.Update(existingURL); err != nil {
		return nil, err
	}

	existingURL.HasPassword = existingURL.PasswordHash != nil
	return existingURL, nil
}

func (s *urlService) DeleteURL(id, userID uuid.UUID) error {
	return s.urlRepo.Delete(id, userID)
}

func (s *urlService) ToggleStatus(id, userID uuid.UUID) (*models.URL, error) {
	existingURL, err := s.urlRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	existingURL.IsActive = !existingURL.IsActive
	if err := s.urlRepo.Update(existingURL); err != nil {
		return nil, err
	}

	existingURL.HasPassword = existingURL.PasswordHash != nil
	return existingURL, nil
}

func (s *urlService) ResolveURL(code string) (*models.URL, bool, error) {
	urlItem, err := s.urlRepo.FindByShortCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, ErrURLNotFound
		}
		return nil, false, err
	}

	if !urlItem.IsActive {
		return nil, false, ErrURLInactive
	}

	if urlItem.ExpiresAt != nil && time.Now().After(*urlItem.ExpiresAt) {
		return nil, false, ErrURLExpired
	}

	if urlItem.ClickLimit != nil && *urlItem.ClickLimit > 0 && urlItem.ClickCount >= int64(*urlItem.ClickLimit) {
		return nil, false, ErrClickLimitReached
	}

	// Check if password required
	if urlItem.PasswordHash != nil && *urlItem.PasswordHash != "" {
		return urlItem, true, nil
	}

	return urlItem, false, nil
}

func (s *urlService) VerifyURLPassword(code, password string) (*models.URL, error) {
	urlItem, err := s.urlRepo.FindByShortCode(code)
	if err != nil {
		return nil, ErrURLNotFound
	}

	if !urlItem.IsActive {
		return nil, ErrURLInactive
	}

	if urlItem.ExpiresAt != nil && time.Now().After(*urlItem.ExpiresAt) {
		return nil, ErrURLExpired
	}

	if urlItem.ClickLimit != nil && *urlItem.ClickLimit > 0 && urlItem.ClickCount >= int64(*urlItem.ClickLimit) {
		return nil, ErrClickLimitReached
	}

	if urlItem.PasswordHash == nil || *urlItem.PasswordHash == "" {
		return urlItem, nil
	}

	if !utils.CheckPasswordHash(password, *urlItem.PasswordHash) {
		return nil, ErrInvalidURLPassword
	}

	return urlItem, nil
}

func (s *urlService) RecordClick(urlID uuid.UUID, click *models.URLClick) {
	go func() {
		_ = s.urlRepo.IncrementClickCount(urlID)
		if click != nil {
			click.URLID = urlID
			_ = s.clickRepo.Create(click)
		}
	}()
}
