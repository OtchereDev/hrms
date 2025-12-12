package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// HolidayRepository handles holiday data operations
type HolidayRepository struct {
	db *gorm.DB
}

// NewHolidayRepository creates a new holiday repository
func NewHolidayRepository(db *gorm.DB) *HolidayRepository {
	return &HolidayRepository{db: db}
}

// GetHolidayList retrieves a holiday list by name
func (r *HolidayRepository) GetHolidayList(ctx context.Context, name string) (*hr.HolidayList, error) {
	var holidayList hr.HolidayList
	err := r.db.WithContext(ctx).
		Preload("Holidays").
		Where("holiday_list_name = ?", name).
		First(&holidayList).Error

	if err != nil {
		return nil, err
	}

	return &holidayList, nil
}

// ListHolidayLists retrieves all holiday lists
func (r *HolidayRepository) ListHolidayLists(ctx context.Context, country string) ([]hr.HolidayList, error) {
	var lists []hr.HolidayList
	query := r.db.WithContext(ctx).Preload("Holidays")

	if country != "" {
		query = query.Where("country = ?", country)
	}

	err := query.Find(&lists).Error
	return lists, err
}

// GetHolidays retrieves holidays for a specific list within a date range
func (r *HolidayRepository) GetHolidays(ctx context.Context, holidayListName string, fromDate, toDate time.Time) ([]hr.Holiday, error) {
	var holidayList hr.HolidayList
	err := r.db.WithContext(ctx).
		Preload("Holidays", "holiday_date >= ? AND holiday_date <= ?", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")).
		Where("holiday_list_name = ?", holidayListName).
		First(&holidayList).Error

	if err != nil {
		return nil, err
	}

	return holidayList.Holidays, nil
}

// IsHoliday checks if a specific date is a holiday
func (r *HolidayRepository) IsHoliday(ctx context.Context, holidayListName string, date time.Time) (bool, *hr.Holiday, error) {
	var holidayList hr.HolidayList
	err := r.db.WithContext(ctx).
		Where("holiday_list_name = ?", holidayListName).
		First(&holidayList).Error

	if err != nil {
		return false, nil, err
	}

	// Check if date is a weekly off
	weekday := date.Weekday()
	if (weekday == time.Sunday && holidayList.WeeklySunday) ||
	   (weekday == time.Monday && holidayList.WeeklyMonday) ||
	   (weekday == time.Tuesday && holidayList.WeeklyTuesday) ||
	   (weekday == time.Wednesday && holidayList.WeeklyWednesday) ||
	   (weekday == time.Thursday && holidayList.WeeklyThursday) ||
	   (weekday == time.Friday && holidayList.WeeklyFriday) ||
	   (weekday == time.Saturday && holidayList.WeeklySaturday) {
		return true, &hr.Holiday{
			HolidayDate: date.Format("2006-01-02"),
			Description: "Weekly Off",
			WeeklyOff:   true,
		}, nil
	}

	// Check if date is in holiday list
	var holiday hr.Holiday
	err = r.db.WithContext(ctx).
		Where("parent_id = ? AND holiday_date = ?", holidayList.ID, date.Format("2006-01-02")).
		First(&holiday).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, err
	}

	return true, &holiday, nil
}

// CreateHolidayList creates a new holiday list
func (r *HolidayRepository) CreateHolidayList(ctx context.Context, holidayList *hr.HolidayList) error {
	return r.db.WithContext(ctx).Create(holidayList).Error
}

// UpdateHolidayList updates an existing holiday list
func (r *HolidayRepository) UpdateHolidayList(ctx context.Context, holidayList *hr.HolidayList) error {
	return r.db.WithContext(ctx).Save(holidayList).Error
}

// DeleteHolidayList deletes a holiday list
func (r *HolidayRepository) DeleteHolidayList(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.HolidayList{}, id).Error
}
