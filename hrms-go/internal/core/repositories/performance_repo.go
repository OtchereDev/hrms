package repositories

import (
	"context"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"gorm.io/gorm"
)

// AppraisalRepository handles appraisal data access
type AppraisalRepository struct {
	db *gorm.DB
}

// NewAppraisalRepository creates a new appraisal repository
func NewAppraisalRepository(db *gorm.DB) *AppraisalRepository {
	return &AppraisalRepository{db: db}
}

// Create creates a new appraisal
func (r *AppraisalRepository) Create(ctx context.Context, appraisal *hr.Appraisal) error {
	return r.db.WithContext(ctx).Create(appraisal).Error
}

// GetByID retrieves an appraisal by ID
func (r *AppraisalRepository) GetByID(ctx context.Context, id uint) (*hr.Appraisal, error) {
	var appraisal hr.Appraisal
	err := r.db.WithContext(ctx).Preload("Goals").First(&appraisal, id).Error
	if err != nil {
		return nil, err
	}
	return &appraisal, nil
}

// Update updates an appraisal
func (r *AppraisalRepository) Update(ctx context.Context, appraisal *hr.Appraisal) error {
	return r.db.WithContext(ctx).Save(appraisal).Error
}

// Delete soft deletes an appraisal
func (r *AppraisalRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.Appraisal{}, id).Error
}

