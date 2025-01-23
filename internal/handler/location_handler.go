package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
