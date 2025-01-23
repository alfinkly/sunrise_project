package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sunrise_project/internal/dao"
	"sunrise_project/internal/platform"
	"sunrise_project/internal/repository"
)

type LocationService struct {
	repo       *repository.LocationRepository
	ipApiURL   string
	ipApiToken string
}

func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{
		repo:       repo,
		ipApiURL:   "https://ipinfo.io",
		ipApiToken: platform.GetIpInfoToken(),
	}
}

func (s *LocationService) GetLocationByIP(ip string) (*dao.Location, error) {
	location, err := s.repo.GetByIP(ip)
	if err != nil {
		return nil, fmt.Errorf("get location by IP from DB: %w", err)
	}

	if location != nil {
		log.Println("location found in DB")
		return location, nil
	}

	log.Println("location not found in DB, making request to external API")
	apiURL := fmt.Sprintf("%s/%s?token=%s", s.ipApiURL, ip, s.ipApiToken)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("get location by IP from API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from IP API: %d", resp.StatusCode)
	}
	fmt.Println(resp.Body)
	var ipInfoResponse struct {
		IP      string `json:"ip"`
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
	}
	err = json.NewDecoder(resp.Body).Decode(&ipInfoResponse)
	if err != nil {
		return nil, fmt.Errorf("decode response from IP API: %w", err)
	}

	location = &dao.Location{
		IP:      ipInfoResponse.IP,
		Country: ipInfoResponse.Country,
		City:    ipInfoResponse.City,
	}

	err = s.repo.Create(location)
	if err != nil {
		log.Printf("Failed to save location to DB: %v", err)
	}

	return location, nil
}

func (s *LocationService) GetAllLocations() ([]dao.Location, error) {
	return s.repo.GetAll()
}

func (s *LocationService) CreateLocation(location *dao.Location) error {
	return s.repo.Create(location)
}

func (s *LocationService) UpdateLocation(location *dao.Location) error {
	return s.repo.Update(location)
}

func (s *LocationService) DeleteLocation(ip string) error {
	return s.repo.Delete(ip)
}
