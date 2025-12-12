package handlers

import (
	"time"

	frappeResponse "github.com/OtchereDev/hrms-go/internal/core/frappe"
	"github.com/OtchereDev/hrms-go/internal/core/services/holiday"
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
		return frappeResponse.SendBadRequest(c, "holiday list name is required")
	}

	holidayList, err := h.holidayService.GetHolidayList(c.Context(), name)
	if err != nil {
		return frappeResponse.SendNotFound(c, "holiday list not found")
	}

	return frappeResponse.SendSuccess(c, holidayList)
}

// ListHolidayLists retrieves all holiday lists
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.list
func (h *HolidayHandler) ListHolidayLists(c *fiber.Ctx) error {
	country := c.Query("country")

	lists, err := h.holidayService.ListHolidayLists(c.Context(), country)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to list holiday lists")
	}

	return frappeResponse.SendSuccess(c, lists)
}

// GetHolidays retrieves holidays within a date range
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_holidays
func (h *HolidayHandler) GetHolidays(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return frappeResponse.SendBadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	holidays, err := h.holidayService.GetHolidays(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get holidays")
	}

	return frappeResponse.SendSuccess(c, holidays)
}

// IsHoliday checks if a specific date is a holiday
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.is_holiday
func (h *HolidayHandler) IsHoliday(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	dateStr := c.Query("date")

	if holidayListName == "" || dateStr == "" {
		return frappeResponse.SendBadRequest(c, "holiday_list and date are required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid date format, expected YYYY-MM-DD")
	}

	isHoliday, holiday, err := h.holidayService.IsHoliday(c.Context(), holidayListName, date)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to check holiday")
	}

	result := map[string]interface{}{
		"is_holiday": isHoliday,
		"holiday":    holiday,
	}

	return frappeResponse.SendSuccess(c, result)
}

// GetWorkingDays calculates working days between two dates
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_working_days
func (h *HolidayHandler) GetWorkingDays(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return frappeResponse.SendBadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	workingDays, err := h.holidayService.GetWorkingDays(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to calculate working days")
	}

	result := map[string]interface{}{
		"working_days": workingDays,
	}

	return frappeResponse.SendSuccess(c, result)
}

// GetHolidaysBetweenDates returns all holidays including weekly offs
// GET /api/method/hrms.hr.doctype.holiday_list.holiday_list.get_holidays_between_dates
func (h *HolidayHandler) GetHolidaysBetweenDates(c *fiber.Ctx) error {
	holidayListName := c.Query("holiday_list")
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")

	if holidayListName == "" || fromDateStr == "" || toDateStr == "" {
		return frappeResponse.SendBadRequest(c, "holiday_list, from_date, and to_date are required")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid from_date format, expected YYYY-MM-DD")
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid to_date format, expected YYYY-MM-DD")
	}

	holidays, err := h.holidayService.GetHolidaysBetweenDates(c.Context(), holidayListName, fromDate, toDate)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get holidays")
	}

	return frappeResponse.SendSuccess(c, holidays)
}
