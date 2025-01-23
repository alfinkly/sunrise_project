package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sunrise_project/internal/dao"
	"sunrise_project/internal/service"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (h *LocationHandler) GetLocationByIP(c *gin.Context) {
	ip := c.GetHeader("X-Forwarded-For")

	location, err := h.service.GetLocationByIP(ip)
	if err != nil {
		log.Printf("Failed to get location for IP %s: %v", ip, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ip":      location.IP,
		"country": location.Country,
		"city":    location.City,
	})
}

func (h *LocationHandler) GetLocationByCustomIP(c *gin.Context) {
	customIP := c.Param("ip")

	location, err := h.service.GetLocationByIP(customIP)
	if err != nil {
		log.Printf("Failed to get location for IP %s: %v", customIP, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ip":      location.IP,
		"country": location.Country,
		"city":    location.City,
	})
}

func (h *LocationHandler) GetAllLocations(c *gin.Context) {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		log.Printf("Failed to get all locations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get locations"})
		return
	}

	c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var newLocation dao.Location
	if err := c.BindJSON(&newLocation); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.service.CreateLocation(&newLocation)
	if err != nil {
		log.Printf("Failed to create location: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create location"})
		return
	}

	c.JSON(http.StatusCreated, newLocation)
}

func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	ip := c.Param("ip")

	existingLocation, err := h.service.GetLocationByIP(ip)
	if err != nil {
		log.Printf("Failed to get location for IP %s: %v", ip, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get location"})
		return
	}

	if existingLocation == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	var updatedLocation dao.Location
	if err := c.BindJSON(&updatedLocation); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedLocation.IP = existingLocation.IP

	err = h.service.UpdateLocation(&updatedLocation)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		} else {
			log.Printf("Failed to update location: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedLocation)
}

func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	ip := c.Param("ip")

	err := h.service.DeleteLocation(ip)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		} else {
			log.Printf("Failed to delete location: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location deleted"})
}
