// File: /quicklynks/backend/internal/controllers/link_controller.go
package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thisisjackii/quicklynks/backend/internal/models"
	"github.com/thisisjackii/quicklynks/backend/internal/utils"
	"gorm.io/gorm"
)

type LinkController struct {
	DB *gorm.DB
}

type CreateLinkInput struct {
	URL string `json:"url" binding:"required,url"`
}

type LinkResponse struct {
	ID          uint      `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int64     `json:"click_count"`
}

// CreateLink godoc
// @Summary Create a new short link
// @Description Creates a new short link for the authenticated user.
// @Tags links
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   input  body   CreateLinkInput  true  "URL to shorten"
// @Success 201 {object} models.Link "Successfully created link"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /links [post]
func (lc *LinkController) CreateLink(c *gin.Context) {
	var input CreateLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// Generate a unique short code
	var shortCode string
	for {
		shortCode = utils.GenerateShortCode()
		var existingLink models.Link
		if err := lc.DB.Where("short_code = ?", shortCode).First(&existingLink).Error; err == gorm.ErrRecordNotFound {
			break // Unique code found
		}
	}

	link := models.Link{
		OriginalURL: input.URL,
		ShortCode:   shortCode,
		UserID:      userID.(uint),
	}

	if err := lc.DB.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}

	c.JSON(http.StatusCreated, link)
}

// GetUserLinks godoc
// @Summary Get all links for the authenticated user
// @Description Retrieves a list of all links created by the current user.
// @Tags links
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} LinkResponse "List of user's links with click counts"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /links [get]
func (lc *LinkController) GetUserLinks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var links []models.Link
	if err := lc.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}

	// For each link, get the click count
	var response []LinkResponse
	for _, link := range links {
		var clickCount int64
		lc.DB.Model(&models.Click{}).Where("link_id = ?", link.ID).Count(&clickCount)
		response = append(response, LinkResponse{
			ID:          link.ID,
			OriginalURL: link.OriginalURL,
			ShortCode:   link.ShortCode,
			UserID:      link.UserID,
			CreatedAt:   link.CreatedAt,
			ClickCount:  clickCount,
		})
	}

	c.JSON(http.StatusOK, response)
}

// RedirectLink godoc
// @Summary Redirect to original URL
// @Description Redirects a short code to its corresponding original URL and tracks the click.
// @Tags redirect
// @Produce  html
// @Param   shortCode  path   string  true  "Short code of the link"
// @Success 302 "Redirects to the original URL"
// @Failure 404 {object} map[string]interface{} "Link not found"
// @Router /{shortCode} [get]
func (lc *LinkController) RedirectLink(c *gin.Context) {
	shortCode := c.Param("shortCode")

	var link models.Link
	if err := lc.DB.Where("short_code = ?", shortCode).First(&link).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Track the click asynchronously to not slow down the redirect
	go func() {
		click := models.Click{
			LinkID:    link.ID,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}
		lc.DB.Create(&click)
	}()

	c.Redirect(http.StatusFound, link.OriginalURL)
}
