package repository

import (
	"errors"
	"sunrise_project/internal/dao"

	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) GetByIP(ip string) (*dao.Location, error) {
	var location dao.Location
	result := r.db.Where("ip = ?", ip).First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &location, nil
}

func (r *LocationRepository) Create(location *dao.Location) error {
	return r.db.Create(location).Error
}

func (r *LocationRepository) GetAll() ([]dao.Location, error) {
	var locations []dao.Location
	result := r.db.Find(&locations)
	if result.Error != nil {
		return nil, result.Error
	}
	return locations, nil
}
