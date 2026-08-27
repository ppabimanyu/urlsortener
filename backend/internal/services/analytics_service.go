package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"backend/internal/models"
	"backend/internal/repositories"

	"github.com/google/uuid"
)

type OverviewAnalytics struct {
	TotalLinks   int64                        `json:"total_links"`
	TotalClicks  int64                        `json:"total_clicks"`
	ActiveLinks  int64                        `json:"active_links"`
	Timeline     []repositories.DateClickStat `json:"timeline"`
	RecentClicks []models.URLClick            `json:"recent_clicks"`
}

type URLAnalyticsResponse struct {
	URL       *models.URL                  `json:"url"`
	Timeline  []repositories.DateClickStat `json:"timeline"`
	Devices   []repositories.GroupCountStat `json:"devices"`
	OS        []repositories.GroupCountStat `json:"os"`
	Browsers  []repositories.GroupCountStat `json:"browsers"`
	Referrers []repositories.GroupCountStat `json:"referrers"`
	Countries []repositories.GroupCountStat `json:"countries"`
	Cities    []repositories.GroupCountStat `json:"cities"`
}

type AnalyticsService interface {
	GetOverview(userID uuid.UUID) (*OverviewAnalytics, error)
	GetURLAnalytics(urlID, userID uuid.UUID, days int) (*URLAnalyticsResponse, error)
	ExportURLAnalyticsCSV(urlID, userID uuid.UUID) ([]byte, string, error)
}

type analyticsService struct {
	urlRepo   repositories.URLRepository
	clickRepo repositories.ClickRepository
}

func NewAnalyticsService(urlRepo repositories.URLRepository, clickRepo repositories.ClickRepository) AnalyticsService {
	return &analyticsService{
		urlRepo:   urlRepo,
		clickRepo: clickRepo,
	}
}

func (s *analyticsService) GetOverview(userID uuid.UUID) (*OverviewAnalytics, error) {
	totalLinks, err := s.urlRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	totalClicks, err := s.urlRepo.TotalClicksByUserID(userID)
	if err != nil {
		return nil, err
	}

	activeFilter := true
	_, activeLinks, err := s.urlRepo.FindByUserID(userID, "", &activeFilter, 1, 1)
	if err != nil {
		return nil, err
	}

	timeline, err := s.clickRepo.GetOverviewTimeline(userID, 14)
	if err != nil {
		return nil, err
	}

	recentClicks, err := s.clickRepo.GetRecentClicksByUserID(userID, 10)
	if err != nil {
		return nil, err
	}

	return &OverviewAnalytics{
		TotalLinks:   totalLinks,
		TotalClicks:  totalClicks,
		ActiveLinks:  activeLinks,
		Timeline:     timeline,
		RecentClicks: recentClicks,
	}, nil
}

func (s *analyticsService) GetURLAnalytics(urlID, userID uuid.UUID, days int) (*URLAnalyticsResponse, error) {
	urlItem, err := s.urlRepo.FindByIDAndUserID(urlID, userID)
	if err != nil {
		return nil, err
	}

	timeline, err := s.clickRepo.GetClicksByDate(urlID, days)
	if err != nil {
		return nil, err
	}

	devices, err := s.clickRepo.GetDeviceStats(urlID)
	if err != nil {
		return nil, err
	}

	osStats, err := s.clickRepo.GetOSStats(urlID)
	if err != nil {
		return nil, err
	}

	browsers, err := s.clickRepo.GetBrowserStats(urlID)
	if err != nil {
		return nil, err
	}

	referrers, err := s.clickRepo.GetReferrerStats(urlID)
	if err != nil {
		return nil, err
	}

	countries, err := s.clickRepo.GetCountryStats(urlID)
	if err != nil {
		return nil, err
	}

	cities, err := s.clickRepo.GetCityStats(urlID)
	if err != nil {
		return nil, err
	}

	return &URLAnalyticsResponse{
		URL:       urlItem,
		Timeline:  timeline,
		Devices:   devices,
		OS:        osStats,
		Browsers:  browsers,
		Referrers: referrers,
		Countries: countries,
		Cities:    cities,
	}, nil
}

func (s *analyticsService) ExportURLAnalyticsCSV(urlID, userID uuid.UUID) ([]byte, string, error) {
	urlItem, err := s.urlRepo.FindByIDAndUserID(urlID, userID)
	if err != nil {
		return nil, "", err
	}

	clicks, err := s.clickRepo.GetAllClicksForExport(urlID)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV Header
	_ = writer.Write([]string{"Click ID", "Clicked At (UTC)", "IP Address", "Device Type", "OS", "Browser", "Referrer", "Country", "Country Code", "City"})

	for _, c := range clicks {
		row := []string{
			fmt.Sprintf("%d", c.ID),
			c.ClickedAt.Format(time.RFC3339),
			c.IPAddress,
			c.DeviceType,
			c.OS,
			c.Browser,
			c.Referrer,
			c.Country,
			c.CountryCode,
			c.City,
		}
		_ = writer.Write(row)
	}

	writer.Flush()
	filename := fmt.Sprintf("analytics_%s_%s.csv", urlItem.ShortCode, time.Now().Format("20060102_150405"))
	return buf.Bytes(), filename, nil
}
