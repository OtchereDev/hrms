package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OtchereDev/hrms-go/internal/api/handlers"
	frappeHandlers "github.com/OtchereDev/hrms-go/internal/api/handlers/frappe"
	"github.com/OtchereDev/hrms-go/internal/api/middleware"
	"github.com/OtchereDev/hrms-go/internal/api/routes"
	"github.com/OtchereDev/hrms-go/internal/config"
	"github.com/OtchereDev/hrms-go/internal/core"
	"github.com/OtchereDev/hrms-go/internal/core/frappe"
	"github.com/OtchereDev/hrms-go/internal/core/services/auth"
	"github.com/OtchereDev/hrms-go/internal/core/services/holiday"
	"github.com/OtchereDev/hrms-go/internal/core/services/shift"
	"github.com/OtchereDev/hrms-go/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	logger.Info("Starting HRMS API Server...")
	logger.Info(fmt.Sprintf("Environment: %s", cfg.App.Environment))

	// Initialize database
	db, err := core.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	logger.Info("Database connected successfully")

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// Seed default data (only in development)
	if cfg.App.Environment == "development" {
		if err := db.SeedDefaultData(); err != nil {
			logger.Warning("Failed to seed default data:", err)
		}
	}

	// Initialize services
	authService := auth.NewAuthService(db.DB, cfg)
	jwtService := auth.NewJWTService(&cfg.JWT)
	holidayService := holiday.NewHolidayService(db.DB)
	shiftService := shift.NewShiftService(db.DB)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, authService)
	permMiddleware := middleware.NewPermissionMiddleware(db.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	employeeHandler := handlers.NewEmployeeHandler(db.DB)
	attendanceHandler := handlers.NewAttendanceHandler(db.DB)
	leaveHandler := handlers.NewLeaveHandler(db.DB)
	payrollHandler := handlers.NewPayrollHandler(db.DB)
	performanceHandler := handlers.NewPerformanceHandler(db.DB)
	holidayHandler := handlers.NewHolidayHandler(holidayService)
	shiftHandler := handlers.NewShiftHandler(shiftService)

	// Initialize Frappe compatibility layer
	doctypeService := frappe.NewDocTypeService(db.DB)
	resourceHandler := frappeHandlers.NewResourceHandler(doctypeService)
	methodHandler := frappeHandlers.NewMethodHandler(
		db.DB,
		attendanceHandler,
		leaveHandler,
		payrollHandler,
		performanceHandler,
		holidayHandler,
		employeeHandler,
		shiftHandler,
	)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(cfg.Server.AllowedOrigins),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Compression middleware
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Rate limiting middleware
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.Server.RateLimit,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	}))

	// Setup routes
	routes.SetupRoutes(app, authHandler, employeeHandler, attendanceHandler, leaveHandler, payrollHandler, performanceHandler, resourceHandler, methodHandler, authMiddleware, permMiddleware)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Server starting on %s", addr))

	// Graceful shutdown
	go func() {
		if err := app.Listen(addr); err != nil {
			logger.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exited successfully")
}

// customErrorHandler handles errors in a standardized way
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
		"error":   err.Error(),
	})
}

// corsOrigins formats CORS origins
func corsOrigins(origins []string) string {
	if len(origins) == 0 {
		return "*"
	}
	result := ""
	for i, origin := range origins {
		if i > 0 {
			result += ","
		}
		result += origin
	}
	return result
}
