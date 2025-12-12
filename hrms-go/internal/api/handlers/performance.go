package handlers

import (
	frappeResponse "github.com/OtchereDev/hrms-go/internal/core/frappe"
	"strconv"
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/repositories"
	"github.com/OtchereDev/hrms-go/internal/core/services/performance"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PerformanceHandler handles performance management requests
type PerformanceHandler struct {
	performanceService *performance.PerformanceService
}

// NewPerformanceHandler creates a new performance handler
func NewPerformanceHandler(db *gorm.DB) *PerformanceHandler {
	return &PerformanceHandler{
		performanceService: performance.NewPerformanceService(db),
	}
}

// ========== Appraisal Operations ==========

// CreateAppraisal creates a new appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.create_appraisal
func (h *PerformanceHandler) CreateAppraisal(c *fiber.Ctx) error {
	var req performance.CreateAppraisalRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	appraisal, err := h.performanceService.CreateAppraisal(c.Context(), &req)
	if err != nil {
		switch err {
		case performance.ErrEmployeeRequired:
			return frappeResponse.SendBadRequest(c, "employee is required")
		case performance.ErrTemplateNotFound:
			return frappeResponse.SendNotFound(c, "appraisal template not found")
		default:
			return frappeResponse.SendInternalError(c, "failed to create appraisal")
		}
	}

	return response.Created(c, appraisal, "Appraisal created successfully")
}

// GetAppraisal retrieves an appraisal by ID
// GET /api/method/hrms.hr.doctype.appraisal.appraisal.get_appraisal
func (h *PerformanceHandler) GetAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	appraisal, err := h.performanceService.GetAppraisalByID(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return frappeResponse.SendNotFound(c, "appraisal not found")
		}
		return frappeResponse.SendInternalError(c, "failed to get appraisal")
	}

	return frappeResponse.SendSuccess(c, appraisal)
}

// UpdateAppraisal updates an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.update_appraisal
func (h *PerformanceHandler) UpdateAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	var req performance.UpdateAppraisalRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	appraisal, err := h.performanceService.UpdateAppraisal(c.Context(), uint(id), &req)
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return frappeResponse.SendNotFound(c, "appraisal not found")
		}
		return frappeResponse.SendInternalError(c, "failed to update appraisal")
	}

	return frappeResponse.SendSuccess(c, appraisal)
}

// SubmitAppraisal submits an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.submit_appraisal
func (h *PerformanceHandler) SubmitAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	appraisal, err := h.performanceService.SubmitAppraisal(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return frappeResponse.SendNotFound(c, "appraisal not found")
		}
		return frappeResponse.SendInternalError(c, "failed to submit appraisal")
	}

	return frappeResponse.SendSuccess(c, appraisal)
}

// CompleteAppraisal completes an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.complete_appraisal
func (h *PerformanceHandler) CompleteAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	finalScoreStr := c.Query("final_score")

	if idStr == "" || finalScoreStr == "" {
		return frappeResponse.SendBadRequest(c, "id and final_score are required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	finalScore, err := strconv.ParseFloat(finalScoreStr, 64)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid final_score")
	}

	appraisal, err := h.performanceService.CompleteAppraisal(c.Context(), uint(id), finalScore)
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return frappeResponse.SendNotFound(c, "appraisal not found")
		}
		return frappeResponse.SendInternalError(c, "failed to complete appraisal")
	}

	return frappeResponse.SendSuccess(c, appraisal)
}

// ListAppraisals retrieves appraisals with filters
// GET /api/method/hrms.hr.doctype.appraisal.appraisal.list_appraisals
func (h *PerformanceHandler) ListAppraisals(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.AppraisalFilters{
		Employee:       c.Query("employee"),
		Company:        c.Query("company"),
		Department:     c.Query("department"),
		Status:         c.Query("status"),
		AppraisalCycle: c.Query("appraisal_cycle"),
	}

	// Parse dates if provided
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filters.EndDate = endDate
		}
	}

	appraisals, total, err := h.performanceService.ListAppraisals(c.Context(), filters, page, pageSize)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to list appraisals")
	}

	return response.Paginated(c, appraisals, page, pageSize, total)
}

// ========== Goal Operations ==========

// CreateGoal creates a new goal
// POST /api/method/hrms.hr.doctype.goal.goal.create_goal
func (h *PerformanceHandler) CreateGoal(c *fiber.Ctx) error {
	var req performance.CreateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	goal, err := h.performanceService.CreateGoal(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "employee is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create goal")
	}

	return response.Created(c, goal, "Goal created successfully")
}

// GetGoal retrieves a goal by ID
// GET /api/method/hrms.hr.doctype.goal.goal.get_goal
func (h *PerformanceHandler) GetGoal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	goal, err := h.performanceService.GetGoalByID(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrGoalNotFound {
			return frappeResponse.SendNotFound(c, "goal not found")
		}
		return frappeResponse.SendInternalError(c, "failed to get goal")
	}

	return frappeResponse.SendSuccess(c, goal)
}

