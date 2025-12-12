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
	employeeHandler *handlers.EmployeeHandler,
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
	authenticated.Get("/hrms.api.get_current_employee_info", employeeHandler.GetCurrentEmployeeInfo)

	// Employee routes
	employeeRoutes := authenticated.Group("")

	// List employees (accessible to HR User and above)
	employeeRoutes.Get("/hrms.api.get_employees",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.ListEmployees,
	)

	employeeRoutes.Get("/hrms.api.get_active_employees",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.GetActiveEmployees,
	)

	employeeRoutes.Get("/hrms.api.get_employees_by_department",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.GetEmployeesByDepartment,
	)

	// Get employee (accessible to all authenticated users)
	employeeRoutes.Get("/hrms.hr.doctype.employee.employee.get_employee",
		employeeHandler.GetEmployee,
	)

	employeeRoutes.Get("/hrms.hr.doctype.employee.employee.get_employee_details",
		employeeHandler.GetEmployeeDetails,
	)

	// Create/Update/Delete employee (HR Manager and above only)
	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.create_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.CreateEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.update_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.UpdateEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.delete_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.DeleteEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.update_employee_status",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.UpdateEmployeeStatus,
	)

	// TODO: Add more routes as modules are implemented
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
