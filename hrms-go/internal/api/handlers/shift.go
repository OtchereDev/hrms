package handlers

import (
	"strconv"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/services/shift"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// ShiftHandler handles shift-related requests
type ShiftHandler struct {
	shiftService *shift.ShiftService
}

// NewShiftHandler creates a new shift handler
func NewShiftHandler(shiftService *shift.ShiftService) *ShiftHandler {
	return &ShiftHandler{
		shiftService: shiftService,
	}
}

// GetShiftType retrieves a shift type by name
// GET /api/method/hrms.hr.doctype.shift_type.shift_type.get
func (h *ShiftHandler) GetShiftType(c *fiber.Ctx) error {
	shiftName := c.Query("shift_name")
	if shiftName == "" {
		return response.BadRequest(c, "shift_name is required")
	}

	shift, err := h.shiftService.GetShiftType(c.Context(), shiftName)
	if err != nil {
		return response.NotFound(c, "shift type not found")
	}

	return response.Success(c, shift, "Shift type retrieved successfully")
}

// ListShiftTypes retrieves all shift types
// GET /api/method/hrms.hr.doctype.shift_type.shift_type.list
func (h *ShiftHandler) ListShiftTypes(c *fiber.Ctx) error {
	shifts, err := h.shiftService.ListShiftTypes(c.Context())
	if err != nil {
		return response.InternalServerError(c, "failed to list shift types")
	}

	return response.Success(c, shifts, "Shift types retrieved successfully")
}

// GetShiftDetails retrieves detailed information about a shift
// GET /api/method/hrms.hr.doctype.shift_type.shift_type.get_shift_details
func (h *ShiftHandler) GetShiftDetails(c *fiber.Ctx) error {
	shiftName := c.Query("shift_name")
	if shiftName == "" {
		return response.BadRequest(c, "shift_name is required")
	}

	details, err := h.shiftService.GetShiftDetails(c.Context(), shiftName)
	if err != nil {
		return response.NotFound(c, "shift type not found")
	}

	return response.Success(c, details, "Shift details retrieved successfully")
}

// AssignShift assigns a shift to an employee
// POST /api/method/hrms.hr.doctype.shift_assignment.shift_assignment.assign_shift
func (h *ShiftHandler) AssignShift(c *fiber.Ctx) error {
	employee := c.Query("employee")
	shiftType := c.Query("shift_type")
	company := c.Query("company")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	approvedBy := c.Query("approved_by")

	if employee == "" || shiftType == "" || company == "" || fromDateStr == "" {
		return response.BadRequest(c, "employee, shift_type, company, and from_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	var toDate *time.Time
	if toDateStr != "" {
		td, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
		}
		toDate = &td
	}

	assignment, err := h.shiftService.AssignShift(c.Context(), employee, shiftType, company, fromDate, toDate, approvedBy)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, assignment, "Shift assigned successfully")
}

// GetCurrentShift retrieves the current shift for an employee
// GET /api/method/hrms.hr.doctype.shift_assignment.shift_assignment.get_current_shift
func (h *ShiftHandler) GetCurrentShift(c *fiber.Ctx) error {
	employee := c.Query("employee")
	dateStr := c.Query("date")

	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return response.BadRequest(c, "invalid date format, expected YYYY-MM-DD")
		}
	} else {
		date = time.Now()
	}

	assignment, shiftType, err := h.shiftService.GetCurrentShift(c.Context(), employee, date)
	if err != nil {
		return response.NotFound(c, "no active shift found for employee")
	}

	result := map[string]interface{}{
		"shift_assignment": assignment,
		"shift_type":       shiftType,
	}

	return response.Success(c, result, "Current shift retrieved successfully")
}

// GetShiftAssignmentsForEmployee retrieves all shift assignments for an employee
// GET /api/method/hrms.hr.doctype.shift_assignment.shift_assignment.get_for_employee
func (h *ShiftHandler) GetShiftAssignmentsForEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	status := c.Query("status")

	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	assignments, err := h.shiftService.GetShiftAssignmentsForEmployee(c.Context(), employee, status)
	if err != nil {
		return response.InternalServerError(c, "failed to get shift assignments")
	}

	return response.Success(c, assignments, "Shift assignments retrieved successfully")
}

// CreateShiftRequest creates a new shift request
// POST /api/method/hrms.hr.doctype.shift_request.shift_request.create
func (h *ShiftHandler) CreateShiftRequest(c *fiber.Ctx) error {
	employee := c.Query("employee")
	shiftType := c.Query("shift_type")
	company := c.Query("company")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	reason := c.Query("reason")

	if employee == "" || shiftType == "" || company == "" || fromDateStr == "" {
		return response.BadRequest(c, "employee, shift_type, company, and from_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	var toDate *time.Time
	if toDateStr != "" {
		td, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
		}
		toDate = &td
	}

	request := &hr.ShiftRequest{
		Employee:  employee,
		ShiftType: shiftType,
		Company:   company,
		FromDate:  fromDate,
		ToDate:    toDate,
	}

	err = h.shiftService.CreateShiftRequest(c.Context(), request)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, request, "Shift request created successfully")
}

// ApproveShiftRequest approves a shift request
// POST /api/method/hrms.hr.doctype.shift_request.shift_request.approve
func (h *ShiftHandler) ApproveShiftRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	approvedBy := c.Query("approved_by")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	request, err := h.shiftService.ApproveShiftRequest(c.Context(), uint(id), approvedBy)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, request, "Shift request approved successfully")
}

// RejectShiftRequest rejects a shift request
// POST /api/method/hrms.hr.doctype.shift_request.shift_request.reject
func (h *ShiftHandler) RejectShiftRequest(c *fiber.Ctx) error {
	idStr := c.Query("id")
	rejectedBy := c.Query("rejected_by")
	reason := c.Query("reason")

	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	request, err := h.shiftService.RejectShiftRequest(c.Context(), uint(id), rejectedBy, reason)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, request, "Shift request rejected successfully")
}

// ListShiftRequests retrieves shift requests with filters
// GET /api/method/hrms.hr.doctype.shift_request.shift_request.list
func (h *ShiftHandler) ListShiftRequests(c *fiber.Ctx) error {
	employee := c.Query("employee")
	status := c.Query("status")

	requests, err := h.shiftService.ListShiftRequests(c.Context(), employee, status)
	if err != nil {
		return response.InternalServerError(c, "failed to list shift requests")
	}

	return response.Success(c, requests, "Shift requests retrieved successfully")
}
