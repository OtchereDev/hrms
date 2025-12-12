package performance

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/hr"
	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"gorm.io/gorm"
)

var (
	ErrAppraisalNotFound       = errors.New("appraisal not found")
	ErrGoalNotFound            = errors.New("goal not found")
	ErrFeedbackNotFound        = errors.New("feedback not found")
	ErrSkillMapNotFound        = errors.New("skill map not found")
	ErrTemplateNotFound        = errors.New("appraisal template not found")
	ErrCycleNotFound           = errors.New("appraisal cycle not found")
	ErrEmployeeRequired        = errors.New("employee is required")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidProgress         = errors.New("progress must be between 0 and 100")
)

// PerformanceService handles performance management business logic
type PerformanceService struct {
	appraisalRepo *repositories.AppraisalRepository
	goalRepo      *repositories.GoalRepository
	feedbackRepo  *repositories.EmployeePerformanceFeedbackRepository
	skillMapRepo  *repositories.EmployeeSkillMapRepository
	templateRepo  *repositories.AppraisalTemplateRepository
	cycleRepo     *repositories.AppraisalCycleRepository
	db            *gorm.DB
}

// NewPerformanceService creates a new performance service
func NewPerformanceService(db *gorm.DB) *PerformanceService {
	return &PerformanceService{
		appraisalRepo: repositories.NewAppraisalRepository(db),
		goalRepo:      repositories.NewGoalRepository(db),
		feedbackRepo:  repositories.NewEmployeePerformanceFeedbackRepository(db),
		skillMapRepo:  repositories.NewEmployeeSkillMapRepository(db),
		templateRepo:  repositories.NewAppraisalTemplateRepository(db),
		cycleRepo:     repositories.NewAppraisalCycleRepository(db),
		db:            db,
	}
}

// ========== Appraisal Operations ==========

// CreateAppraisalRequest represents a request to create an appraisal
type CreateAppraisalRequest struct {
	Employee          string    `json:"employee" validate:"required"`
	Company           string    `json:"company" validate:"required"`
	Department        string    `json:"department"`
	Designation       string    `json:"designation"`
	AppraisalCycle    string    `json:"appraisal_cycle"`
	StartDate         time.Time `json:"start_date" validate:"required"`
	EndDate           time.Time `json:"end_date" validate:"required"`
	AppraisalTemplate string    `json:"appraisal_template"`
}

// UpdateAppraisalRequest represents a request to update an appraisal
type UpdateAppraisalRequest struct {
	TotalScore      *float64 `json:"total_score"`
	SelfAppraisal   *string  `json:"self_appraisal"`
	ManagerFeedback *string  `json:"manager_feedback"`
	FinalScore      *float64 `json:"final_score"`
}

// CreateAppraisal creates a new appraisal
func (s *PerformanceService) CreateAppraisal(ctx context.Context, req *CreateAppraisalRequest) (*hr.Appraisal, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}
	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	// Get template if specified
	var goals []hr.AppraisalGoal
	if req.AppraisalTemplate != "" {
		template, err := s.templateRepo.GetByName(ctx, req.AppraisalTemplate)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrTemplateNotFound
			}
			return nil, fmt.Errorf("failed to get template: %w", err)
		}

		// Copy goals from template
		for _, tGoal := range template.Goals {
			goals = append(goals, hr.AppraisalGoal{
				KRA:       tGoal.KRA,
				KRATitle:  tGoal.KRATitle,
				Weightage: tGoal.Weightage,
			})
		}
	}

	appraisal := &hr.Appraisal{
		Employee:          req.Employee,
		Company:           req.Company,
		Department:        req.Department,
		Designation:       req.Designation,
		AppraisalCycle:    req.AppraisalCycle,
		StartDate:         req.StartDate,
		EndDate:           req.EndDate,
		AppraisalTemplate: req.AppraisalTemplate,
		Status:            "Draft",
		Goals:             goals,
	}

	if err := s.appraisalRepo.Create(ctx, appraisal); err != nil {
		return nil, fmt.Errorf("failed to create appraisal: %w", err)
	}

	return appraisal, nil
}

// GetAppraisalByID retrieves an appraisal by ID
func (s *PerformanceService) GetAppraisalByID(ctx context.Context, id uint) (*hr.Appraisal, error) {
	appraisal, err := s.appraisalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppraisalNotFound
		}
		return nil, fmt.Errorf("failed to get appraisal: %w", err)
	}
	return appraisal, nil
}