// List retrieves appraisals with pagination and filters
func (r *AppraisalRepository) List(ctx context.Context, filters AppraisalFilters, page, pageSize int) ([]hr.Appraisal, int64, error) {
	var appraisals []hr.Appraisal
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.Appraisal{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.AppraisalCycle != "" {
		query = query.Where("appraisal_cycle = ?", filters.AppraisalCycle)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("start_date >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("end_date <= ?", filters.EndDate)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("start_date DESC").Find(&appraisals).Error

	return appraisals, total, err
}

// AppraisalFilters represents filters for appraisal queries
type AppraisalFilters struct {
	Employee       string
	Company        string
	Department     string
	Status         string
	AppraisalCycle string
	StartDate      time.Time
	EndDate        time.Time
}

// GoalRepository handles goal data access
type GoalRepository struct {
	db *gorm.DB
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

// Create creates a new goal
func (r *GoalRepository) Create(ctx context.Context, goal *hr.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

// GetByID retrieves a goal by ID
func (r *GoalRepository) GetByID(ctx context.Context, id uint) (*hr.Goal, error) {
	var goal hr.Goal
	err := r.db.WithContext(ctx).First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// Update updates a goal
func (r *GoalRepository) Update(ctx context.Context, goal *hr.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

// Delete soft deletes a goal
func (r *GoalRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&hr.Goal{}, id).Error
}

// List retrieves goals with pagination and filters
func (r *GoalRepository) List(ctx context.Context, filters GoalFilters, page, pageSize int) ([]hr.Goal, int64, error) {
	var goals []hr.Goal
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.Goal{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.AssignedBy != "" {
		query = query.Where("assigned_by = ?", filters.AssignedBy)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("start_date DESC").Find(&goals).Error

	return goals, total, err
}

// GetActiveGoalsForEmployee retrieves active goals for an employee
func (r *GoalRepository) GetActiveGoalsForEmployee(ctx context.Context, employee string) ([]hr.Goal, error) {
	var goals []hr.Goal
	err := r.db.WithContext(ctx).
		Where("employee = ? AND status NOT IN ?", employee, []string{"Archived", "Failed"}).
		Order("start_date DESC").
		Find(&goals).Error
	return goals, err
}

// GoalFilters represents filters for goal queries
type GoalFilters struct {
	Employee   string
	Company    string
	Status     string
	AssignedBy string
}

// EmployeePerformanceFeedbackRepository handles 360-degree feedback
type EmployeePerformanceFeedbackRepository struct {
	db *gorm.DB
}

// NewEmployeePerformanceFeedbackRepository creates a new feedback repository
func NewEmployeePerformanceFeedbackRepository(db *gorm.DB) *EmployeePerformanceFeedbackRepository {
	return &EmployeePerformanceFeedbackRepository{db: db}
}

// Create creates a new feedback record
func (r *EmployeePerformanceFeedbackRepository) Create(ctx context.Context, feedback *hr.EmployeePerformanceFeedback) error {
	return r.db.WithContext(ctx).Create(feedback).Error
}

// GetByID retrieves a feedback record by ID
func (r *EmployeePerformanceFeedbackRepository) GetByID(ctx context.Context, id uint) (*hr.EmployeePerformanceFeedback, error) {
	var feedback hr.EmployeePerformanceFeedback
	err := r.db.WithContext(ctx).Preload("FeedbackCriteria").First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// Update updates a feedback record
func (r *EmployeePerformanceFeedbackRepository) Update(ctx context.Context, feedback *hr.EmployeePerformanceFeedback) error {
	return r.db.WithContext(ctx).Save(feedback).Error
}

// List retrieves feedback records with filters
func (r *EmployeePerformanceFeedbackRepository) List(ctx context.Context, filters FeedbackFilters, page, pageSize int) ([]hr.EmployeePerformanceFeedback, int64, error) {
	var feedbacks []hr.EmployeePerformanceFeedback
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.EmployeePerformanceFeedback{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Reviewer != "" {
		query = query.Where("reviewer = ?", filters.Reviewer)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.AppraisalCycle != "" {
		query = query.Where("appraisal_cycle = ?", filters.AppraisalCycle)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("feedback_date DESC").Find(&feedbacks).Error

	return feedbacks, total, err
}

// GetFeedbackForEmployee retrieves all feedback for an employee
func (r *EmployeePerformanceFeedbackRepository) GetFeedbackForEmployee(ctx context.Context, employee string, cycle string) ([]hr.EmployeePerformanceFeedback, error) {
	var feedbacks []hr.EmployeePerformanceFeedback
	query := r.db.WithContext(ctx).Where("employee = ?", employee)

	if cycle != "" {
		query = query.Where("appraisal_cycle = ?", cycle)
	}

	err := query.Order("feedback_date DESC").Find(&feedbacks).Error
	return feedbacks, err
}

// FeedbackFilters represents filters for feedback queries
type FeedbackFilters struct {
	Employee       string
	Reviewer       string
	Company        string
	AppraisalCycle string
}

// EmployeeSkillMapRepository handles skill map data access
type EmployeeSkillMapRepository struct {
	db *gorm.DB
}

// NewEmployeeSkillMapRepository creates a new skill map repository
func NewEmployeeSkillMapRepository(db *gorm.DB) *EmployeeSkillMapRepository {
	return &EmployeeSkillMapRepository{db: db}
}

// Create creates a new skill map
func (r *EmployeeSkillMapRepository) Create(ctx context.Context, skillMap *hr.EmployeeSkillMap) error {
	return r.db.WithContext(ctx).Create(skillMap).Error
}

// GetByID retrieves a skill map by ID
func (r *EmployeeSkillMapRepository) GetByID(ctx context.Context, id uint) (*hr.EmployeeSkillMap, error) {
	var skillMap hr.EmployeeSkillMap
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Trainings").
		First(&skillMap, id).Error
	if err != nil {
		return nil, err
	}
	return &skillMap, nil
}

// GetByEmployee retrieves skill map for an employee
func (r *EmployeeSkillMapRepository) GetByEmployee(ctx context.Context, employee string) (*hr.EmployeeSkillMap, error) {
	var skillMap hr.EmployeeSkillMap
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Trainings").
		Where("employee = ?", employee).
		First(&skillMap).Error
	if err != nil {
		return nil, err
	}
	return &skillMap, nil
}

// Update updates a skill map
func (r *EmployeeSkillMapRepository) Update(ctx context.Context, skillMap *hr.EmployeeSkillMap) error {
	return r.db.WithContext(ctx).Save(skillMap).Error
}

// List retrieves skill maps with filters
func (r *EmployeeSkillMapRepository) List(ctx context.Context, filters SkillMapFilters, page, pageSize int) ([]hr.EmployeeSkillMap, int64, error) {
	var skillMaps []hr.EmployeeSkillMap
	var total int64

	query := r.db.WithContext(ctx).Model(&hr.EmployeeSkillMap{})

	// Apply filters
	if filters.Employee != "" {
		query = query.Where("employee = ?", filters.Employee)
	}
	if filters.Company != "" {
		query = query.Where("company = ?", filters.Company)
	}
	if filters.Department != "" {
		query = query.Where("department = ?", filters.Department)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&skillMaps).Error

	return skillMaps, total, err
}

// SkillMapFilters represents filters for skill map queries
type SkillMapFilters struct {
	Employee   string
	Company    string
	Department string
}

// AppraisalTemplateRepository handles appraisal template data access
type AppraisalTemplateRepository struct {
	db *gorm.DB
}

// NewAppraisalTemplateRepository creates a new template repository
func NewAppraisalTemplateRepository(db *gorm.DB) *AppraisalTemplateRepository {
	return &AppraisalTemplateRepository{db: db}
}

// Create creates a new appraisal template
func (r *AppraisalTemplateRepository) Create(ctx context.Context, template *hr.AppraisalTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID retrieves an appraisal template by ID
func (r *AppraisalTemplateRepository) GetByID(ctx context.Context, id uint) (*hr.AppraisalTemplate, error) {
	var template hr.AppraisalTemplate
	err := r.db.WithContext(ctx).
		Preload("Goals").
		Preload("RatingScale").
		First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetByName retrieves an appraisal template by name
func (r *AppraisalTemplateRepository) GetByName(ctx context.Context, name string) (*hr.AppraisalTemplate, error) {
	var template hr.AppraisalTemplate
	err := r.db.WithContext(ctx).
		Preload("Goals").
		Preload("RatingScale").
		Where("template_name = ?", name).
		First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update updates an appraisal template
func (r *AppraisalTemplateRepository) Update(ctx context.Context, template *hr.AppraisalTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// ListActive retrieves all active appraisal templates
func (r *AppraisalTemplateRepository) ListActive(ctx context.Context, company string) ([]hr.AppraisalTemplate, error) {
	var templates []hr.AppraisalTemplate
	query := r.db.WithContext(ctx)

	if company != "" {
		query = query.Where("company = ?", company)
	}

	err := query.Order("template_name ASC").Find(&templates).Error
	return templates, err
}

// AppraisalCycleRepository handles appraisal cycle data access
type AppraisalCycleRepository struct {
	db *gorm.DB
}

// NewAppraisalCycleRepository creates a new cycle repository
func NewAppraisalCycleRepository(db *gorm.DB) *AppraisalCycleRepository {
	return &AppraisalCycleRepository{db: db}
}

// Create creates a new appraisal cycle
func (r *AppraisalCycleRepository) Create(ctx context.Context, cycle *hr.AppraisalCycle) error {
	return r.db.WithContext(ctx).Create(cycle).Error
}

// GetByID retrieves an appraisal cycle by ID
func (r *AppraisalCycleRepository) GetByID(ctx context.Context, id uint) (*hr.AppraisalCycle, error) {
	var cycle hr.AppraisalCycle
	err := r.db.WithContext(ctx).First(&cycle, id).Error
	if err != nil {
		return nil, err
	}
	return &cycle, nil
}

// GetByName retrieves an appraisal cycle by name
func (r *AppraisalCycleRepository) GetByName(ctx context.Context, name string) (*hr.AppraisalCycle, error) {
	var cycle hr.AppraisalCycle
	err := r.db.WithContext(ctx).Where("cycle_name = ?", name).First(&cycle).Error
	if err != nil {
		return nil, err
	}
	return &cycle, nil
}

// Update updates an appraisal cycle
func (r *AppraisalCycleRepository) Update(ctx context.Context, cycle *hr.AppraisalCycle) error {
	return r.db.WithContext(ctx).Save(cycle).Error
}

// GetActiveCycles retrieves all active appraisal cycles
func (r *AppraisalCycleRepository) GetActiveCycles(ctx context.Context, company string) ([]hr.AppraisalCycle, error) {
	var cycles []hr.AppraisalCycle
	query := r.db.WithContext(ctx).Where("status = ?", "Active")

	if company != "" {
		query = query.Where("company = ?", company)
	}

	err := query.Order("start_date DESC").Find(&cycles).Error
	return cycles, err
}
