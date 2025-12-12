package holiday

import (
	"context"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"gorm.io/gorm"
)

// HolidayService handles holiday business logic
type HolidayService struct {
	holidayRepo *repositories.HolidayRepository
	db          *gorm.DB
}

// NewHolidayService creates a new holiday service
func NewHolidayService(db *gorm.DB) *HolidayService {
	return &HolidayService{
		holidayRepo: repositories.NewHolidayRepository(db),
		db:          db,
	}
}

// GetHolidayList retrieves a holiday list by name
func (s *HolidayService) GetHolidayList(ctx context.Context, name string) (*hr.HolidayList, error) {
	return s.holidayRepo.GetHolidayList(ctx, name)
}

// ListHolidayLists retrieves all holiday lists
func (s *HolidayService) ListHolidayLists(ctx context.Context, country string) ([]hr.HolidayList, error) {
	return s.holidayRepo.ListHolidayLists(ctx, country)
}

// GetHolidays retrieves holidays within a date range
func (s *HolidayService) GetHolidays(ctx context.Context, holidayListName string, fromDate, toDate time.Time) ([]hr.Holiday, error) {
	return s.holidayRepo.GetHolidays(ctx, holidayListName, fromDate, toDate)
}

// IsHoliday checks if a specific date is a holiday
func (s *HolidayService) IsHoliday(ctx context.Context, holidayListName string, date time.Time) (bool, *hr.Holiday, error) {
	return s.holidayRepo.IsHoliday(ctx, holidayListName, date)
}

// GetWorkingDays calculates the number of working days between two dates
func (s *HolidayService) GetWorkingDays(ctx context.Context, holidayListName string, fromDate, toDate time.Time) (int, error) {
	if fromDate.After(toDate) {
		return 0, fmt.Errorf("from_date cannot be after to_date")
	}

	// Get all holidays in the date range
	holidays, err := s.holidayRepo.GetHolidays(ctx, holidayListName, fromDate, toDate)
	if err != nil {
		return 0, err
	}

	// Create a map of holiday dates for quick lookup
	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.HolidayDate] = true
	}

	// Get holiday list for weekly offs
	holidayList, err := s.holidayRepo.GetHolidayList(ctx, holidayListName)
	if err != nil {
		return 0, err
	}

	// Count working days
	workingDays := 0
	currentDate := fromDate

	for currentDate.Before(toDate) || currentDate.Equal(toDate) {
		// Check if it's a weekly off
		weekday := currentDate.Weekday()
		isWeeklyOff := (weekday == time.Sunday && holidayList.WeeklySunday) ||
			(weekday == time.Monday && holidayList.WeeklyMonday) ||
			(weekday == time.Tuesday && holidayList.WeeklyTuesday) ||
			(weekday == time.Wednesday && holidayList.WeeklyWednesday) ||
			(weekday == time.Thursday && holidayList.WeeklyThursday) ||
			(weekday == time.Friday && holidayList.WeeklyFriday) ||
			(weekday == time.Saturday && holidayList.WeeklySaturday)

		// Check if it's a holiday
		dateStr := currentDate.Format("2006-01-02")
		isHoliday := holidayMap[dateStr]

		// Count if it's not a weekly off and not a holiday
		if !isWeeklyOff && !isHoliday {
			workingDays++
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return workingDays, nil
}

// GetHolidaysBetweenDates returns all holidays (including weekly offs) between two dates
func (s *HolidayService) GetHolidaysBetweenDates(ctx context.Context, holidayListName string, fromDate, toDate time.Time) ([]hr.Holiday, error) {
	// Get holidays from database
	holidays, err := s.holidayRepo.GetHolidays(ctx, holidayListName, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	// Get holiday list for weekly offs
	holidayList, err := s.holidayRepo.GetHolidayList(ctx, holidayListName)
	if err != nil {
		return nil, err
	}

	// Add weekly offs
	allHolidays := make([]hr.Holiday, 0, len(holidays))
	allHolidays = append(allHolidays, holidays...)

	// Iterate through dates and add weekly offs
	currentDate := fromDate
	for currentDate.Before(toDate) || currentDate.Equal(toDate) {
		weekday := currentDate.Weekday()
		isWeeklyOff := (weekday == time.Sunday && holidayList.WeeklySunday) ||
			(weekday == time.Monday && holidayList.WeeklyMonday) ||
			(weekday == time.Tuesday && holidayList.WeeklyTuesday) ||
			(weekday == time.Wednesday && holidayList.WeeklyWednesday) ||
			(weekday == time.Thursday && holidayList.WeeklyThursday) ||
			(weekday == time.Friday && holidayList.WeeklyFriday) ||
			(weekday == time.Saturday && holidayList.WeeklySaturday)

		if isWeeklyOff {
			// Check if this date is not already in holidays
			dateStr := currentDate.Format("2006-01-02")
			exists := false
			for _, h := range holidays {
				if h.HolidayDate == dateStr {
					exists = true
					break
				}
			}

			if !exists {
				allHolidays = append(allHolidays, hr.Holiday{
					HolidayDate: dateStr,
					Description: "Weekly Off",
					WeeklyOff:   true,
				})
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return allHolidays, nil
}

// CreateHolidayList creates a new holiday list
func (s *HolidayService) CreateHolidayList(ctx context.Context, holidayList *hr.HolidayList) error {
	// Count total holidays
	holidayList.TotalHolidays = len(holidayList.Holidays)
	return s.holidayRepo.CreateHolidayList(ctx, holidayList)
}

// UpdateHolidayList updates an existing holiday list
func (s *HolidayService) UpdateHolidayList(ctx context.Context, holidayList *hr.HolidayList) error {
	// Recalculate total holidays
	holidayList.TotalHolidays = len(holidayList.Holidays)
	return s.holidayRepo.UpdateHolidayList(ctx, holidayList)
}

// DeleteHolidayList deletes a holiday list
func (s *HolidayService) DeleteHolidayList(ctx context.Context, id uint) error {
	return s.holidayRepo.DeleteHolidayList(ctx, id)
}
