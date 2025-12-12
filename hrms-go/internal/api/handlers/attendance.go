package handlers

import (
	"strconv"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"github.com/OtchereDev/hrms-go/internal/core/services/attendance"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AttendanceHandler handles attendance-related requests
type AttendanceHandler struct {
	attendanceService *attendance.AttendanceService
}

// NewAttendanceHandler creates a new attendance handler
func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendance.NewAttendanceService(db),
	}
}

// MarkAttendance marks attendance for an employee
// POST /api/method/hrms.hr.doctype.attendance.attendance.mark_attendance
func (h *AttendanceHandler) MarkAttendance(c *fiber.Ctx) error {
	var req attendance.MarkAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	att, err := h.attendanceService.MarkAttendance(c.Context(), &req)
	if err != nil {
		switch err {
		case attendance.ErrAttendanceExists:
			return response.Conflict(c, "attendance already exists for this date")
		case attendance.ErrEmployeeRequired, attendance.ErrAttendanceDateRequired, attendance.ErrInvalidStatus:
			return response.BadRequest(c, err.Error())
		default:
			return response.InternalServerError(c, "failed to mark attendance")
		}
	}

	return response.Created(c, att, "Attendance marked successfully")
}

// GetAttendance retrieves an attendance record
// GET /api/method/hrms.hr.doctype.attendance.attendance.get_attendance
func (h *AttendanceHandler) GetAttendance(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	att, err := h.attendanceService.GetAttendanceByID(c.Context(), uint(id))
	if err != nil {
		if err == attendance.ErrAttendanceNotFound {
			return response.NotFound(c, "attendance record not found")
		}
		return response.InternalServerError(c, "failed to get attendance")
	}

	return response.Success(c, att, "success")
}

// ListAttendance retrieves attendance records with filters
// GET /api/method/hrms.hr.doctype.attendance.attendance.list_attendance
func (h *AttendanceHandler) ListAttendance(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.AttendanceFilters{
		Employee:   c.Query("employee"),
		Company:    c.Query("company"),
		Department: c.Query("department"),
		Status:     c.Query("status"),
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

	attendances, total, err := h.attendanceService.ListAttendance(c.Context(), filters, page, pageSize)
	if err != nil {
		return response.InternalServerError(c, "failed to list attendance")
	}

	return response.Paginated(c, attendances, page, pageSize, total)
}

// GetMonthlyAttendance retrieves monthly attendance for an employee
// GET /api/method/hrms.hr.doctype.attendance.attendance.get_monthly_attendance
func (h *AttendanceHandler) GetMonthlyAttendance(c *fiber.Ctx) error {
	employee := c.Query("employee")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if employee == "" || yearStr == "" || monthStr == "" {
		return response.BadRequest(c, "employee, year, and month are required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return response.BadRequest(c, "invalid year")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return response.BadRequest(c, "invalid month")
	}

	attendances, err := h.attendanceService.GetMonthlyAttendance(c.Context(), employee, year, month)
	if err != nil {
		return response.InternalServerError(c, "failed to get monthly attendance")
	}

	return response.Success(c, attendances, "success")
}

// Checkin handles employee check-in/check-out
// POST /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin
func (h *AttendanceHandler) Checkin(c *fiber.Ctx) error {
	var req attendance.CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	checkin, err := h.attendanceService.Checkin(c.Context(), &req)
	if err != nil {
		switch err {
		case attendance.ErrEmployeeRequired:
			return response.BadRequest(c, "employee is required")
		case attendance.ErrNotCheckedIn:
			return response.BadRequest(c, "employee not checked in")
		case attendance.ErrAlreadyCheckedIn:
			return response.Conflict(c, "employee already checked in")
		default:
			return response.InternalServerError(c, "failed to process checkin")
		}
	}

	return response.Created(c, checkin, "Checkin recorded successfully")
}

// GetTodayCheckins retrieves today's check-ins for an employee
// GET /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins
func (h *AttendanceHandler) GetTodayCheckins(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	checkins, err := h.attendanceService.GetTodayCheckinsForEmployee(c.Context(), employee)
	if err != nil {
		return response.InternalServerError(c, "failed to get checkins")
	}

	return response.Success(c, checkins, "success")
}

// CreateAttendanceRequest creates an attendance correction request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.create_request
func (h *AttendanceHandler) CreateAttendanceRequest(c *fiber.Ctx) error {
	var req attendance.AttendanceRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	attendanceReq, err := h.attendanceService.CreateAttendanceRequest(c.Context(), &req)
	if err != nil {
		if err == attendance.ErrEmployeeRequired {
			return response.BadRequest(c, err.Error())
		}
		return response.InternalServerError(c, "failed to create attendance request")
	}

	return response.Created(c, attendanceReq, "Attendance request created successfully")
}

// ApproveAttendanceRequest approves an attendance request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.approve_request
func (h *AttendanceHandler) ApproveAttendanceRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	attendanceReq, err := h.attendanceService.ApproveAttendanceRequest(c.Context(), uint(id), approver)
	if err != nil {
		if err == attendance.ErrAttendanceRequestNotFound {
			return response.NotFound(c, "attendance request not found")
		}
		return response.InternalServerError(c, "failed to approve attendance request")
	}

	return response.Success(c, attendanceReq, "Attendance request approved successfully")
}

// RejectAttendanceRequest rejects an attendance request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.reject_request
func (h *AttendanceHandler) RejectAttendanceRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	attendanceReq, err := h.attendanceService.RejectAttendanceRequest(c.Context(), uint(id), approver)
	if err != nil {
		if err == attendance.ErrAttendanceRequestNotFound {
			return response.NotFound(c, "attendance request not found")
		}
		return response.InternalServerError(c, "failed to reject attendance request")
	}

	return response.Success(c, attendanceReq, "Attendance request rejected successfully")
}
