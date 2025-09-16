package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserHandler handles user-related requests
type UserHandler struct {
	users map[int]*User
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		users: make(map[int]*User),
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data")
	}

	if user.Name == "" || user.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name and email are required")
	}

	user.ID = len(h.users) + 1
	h.users[user.ID] = user

	return c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	// Simple validation for demo
	if id == "1" {
		user := &User{ID: 1, Name: "John Doe", Email: "john@example.com"}
		return c.JSON(http.StatusOK, user)
	}

	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}

// GetUsers returns all users
func (h *UserHandler) GetUsers(c echo.Context) error {
	users := make([]*User, 0, len(h.users))
	for _, user := range h.users {
		users = append(users, user)
	}
	return c.JSON(http.StatusOK, users)
}

// Health check endpoint
func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"message": "Service is healthy",
	})
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize handler
	userHandler := NewUserHandler()

	// Routes
	e.GET("/health", health)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users", userHandler.GetUsers)
	e.GET("/users/:id", userHandler.GetUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}