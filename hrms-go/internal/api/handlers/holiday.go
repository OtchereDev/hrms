package handlers

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/services/holiday"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// HolidayHandler handles holiday-related requests
type HolidayHandler struct {
	holidayService *holiday.HolidayService
}

// NewHolidayHandler creates a new holiday handler
func NewHolidayHandler(holidayService *holiday.HolidayService) *HolidayHandler {
	return &HolidayHandler{
		holidayService: holidayService,
	}
}

// GetHolidayList retrieves a holiday list by name
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get
func (h *HolidayHandler) GetHolidayList(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return response.BadRequest(c, "holiday list name is required")
	}

	holidayList, err := h.holidayService.GetHolidayList(c.Context(), name)
	if err != nil {
		return response.NotFound(c, "holiday list not found")
	}

	return response.Success(c, holidayList, "Holiday list retrieved successfully")
}

// ListHolidayLists retrieves all holiday lists
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.list
func (h *HolidayHandler) ListHolidayLists(c *fiber.Ctx) error {
	country := c.Query("country")

	lists, err := h.holidayService.ListHolidayLists(c.Context(), country)
	if err != nil {
		return response.InternalServerError(c, "failed to list holiday lists")
	}

	return response.Success(c, lists, "Holiday lists retrieved successfully")
}

// GetHolidays retrieves holidays within a date range
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_holidays
func (h *HolidayHandler) GetHolidays(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return response.BadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	holidays, err := h.holidayService.GetHolidays(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return response.InternalServerError(c, "failed to get holidays")
	}

	return response.Success(c, holidays, "Holidays retrieved successfully")
}

// IsHoliday checks if a specific date is a holiday
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.is_holiday
func (h *HolidayHandler) IsHoliday(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	dateStr := c.Query("date")

	if holidayListName == "" || dateStr == "" {
		return response.BadRequest(c, "holiday_list and date are required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return response.BadRequest(c, "invalid date format, expected YYYY-MM-DD")
	}

	isHoliday, holiday, err := h.holidayService.IsHoliday(c.Context(), holidayListName, date)
	if err != nil {
		return response.InternalServerError(c, "failed to check holiday")
	}

	result := map[string]interface{}{
		"is_holiday": isHoliday,
		"holiday":    holiday,
	}

	return response.Success(c, result, "Holiday check completed successfully")
}

// GetWorkingDays calculates working days between two dates
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_working_days
func (h *HolidayHandler) GetWorkingDays(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return response.BadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	workingDays, err := h.holidayService.GetWorkingDays(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return response.InternalServerError(c, "failed to calculate working days")
	}

	result := map[string]interface{}{
		"working_days": workingDays,
	}

	return response.Success(c, result, "Working days calculated successfully")
}

// GetHolidaysBetweenDates returns all holidays including weekly offs
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_holidays_between_dates
func (h *HolidayHandler) GetHolidaysBetweenDates(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return response.BadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return response.BadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	holidays, err := h.holidayService.GetHolidaysBetweenDates(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return response.InternalServerError(c, "failed to get holidays")
	}

	return response.Success(c, holidays, "Holidays retrieved successfully")
}
