# File: backend/api/handlers.go
package api

import (
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/thisisjackii/quicklynks/backend/internal/models"
	"github.com/thisisjackii/quicklynks/backend/internal/store"
	"github.com/thisisjackii/quicklynks/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Handler holds all dependencies for API handlers.
type Handler struct {
	userStore *store.UserStore
	linkStore *store.LinkStore
	jwtSecret string
}

// NewHandler creates a new handler with its dependencies.
func NewHandler(us *store.UserStore, ls *store.LinkStore, secret string) *Handler {
	return &Handler{
		userStore: us,
		linkStore: ls,
		jwtSecret: secret,
	}
}

// --- User Handlers ---

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser handles new user registration.
func (h *Handler) RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Basic validation
	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters long"})
	}
	if req.Username == "" || req.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and email are required"})
	}

	// Check if username or email already exists
	if existing, _ := h.userStore.GetByUsername(req.Username); existing != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already taken"})
	}
	if existing, _ := h.userStore.GetByEmail(req.Email); existing != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Email already registered"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not process password"})
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.userStore.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user"})
	}

	user.PasswordHash = "" // Clear hash before sending response
	return c.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser handles user authentication.
func (h *Handler) LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	user, err := h.userStore.GetByEmail(req.Email)
	if err != nil || user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Create JWT token
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HttpOnly = true
	cookie.Path = "/"
	// Set SameSite=Strict for better security
	cookie.SameSite = http.SameSiteStrictMode 
	// In production, set Secure=true
	// cookie.Secure = true
	c.SetCookie(cookie)

	user.PasswordHash = ""
	return c.JSON(http.StatusOK, user)
}

// --- Link Handlers ---

type CreateLinkRequest struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// CreateLink adds a new link for the authenticated user.
func (h *Handler) CreateLink(c echo.Context) error {
	var req CreateLinkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}
	if req.Title == "" || req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title and URL are required"})
	}
	
	userIDStr := c.Get("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	link := &models.Link{
		UserID: userID,
		Title:  req.Title,
		URL:    req.URL,
	}

	if err := h.linkStore.Create(link); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create link"})
	}

	return c.JSON(http.StatusCreated, link)
}

// GetMyLinks retrieves all links for the authenticated user.
func (h *Handler) GetMyLinks(c echo.Context) error {
	userIDStr := c.Get("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	links, err := h.linkStore.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch links"})
	}

	return c.JSON(http.StatusOK, links)
}

// DeleteLink removes a link for the authenticated user.
func (h *Handler) DeleteLink(c echo.Context) error {
	linkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid link ID"})
	}

	userIDStr := c.Get("userID").(string)
	userID, _ := strconv.Atoi(userIDStr)

	err = h.linkStore.Delete(linkID, userID)
	if err != nil {
		// This could be because the link doesn't exist, or it belongs to another user.
		// For security, we return a generic 404.
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// --- Public Profile Handler ---

type ProfileResponse struct {
	Username string         `json:"username"`
	Links    []models.Link `json:"links"`
}

// GetProfile retrieves a user's public profile data.
func (h *Handler) GetProfile(c echo.Context) error {
	username := c.Param("username")
	
	user, err := h.userStore.GetByUsername(username)
	if err != nil || user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	links, err := h.linkStore.GetByUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch profile links"})
	}
	
	// Ensure links are not nil for JSON marshalling
	if links == nil {
		links = []models.Link{}
	}

	response := ProfileResponse{
		Username: user.Username,
		Links:    links,
	}

	return c.JSON(http.StatusOK, response)
}