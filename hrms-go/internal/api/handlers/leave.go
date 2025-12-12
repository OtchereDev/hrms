package handlers

import (
	"strconv"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"github.com/OtchereDev/hrms-go/internal/core/services/leave"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// LeaveHandler handles leave-related requests
type LeaveHandler struct {
	leaveService *leave.LeaveService
}

// NewLeaveHandler creates a new leave handler
func NewLeaveHandler(db *gorm.DB) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leave.NewLeaveService(db),
	}
}

// ApplyLeave creates a new leave application
// POST /api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave
func (h *LeaveHandler) ApplyLeave(c *fiber.Ctx) error {
	var req leave.ApplyLeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	leaveApp, err := h.leaveService.ApplyLeave(c.Context(), &req)
	if err != nil {
		switch err {
		case leave.ErrInsufficientLeaveBalance:
			return response.BadRequest(c, "insufficient leave balance")
		case leave.ErrOverlappingLeaves:
			return response.Conflict(c, "overlapping leave applications exist")
		case leave.ErrEmployeeRequired, leave.ErrLeaveTypeRequired, leave.ErrInvalidDates:
			return response.BadRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to apply for leave")
		}
	}

	return response.Created(c, leaveApp, "Leave application created successfully")
}

// GetLeaveApplication retrieves a leave application
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.get_leave_application
func (h *LeaveHandler) GetLeaveApplication(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	leaveApp, err := h.leaveService.GetLeaveApplicationByID(c.Context(), uint(id))
	if err != nil {
		if err == leave.ErrLeaveApplicationNotFound {
			return response.NotFound(c, "leave application not found")
		}
		return response.InternalServerError(c, "failed to get leave application")
	}

	return response.Success(c, leaveApp, "success")
}

// ListLeaveApplications retrieves leave applications with filters
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.list_leave_applications
func (h *LeaveHandler) ListLeaveApplications(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.LeaveApplicationFilters{
		Employee:      c.Query("employee"),
		Company:       c.Query("company"),
		LeaveType:     c.Query("leave_type"),
		Status:        c.Query("status"),
		LeaveApprover: c.Query("leave_approver"),
	}

	// Parse dates if provided
	if fromDateStr := c.Query("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filters.FromDate = fromDate
		}
	}
	if toDateStr := c.Query("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filters.ToDate = toDate
		}
	}

	leaves, total, err := h.leaveService.ListLeaveApplications(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list leave applications")
	}

	return response.Paginated(c, leaves, page, pageSize, total)
}

// ApproveLeaveApplication approves a leave application
// POST /api/method/hrms.hr.doctype.leave_application.leave_application.approve_leave
func (h *LeaveHandler) ApproveLeaveApplication(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	leaveApp, err := h.leaveService.ApproveLeaveApplication(c.Context(), uint(id), approver)
	if err != nil {
		switch err {
		case leave.ErrLeaveApplicationNotFound:
			return response.NotFound(c, "leave application not found")
		case leave.ErrInsufficientLeaveBalance:
			return response.BadRequest(c, "insufficient leave balance")
		default:
			return response.InternalServerError(c, "failed to approve leave application")
		}
	}

	return response.Success(c, leaveApp, "Leave application approved successfully")
}

// RejectLeaveApplication rejects a leave application
// POST /api/method/hrms.hr.doctype.leave_application.leave_application.reject_leave
func (h *LeaveHandler) RejectLeaveApplication(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")
	reason := c.Query("reason")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	leaveApp, err := h.leaveService.RejectLeaveApplication(c.Context(), uint(id), approver, reason)
	if err != nil {
		if err == leave.ErrLeaveApplicationNotFound {
			return response.NotFound(c, "leave application not found")
		}
		return response.InternalServerError(c, "failed to reject leave application")
	}

	return response.Success(c, leaveApp, "Leave application rejected successfully")
}

// CancelLeaveApplication cancels a leave application
// POST /api/method/hrms.hr.doctype.leave_application.leave_application.cancel_leave
func (h *LeaveHandler) CancelLeaveApplication(c *fiber.Ctx) error {
	idStr := c.Query("id")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	leaveApp, err := h.leaveService.CancelLeaveApplication(c.Context(), uint(id))
	if err != nil {
		if err == leave.ErrLeaveApplicationNotFound {
			return response.NotFound(c, "leave application not found")
		}
		return response.InternalServerError(c, "failed to cancel leave application")
	}

	return response.Success(c, leaveApp, "Leave application cancelled successfully")
}

// AllocateLeave allocates leaves to an employee
// POST /api/method/hrms.hr.doctype.leave_allocation.leave_allocation.allocate_leave
func (h *LeaveHandler) AllocateLeave(c *fiber.Ctx) error {
	var req leave.AllocateLeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	allocation, err := h.leaveService.AllocateLeave(c.Context(), &req)
	if err != nil {
		switch err {
		case leave.ErrEmployeeRequired, leave.ErrLeaveTypeRequired, leave.ErrInvalidDates:
			return response.BadRequest(c, err.Error())
		case leave.ErrLeaveTypeNotFound:
			return response.NotFound(c, "leave type not found")
		default:
			return response.InternalServerError(c, "failed to allocate leave")
		}
	}

	return response.Created(c, allocation, "Leave allocated successfully")
}