// UpdateGoal updates a goal
// POST /api/method/hrms.hr.doctype.goal.goal.update_goal
func (h *PerformanceHandler) UpdateGoal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return frappeResponse.SendBadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return frappeResponse.SendBadRequest(c, "invalid id")
	}

	var req performance.UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	goal, err := h.performanceService.UpdateGoal(c.Context(), uint(id), &req)
	if err != nil {
		switch err {
		case performance.ErrGoalNotFound:
			return frappeResponse.SendNotFound(c, "goal not found")
		case performance.ErrInvalidProgress:
			return frappeResponse.SendBadRequest(c, "progress must be between 0 and 100")
		default:
			return frappeResponse.SendInternalError(c, "failed to update goal")
		}
	}

	return frappeResponse.SendSuccess(c, goal)
}

// ListGoals retrieves goals with filters
// GET /api/method/hrms.hr.doctype.goal.goal.list_goals
func (h *PerformanceHandler) ListGoals(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	filters := repositories.GoalFilters{
		Employee:   c.Query("employee"),
		Company:    c.Query("company"),
		Status:     c.Query("status"),
		AssignedBy: c.Query("assigned_by"),
	}

	goals, total, err := h.performanceService.ListGoals(c.Context(), filters, page, pageSize)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to list goals")
	}

	return response.Paginated(c, goals, page, pageSize, total)
}

// GetActiveGoalsForEmployee retrieves active goals for an employee
// GET /api/method/hrms.hr.doctype.goal.goal.get_active_goals
func (h *PerformanceHandler) GetActiveGoalsForEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return frappeResponse.SendBadRequest(c, "employee is required")
	}

	goals, err := h.performanceService.GetActiveGoalsForEmployee(c.Context(), employee)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get active goals")
	}

	return frappeResponse.SendSuccess(c, goals)
}

// ========== 360-Degree Feedback Operations ==========

// CreateFeedback creates a new 360-degree feedback
// POST /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.create_feedback
func (h *PerformanceHandler) CreateFeedback(c *fiber.Ctx) error {
	var req performance.CreateFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	feedback, err := h.performanceService.CreateFeedback(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "employee is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create feedback")
	}

	return response.Created(c, feedback, "Feedback created successfully")
}

// GetFeedbackForEmployee retrieves all feedback for an employee
// GET /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_feedback
func (h *PerformanceHandler) GetFeedbackForEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	cycle := c.Query("cycle")

	if employee == "" {
		return frappeResponse.SendBadRequest(c, "employee is required")
	}

	feedbacks, err := h.performanceService.GetFeedbackForEmployee(c.Context(), employee, cycle)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get feedback")
	}

	return frappeResponse.SendSuccess(c, feedbacks)
}

// ========== Skill Map Operations ==========

// CreateOrUpdateSkillMap creates or updates an employee's skill map
// POST /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.create_or_update
func (h *PerformanceHandler) CreateOrUpdateSkillMap(c *fiber.Ctx) error {
	var req performance.CreateOrUpdateSkillMapRequest
	if err := c.BodyParser(&req); err != nil {
		return frappeResponse.SendBadRequest(c, "invalid request body")
	}

	skillMap, err := h.performanceService.CreateOrUpdateSkillMap(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return frappeResponse.SendBadRequest(c, "employee is required")
		}
		return frappeResponse.SendInternalError(c, "failed to create or update skill map")
	}

	return frappeResponse.SendSuccess(c, skillMap)
}

// GetSkillMapByEmployee retrieves skill map for an employee
// GET /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.get_skill_map
func (h *PerformanceHandler) GetSkillMapByEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return frappeResponse.SendBadRequest(c, "employee is required")
	}

	skillMap, err := h.performanceService.GetSkillMapByEmployee(c.Context(), employee)
	if err != nil {
		if err == performance.ErrSkillMapNotFound {
			return frappeResponse.SendNotFound(c, "skill map not found")
		}
		return frappeResponse.SendInternalError(c, "failed to get skill map")
	}

	return frappeResponse.SendSuccess(c, skillMap)
}

// ========== Template & Cycle Operations ==========

// GetActiveAppraisalCycles retrieves active appraisal cycles
// GET /api/method/hrms.hr.doctype.appraisal_cycle.appraisal_cycle.get_active_cycles
func (h *PerformanceHandler) GetActiveAppraisalCycles(c *fiber.Ctx) error {
	company := c.Query("company")

	cycles, err := h.performanceService.GetActiveAppraisalCycles(c.Context(), company)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get active cycles")
	}

	return frappeResponse.SendSuccess(c, cycles)
}

// GetActiveAppraisalTemplates retrieves active appraisal templates
// GET /api/method/hrms.hr.doctype.appraisal_template.appraisal_template.get_active_templates
func (h *PerformanceHandler) GetActiveAppraisalTemplates(c *fiber.Ctx) error {
	company := c.Query("company")

	templates, err := h.performanceService.GetActiveAppraisalTemplates(c.Context(), company)
	if err != nil {
		return frappeResponse.SendInternalError(c, "failed to get active templates")
	}

	return frappeResponse.SendSuccess(c, templates)
}