// UpdateAppraisal updates an appraisal
func (s *PerformanceService) UpdateAppraisal(ctx context.Context, id uint, req *UpdateAppraisalRequest) (*hr.Appraisal, error) {
	appraisal, err := s.appraisalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppraisalNotFound
		}
		return nil, fmt.Errorf("failed to get appraisal: %w", err)
	}

	// Update fields if provided
	if req.TotalScore != nil {
		appraisal.TotalScore = *req.TotalScore
	}
	if req.SelfAppraisal != nil {
		appraisal.SelfAppraisal = *req.SelfAppraisal
	}
	if req.ManagerFeedback != nil {
		appraisal.ManagerFeedback = *req.ManagerFeedback
	}
	if req.FinalScore != nil {
		appraisal.FinalScore = *req.FinalScore
	}

	if err := s.appraisalRepo.Update(ctx, appraisal); err != nil {
		return nil, fmt.Errorf("failed to update appraisal: %w", err)
	}

	return appraisal, nil
}

// SubmitAppraisal submits an appraisal
func (s *PerformanceService) SubmitAppraisal(ctx context.Context, id uint) (*hr.Appraisal, error) {
	appraisal, err := s.appraisalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppraisalNotFound
		}
		return nil, fmt.Errorf("failed to get appraisal: %w", err)
	}

	if appraisal.Status != "Draft" {
		return nil, fmt.Errorf("cannot submit appraisal with status: %s", appraisal.Status)
	}

	appraisal.Status = "Submitted"

	if err := s.appraisalRepo.Update(ctx, appraisal); err != nil {
		return nil, fmt.Errorf("failed to submit appraisal: %w", err)
	}

	return appraisal, nil
}

// CompleteAppraisal completes an appraisal
func (s *PerformanceService) CompleteAppraisal(ctx context.Context, id uint, finalScore float64) (*hr.Appraisal, error) {
	appraisal, err := s.appraisalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppraisalNotFound
		}
		return nil, fmt.Errorf("failed to get appraisal: %w", err)
	}

	if appraisal.Status != "Submitted" {
		return nil, fmt.Errorf("cannot complete appraisal with status: %s", appraisal.Status)
	}

	appraisal.Status = "Completed"
	appraisal.FinalScore = finalScore

	if err := s.appraisalRepo.Update(ctx, appraisal); err != nil {
		return nil, fmt.Errorf("failed to complete appraisal: %w", err)
	}

	return appraisal, nil
}

// ListAppraisals retrieves appraisals with filters
func (s *PerformanceService) ListAppraisals(ctx context.Context, filters repositories.AppraisalFilters, page, pageSize int) ([]hr.Appraisal, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.appraisalRepo.List(ctx, filters, page, pageSize)
}

// ========== Goal Operations ==========

// CreateGoalRequest represents a request to create a goal
type CreateGoalRequest struct {
	GoalName      string    `json:"goal_name" validate:"required"`
	Description   string    `json:"description"`
	Employee      string    `json:"employee" validate:"required"`
	AssignedBy    string    `json:"assigned_by"`
	Company       string    `json:"company" validate:"required"`
	StartDate     time.Time `json:"start_date" validate:"required"`
	EndDate       time.Time `json:"end_date" validate:"required"`
	KPIGoal       bool      `json:"kpi_goal"`
	AlignWithGoal string    `json:"align_with_goal"`
}

// UpdateGoalRequest represents a request to update a goal
type UpdateGoalRequest struct {
	Description *string  `json:"description"`
	Progress    *float64 `json:"progress"`
	Status      *string  `json:"status"`
}

// CreateGoal creates a new goal
func (s *PerformanceService) CreateGoal(ctx context.Context, req *CreateGoalRequest) (*hr.Goal, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}
	if req.GoalName == "" {
		return nil, fmt.Errorf("goal name is required")
	}
	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	goal := &hr.Goal{
		GoalName:      req.GoalName,
		Description:   req.Description,
		Employee:      req.Employee,
		AssignedBy:    req.AssignedBy,
		Company:       req.Company,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Status:        "Pending",
		Progress:      0,
		KPIGoal:       req.KPIGoal,
		AlignWithGoal: req.AlignWithGoal,
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goal, nil
}

// GetGoalByID retrieves a goal by ID
func (s *PerformanceService) GetGoalByID(ctx context.Context, id uint) (*hr.Goal, error) {
	goal, err := s.goalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}
	return goal, nil
}

// UpdateGoal updates a goal
func (s *PerformanceService) UpdateGoal(ctx context.Context, id uint, req *UpdateGoalRequest) (*hr.Goal, error) {
	goal, err := s.goalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGoalNotFound
		}
		return nil, fmt.Errorf("failed to get goal: %w", err)
	}

	// Update fields if provided
	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.Progress != nil {
		if *req.Progress < 0 || *req.Progress > 100 {
			return nil, ErrInvalidProgress
		}
		goal.Progress = *req.Progress

		// Auto-update status based on progress
		if *req.Progress == 100 {
			goal.Status = "Achieved"
		} else if *req.Progress > 0 && *req.Progress < 100 {
			goal.Status = "Partially Achieved"
		}
	}
	if req.Status != nil {
		goal.Status = *req.Status
	}

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	return goal, nil
}