// GetLeaveBalance retrieves leave balance for an employee
// GET /api/method/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance
func (h *LeaveHandler) GetLeaveBalance(c *fiber.Ctx) error {
	employee := c.Query("employee")
	leaveType := c.Query("leave_type")

	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	if leaveType != "" {
		// Get balance for specific leave type
		balance, err := h.leaveService.GetLeaveBalance(c.Context(), employee, leaveType)
		if err != nil {
			return response.InternalServerError(c, "failed to get leave balance")
		}

		return response.Success(c, fiber.Map{
			"employee":   employee,
			"leave_type": leaveType,
			"balance":    balance,
		}, "success")
	}

	// Get all leave balances
	balances, err := h.leaveService.GetAllLeaveBalances(c.Context(), employee)
	if err != nil {
		return response.InternalServerError(c, "failed to get leave balances")
	}

	return response.Success(c, fiber.Map{
		"employee": employee,
		"balances": balances,
	}, "success")
}

// GetActiveLeaveTypes retrieves all active leave types
// GET /api/method/hrms.hr.doctype.leave_type.leave_type.get_active_leave_types
func (h *LeaveHandler) GetActiveLeaveTypes(c *fiber.Ctx) error {
	leaveTypes, err := h.leaveService.GetActiveLeaveTypes(c.Context())
	if err != nil {
		return response.InternalServerError(c, "failed to get leave types")
	}

	return response.Success(c, leaveTypes, "success")
}

// CreateLeaveEncashment creates a leave encashment request
// POST /api/method/hrms.hr.doctype.leave_encashment.leave_encashment.create_encashment
func (h *LeaveHandler) CreateLeaveEncashment(c *fiber.Ctx) error {
	employee := c.Query("employee")
	leaveType := c.Query("leave_type")
	encashableDaysStr := c.Query("encashable_days")

	if employee == "" || leaveType == "" || encashableDaysStr == "" {
		return response.BadRequest(c, "employee, leave_type, and encashable_days are required")
	}

	encashableDays, err := strconv.ParseFloat(encashableDaysStr, 64)
	if err != nil {
		return response.BadRequest(c, "invalid encashable_days")
	}

	encashment, err := h.leaveService.CreateLeaveEncashment(c.Context(), employee, leaveType, encashableDays)
	if err != nil {
		switch err {
		case leave.ErrInsufficientLeaveBalance:
			return response.BadRequest(c, "insufficient leave balance")
		case leave.ErrEmployeeRequired, leave.ErrLeaveTypeRequired:
			return response.BadRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to create leave encashment")
		}
	}

	return response.Created(c, encashment, "Leave encashment created successfully")
}

// GetLeaveDetails returns leave details for an employee on a specific date
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.get_leave_details
func (h *LeaveHandler) GetLeaveDetails(c *fiber.Ctx) error {
	employee := c.Query("employee")
	dateStr := c.Query("date")
	forSalarySlip := c.Query("for_salary_slip") == "1"

	if employee == "" || dateStr == "" {
		return response.BadRequest(c, "employee and date are required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return response.BadRequest(c, "invalid date format, expected YYYY-MM-DD")
	}

	leaveDetails, err := h.leaveService.GetLeaveDetails(c.Context(), employee, date, forSalarySlip)
	if err != nil {
		return response.InternalServerError(c, "failed to get leave details")
	}

	return response.Success(c, leaveDetails, "Leave details retrieved successfully")
}

// GetNumberOfLeaveDays calculates the number of leave days
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.get_number_of_leave_days
func (h *LeaveHandler) GetNumberOfLeaveDays(c *fiber.Ctx) error {
	employee := c.Query("employee")
	leaveType := c.Query("leave_type")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	halfDay := c.Query("half_day") == "1"
	halfDayDateStr := c.Query("half_day_date")

	if employee == "" || leaveType == "" || fromDateStr == "" || toDateStr == "" {
		return response.BadRequest(c, "employee, leave_type, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	var halfDayDate *time.Time
	if halfDayDateStr != "" {
		hdd, err := time.Parse("2006-01-02", halfDayDateStr)
		if err != nil {
			return response.BadRequest(c, "invalid half_day_date format, expected YYYY-MM-DD")
		}
		halfDayDate = &hdd
	}

	days, err := h.leaveService.GetNumberOfLeaveDays(c.Context(), employee, leaveType, fromDate, toDate, halfDay, halfDayDate)
	if err != nil {
		return response.InternalServerError(c, "failed to calculate leave days")
	}

	return response.Success(c, map[string]interface{}{"leave_days": days}, "Leave days calculated successfully")
}

// GetLeaveBalanceOn returns the leave balance on a specific date
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.get_leave_balance_on
func (h *LeaveHandler) GetLeaveBalanceOn(c *fiber.Ctx) error {
	employee := c.Query("employee")
	leaveType := c.Query("leave_type")
	dateStr := c.Query("date")

	if employee == "" || leaveType == "" || dateStr == "" {
		return response.BadRequest(c, "employee, leave_type, and date are required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return response.BadRequest(c, "invalid date format, expected YYYY-MM-DD")
	}

	balance, err := h.leaveService.GetLeaveBalanceOn(c.Context(), employee, leaveType, date)
	if err != nil {
		return response.InternalServerError(c, "failed to get leave balance")
	}

	return response.Success(c, map[string]interface{}{"leave_balance": balance}, "Leave balance retrieved successfully")
}

// GetLeavesForPeriod returns all approved leaves for an employee in a date range
// GET /api/method/hrms.hr.doctype.leave_application.leave_application.get_leaves_for_period
func (h *LeaveHandler) GetLeavesForPeriod(c *fiber.Ctx) error {
	employee := c.Query("employee")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if employee == "" || fromDateStr == "" || toDateStr == "" {
		return response.BadRequest(c, "employee, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	leaves, err := h.leaveService.GetLeavesForPeriod(c.Context(), employee, fromDate, toDate)
	if err != nil {
		return response.InternalServerError(c, "failed to get leaves for period")
	}

	return response.Success(c, leaves, "Leaves for period retrieved successfully")
}
