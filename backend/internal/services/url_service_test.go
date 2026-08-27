package services

import (
	"errors"
	"testing"
	"time"

	"backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockURLRepo struct {
	urls map[uuid.UUID]*models.URL
}

func newMockURLRepo() *mockURLRepo {
	return &mockURLRepo{urls: make(map[uuid.UUID]*models.URL)}
}

func (m *mockURLRepo) Create(url *models.URL) error {
	if url.ID == uuid.Nil {
		url.ID = uuid.New()
	}
	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()
	m.urls[url.ID] = url
	return nil
}

func (m *mockURLRepo) FindByID(id uuid.UUID) (*models.URL, error) {
	url, ok := m.urls[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return url, nil
}

func (m *mockURLRepo) FindByIDAndUserID(id, userID uuid.UUID) (*models.URL, error) {
	url, ok := m.urls[id]
	if !ok || url.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	return url, nil
}

func (m *mockURLRepo) FindByShortCode(code string) (*models.URL, error) {
	for _, url := range m.urls {
		if url.ShortCode == code {
			return url, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockURLRepo) FindByUserID(userID uuid.UUID, search string, activeFilter *bool, page, limit int) ([]models.URL, int64, error) {
	var result []models.URL
	for _, u := range m.urls {
		if u.UserID != userID {
			continue
		}
		if activeFilter != nil && u.IsActive != *activeFilter {
			continue
		}
		result = append(result, *u)
	}
	return result, int64(len(result)), nil
}

func (m *mockURLRepo) Update(url *models.URL) error {
	m.urls[url.ID] = url
	return nil
}

func (m *mockURLRepo) Delete(id, userID uuid.UUID) error {
	delete(m.urls, id)
	return nil
}

func (m *mockURLRepo) IncrementClickCount(id uuid.UUID) error {
	if u, ok := m.urls[id]; ok {
		u.ClickCount++
	}
	return nil
}

func (m *mockURLRepo) ExistsByShortCode(code string) (bool, error) {
	for _, u := range m.urls {
		if u.ShortCode == code {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockURLRepo) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	for _, u := range m.urls {
		if u.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *mockURLRepo) TotalClicksByUserID(userID uuid.UUID) (int64, error) {
	var total int64
	for _, u := range m.urls {
		if u.UserID == userID {
			total += u.ClickCount
		}
	}
	return total, nil
}

func TestURLActiveInactiveFeature(t *testing.T) {
	mockRepo := newMockURLRepo()
	service := NewURLService(mockRepo, nil)

	userID := uuid.New()

	// 1. Create URL
	urlItem, err := service.CreateURL(userID, "https://google.com", "test-toggle", "Test Link", "", nil, nil)
	if err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}

	if !urlItem.IsActive {
		t.Fatalf("Expected new URL to be active by default, got %v", urlItem.IsActive)
	}

	// 2. Resolve active URL
	resolved, isPwRequired, err := service.ResolveURL("test-toggle")
	if err != nil {
		t.Fatalf("Expected resolve to succeed, got error: %v", err)
	}
	if isPwRequired {
		t.Fatalf("Expected password not required")
	}
	if resolved.OriginalURL != "https://google.com" {
		t.Fatalf("Expected original URL https://google.com, got %s", resolved.OriginalURL)
	}

	// 3. Toggle Status to Inactive
	toggled, err := service.ToggleStatus(urlItem.ID, userID)
	if err != nil {
		t.Fatalf("Failed to toggle status: %v", err)
	}
	if toggled.IsActive {
		t.Fatalf("Expected URL to become inactive, but is_active is true")
	}

	// 4. Resolve inactive URL should return ErrURLInactive
	_, _, err = service.ResolveURL("test-toggle")
	if !errors.Is(err, ErrURLInactive) {
		t.Fatalf("Expected ErrURLInactive, got %v", err)
	}

	// 5. Verify Password for inactive URL should return ErrURLInactive
	_, err = service.VerifyURLPassword("test-toggle", "")
	if !errors.Is(err, ErrURLInactive) {
		t.Fatalf("Expected ErrURLInactive on verify password, got %v", err)
	}

	// 6. Filter by active status
	activeTrue := true
	activeFalse := false

	activeList, totalActive, _ := service.GetURLs(userID, "", &activeTrue, 1, 10)
	if totalActive != 0 || len(activeList) != 0 {
		t.Fatalf("Expected 0 active URLs, got %d", totalActive)
	}

	inactiveList, totalInactive, _ := service.GetURLs(userID, "", &activeFalse, 1, 10)
	if totalInactive != 1 || len(inactiveList) != 1 {
		t.Fatalf("Expected 1 inactive URL, got %d", totalInactive)
	}

	// 7. Toggle back to Active
	toggledBack, err := service.ToggleStatus(urlItem.ID, userID)
	if err != nil {
		t.Fatalf("Failed to toggle back: %v", err)
	}
	if !toggledBack.IsActive {
		t.Fatalf("Expected URL to become active again, but is_active is false")
	}

	// 8. Resolve active again
	resolvedAgain, _, err := service.ResolveURL("test-toggle")
	if err != nil {
		t.Fatalf("Expected resolve to succeed after re-activating, got: %v", err)
	}
	if resolvedAgain.OriginalURL != "https://google.com" {
		t.Fatalf("Expected original URL https://google.com, got %s", resolvedAgain.OriginalURL)
	}
}