// ListGoals retrieves goals with filters
func (s *PerformanceService) ListGoals(ctx context.Context, filters repositories.GoalFilters, page, pageSize int) ([]hr.Goal, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	return s.goalRepo.List(ctx, filters, page, pageSize)
}

// GetActiveGoalsForEmployee retrieves active goals for an employee
func (s *PerformanceService) GetActiveGoalsForEmployee(ctx context.Context, employee string) ([]hr.Goal, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	return s.goalRepo.GetActiveGoalsForEmployee(ctx, employee)
}

// ========== 360-Degree Feedback Operations ==========

// CreateFeedbackRequest represents a request to create feedback
type CreateFeedbackRequest struct {
	Employee       string    `json:"employee" validate:"required"`
	Company        string    `json:"company" validate:"required"`
	Reviewer       string    `json:"reviewer" validate:"required"`
	FeedbackDate   time.Time `json:"feedback_date" validate:"required"`
	AppraisalCycle string    `json:"appraisal_cycle"`
	OverallFeedback string   `json:"overall_feedback"`
	OverallRating  float64   `json:"overall_rating"`
}

// CreateFeedback creates a new 360-degree feedback
func (s *PerformanceService) CreateFeedback(ctx context.Context, req *CreateFeedbackRequest) (*hr.EmployeePerformanceFeedback, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}
	if req.Reviewer == "" {
		return nil, fmt.Errorf("reviewer is required")
	}

	feedback := &hr.EmployeePerformanceFeedback{
		Employee:        req.Employee,
		Company:         req.Company,
		Reviewer:        req.Reviewer,
		FeedbackDate:    req.FeedbackDate,
		AppraisalCycle:  req.AppraisalCycle,
		OverallFeedback: req.OverallFeedback,
		OverallRating:   req.OverallRating,
	}

	if err := s.feedbackRepo.Create(ctx, feedback); err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}

	return feedback, nil
}

// GetFeedbackForEmployee retrieves all feedback for an employee
func (s *PerformanceService) GetFeedbackForEmployee(ctx context.Context, employee, cycle string) ([]hr.EmployeePerformanceFeedback, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	return s.feedbackRepo.GetFeedbackForEmployee(ctx, employee, cycle)
}

// ========== Skill Map Operations ==========

// CreateOrUpdateSkillMapRequest represents a request to create/update skill map
type CreateOrUpdateSkillMapRequest struct {
	Employee   string `json:"employee" validate:"required"`
	Company    string `json:"company" validate:"required"`
	Department string `json:"department"`
	Designation string `json:"designation"`
}

// CreateOrUpdateSkillMap creates or updates an employee's skill map
func (s *PerformanceService) CreateOrUpdateSkillMap(ctx context.Context, req *CreateOrUpdateSkillMapRequest) (*hr.EmployeeSkillMap, error) {
	if req.Employee == "" {
		return nil, ErrEmployeeRequired
	}

	// Check if skill map exists
	existing, err := s.skillMapRepo.GetByEmployee(ctx, req.Employee)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing skill map: %w", err)
	}

	if existing != nil {
		// Update existing
		existing.Department = req.Department
		existing.Designation = req.Designation
		if err := s.skillMapRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update skill map: %w", err)
		}
		return existing, nil
	}

	// Create new
	skillMap := &hr.EmployeeSkillMap{
		Employee:    req.Employee,
		Company:     req.Company,
		Department:  req.Department,
		Designation: req.Designation,
	}

	if err := s.skillMapRepo.Create(ctx, skillMap); err != nil {
		return nil, fmt.Errorf("failed to create skill map: %w", err)
	}

	return skillMap, nil
}

// GetSkillMapByEmployee retrieves skill map for an employee
func (s *PerformanceService) GetSkillMapByEmployee(ctx context.Context, employee string) (*hr.EmployeeSkillMap, error) {
	if employee == "" {
		return nil, ErrEmployeeRequired
	}

	skillMap, err := s.skillMapRepo.GetByEmployee(ctx, employee)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSkillMapNotFound
		}
		return nil, fmt.Errorf("failed to get skill map: %w", err)
	}

	return skillMap, nil
}

// ========== Appraisal Template & Cycle Operations ==========

// GetActiveAppraisalCycles retrieves active appraisal cycles
func (s *PerformanceService) GetActiveAppraisalCycles(ctx context.Context, company string) ([]hr.AppraisalCycle, error) {
	return s.cycleRepo.GetActiveCycles(ctx, company)
}

// GetActiveAppraisalTemplates retrieves active appraisal templates
func (s *PerformanceService) GetActiveAppraisalTemplates(ctx context.Context, company string) ([]hr.AppraisalTemplate, error) {
	return s.templateRepo.ListActive(ctx, company)
}
