package routes

import (
	"github.com/OtchereDev/hrms-go/internal/api/handlers"
	"github.com/OtchereDev/hrms-go/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	permMiddleware *middleware.PermissionMiddleware,
) {
	// API version prefix
	api := app.Group("/api")

	// Method routes (Frappe-compatible API pattern)
	method := api.Group("/method")

	// Public authentication routes
	method.Post("/login", authHandler.Login)
	method.Post("/logout", authMiddleware.Authenticate, authHandler.Logout)
	method.Post("/refresh_token", authHandler.RefreshToken)

	// Protected routes - require authentication
	authenticated := method.Use(authMiddleware.Authenticate)

	// User info routes
	authenticated.Get("/hrms.api.get_current_user_info", authHandler.GetCurrentUserInfo)
	authenticated.Get("/hrms.api.get_current_employee_info", authHandler.GetCurrentEmployeeInfo)

	// TODO: Add more routes as modules are implemented
	// - Employee routes
	// - Attendance routes
	// - Leave routes
	// - Payroll routes
	// - etc.

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "hrms-api",
		})
	})

	// API info route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "HRMS API",
			"version": "1.0.0",
			"status":  "running",
		})
	})
}
