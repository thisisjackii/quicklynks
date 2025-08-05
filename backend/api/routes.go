# File: backend/api/routes.go
package api

import "github.com/labstack/echo/v4"

// RegisterRoutes sets up all the API routes.
func RegisterRoutes(e *echo.Echo, h *Handler) {
	v1 := e.Group("/api/v1")

	// Public routes
	v1.POST("/users/register", h.RegisterUser)
	v1.POST("/users/login", h.LoginUser)
	v1.GET("/profiles/:username", h.GetProfile)

	// Authenticated routes
	authGroup := v1.Group("/me")
	authGroup.Use(h.AuthMiddleware)
	authGroup.GET("/links", h.GetMyLinks)
	authGroup.POST("/links", h.CreateLink)
	authGroup.DELETE("/links/:id", h.DeleteLink)
}