package handlers

import (
	"encoding/json"
	frappeResponse "github.com/OtchereDev/hrms-go/internal/core/frappe"
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
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	att, err := h.attendanceService.MarkAttendance(c.Context(), &req)
	if err != nil {
		switch err {
		case attendance.ErrAttendanceExists:
			return response.Conflict(c, "attendance already exists for this date")
		case attendance.ErrEmployeeRequired, attendance.ErrAttendanceDateRequired, attendance.ErrInvalidStatus:
			return frappeResponse.SendBadRequest(c, err.Error())
		default:
			return frappeResponse.SendInternalError(c, "failed to mark attendance")
		}
	}

	return response.Created(c, att, "Attendance marked successfully")
}

// GetAttendance retrieves an attendance record
// GET /api/method/hrms.hr.doctype.attendance.attendance.get_attendance
func (h *AttendanceHandler) GetAttendance(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	att, err := h.attendanceService.GetAttendanceByID(c.Context(), uint(id))
	if err != nil {
		if err == attendance.ErrAttendanceNotFound {
			return frappeResponse.SendNotFound(c, "attendance record not found")
		}
		return frappeResponse.SendInternalError(c, "failed to get attendance")
	}

	return frappeResponse.SendSuccess(c, att)
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
		return frappeResponse.SendInternalError(c, "failed to list attendance")
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
		return frappeResponse.SendBadRequest(c, "employee, year, and month are required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid year")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid month")
	}

	attendances, err := h.attendanceService.GetMonthlyAttendance(c.Context(), employee, year, month)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get monthly attendance")
	}

	return frappeResponse.SendSuccess(c, attendances)
}

// Checkin handles employee check-in/check-out
// POST /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.checkin
func (h *AttendanceHandler) Checkin(c *fiber.Ctx) error {
	var req attendance.CheckinRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	checkin, err := h.attendanceService.Checkin(c.Context(), &req)
	if err != nil {
		switch err {
		case attendance.ErrEmployeeRequired:
			return frappeResponse.SendBadRequest(c, "employee is required")
		case attendance.ErrNotCheckedIn:
			return frappeResponse.SendBadRequest(c, "employee not checked in")
		case attendance.ErrAlreadyCheckedIn:
			return response.Conflict(c, "employee already checked in")
		default:
			return frappeResponse.SendInternalError(c, "failed to process checkin")
		}
	}

	return response.Created(c, checkin, "Checkin recorded successfully")
}

// GetTodayCheckins retrieves today's check-ins for an employee
// GET /api/method/hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins
func (h *AttendanceHandler) GetTodayCheckins(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return frappeResponse.SendBadRequest(c, "employee is required")
	}

	checkins, err := h.attendanceService.GetTodayCheckinsForEmployee(c.Context(), employee)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get checkins")
	}

	return frappeResponse.SendSuccess(c, checkins)
}

// CreateAttendanceRequest creates an attendance correction request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.create_request
func (h *AttendanceHandler) CreateAttendanceRequest(c *fiber.Ctx) error {
	var req attendance.AttendanceRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	attendanceReq, err := h.attendanceService.CreateAttendanceRequest(c.Context(), &req)
	if err != nil {
		if err == attendance.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, err.Error())
		}
		return frappeResponse.SendInternalError(c, "failed to create attendance request")
	}

	return response.Created(c, attendanceReq, "Attendance request created successfully")
}

// ApproveAttendanceRequest approves an attendance request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.approve_request
func (h *AttendanceHandler) ApproveAttendanceRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")

	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	attendanceReq, err := h.attendanceService.ApproveAttendanceRequest(c.Context(), uint(id), approver)
	if err != nil {
		if err == attendance.ErrAttendanceRequestNotFound {
			return frappeResponse.SendNotFound(c, "attendance request not found")
		}
		return frappeResponse.SendInternalError(c, "failed to approve attendance request")
	}

	return frappeResponse.SendSuccess(c, attendanceReq)
}

// RejectAttendanceRequest rejects an attendance request
// POST /api/method/hrms.hr.doctype.attendance_request.attendance_request.reject_request
func (h *AttendanceHandler) RejectAttendanceRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approver := c.Query("approver")

	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	attendanceReq, err := h.attendanceService.RejectAttendanceRequest(c.Context(), uint(id), approver)
	if err != nil {
		if err == attendance.ErrAttendanceRequestNotFound {
			return frappeResponse.SendNotFound(c, "attendance request not found")
		}
		return frappeResponse.SendInternalError(c, "failed to reject attendance request")
	}

	return frappeResponse.SendSuccess(c, attendanceReq)
}

// GetUnmarkedDays returns days without attendance records
// GET /api/method/hrms.hr.doctype.attendance.attendance.get_unmarked_days
func (h *AttendanceHandler) GetUnmarkedDays(c *fiber.Ctx) error {
	employee := c.Query("employee")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	excludeHolidays := c.Query("exclude_holidays") == "1"

	if employee == "" || fromDateStr == "" || toDateStr == "" {
		return frappeResponse.SendBadRequest(c, "employee, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	unmarkedDays, err := h.attendanceService.GetUnmarkedDays(c.Context(), employee, fromDate, toDate, excludeHolidays)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get unmarked days")
	}

	return frappeResponse.SendSuccess(c, unmarkedDays)
}

// MarkBulkAttendance marks attendance for multiple days
// POST /api/method/hrms.hr.doctype.attendance.attendance.mark_bulk_attendance
func (h *AttendanceHandler) MarkBulkAttendance(c *fiber.Ctx) error {
	var req attendance.MarkBulkAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	if err := h.attendanceService.MarkBulkAttendance(c.Context(), &req); err != nil {
		return frappeResponse.SendInternalError(c, "failed to mark bulk attendance")
	}

	return frappeResponse.SendSuccess(c, nil)
}

// GetEvents returns attendance events for calendar view
// GET /api/method/hrms.hr.doctype.attendance.attendance.get_events
func (h *AttendanceHandler) GetEvents(c *fiber.Ctx) error {
	employee := c.Query("employee")
	startStr := c.Query("start")
	endStr := c.Query("end")
	filtersStr := c.Query("filters")

	if employee == "" || startStr == "" || endStr == "" {
		return frappeResponse.SendBadRequest(c, "employee, start, and end are required")
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid start date format, expected YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid end date format, expected YYYY-MM-DD")
	}

	// Parse filters if provided
	filters := make(map[string]interface{})
	if filtersStr != "" {
		if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
			return frappeResponse.SendBadRequest(c, "invalid filters format")
		}
	}

	events, err := h.attendanceService.GetEvents(c.Context(), employee, start, end, filters)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get events")
	}

	return frappeResponse.SendSuccess(c, events)
}
