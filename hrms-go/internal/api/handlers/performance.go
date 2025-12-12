package handlers

import (
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
		return response.BadRequest(c, "invalid request body")
	}

	appraisal, err := h.performanceService.CreateAppraisal(c.Context(), &req)
	if err != nil {
		switch err {
		case performance.ErrEmployeeRequired:
			return response.BadRequest(c, "employee is required")
		case performance.ErrTemplateNotFound:
			return response.NotFound(c, "appraisal template not found")
		default:
			return response.InternalServerError(c, "failed to create appraisal")
		}
	}

	return response.Created(c, appraisal, "Appraisal created successfully")
}

// GetAppraisal retrieves an appraisal by ID
// GET /api/method/hrms.hr.doctype.appraisal.appraisal.get_appraisal
func (h *PerformanceHandler) GetAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	appraisal, err := h.performanceService.GetAppraisalByID(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return response.NotFound(c, "appraisal not found")
		}
		return response.InternalServerError(c, "failed to get appraisal")
	}

	return response.Success(c, appraisal, "success")
}

// UpdateAppraisal updates an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.update_appraisal
func (h *PerformanceHandler) UpdateAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	var req performance.UpdateAppraisalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	appraisal, err := h.performanceService.UpdateAppraisal(c.Context(), uint(id), &req)
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return response.NotFound(c, "appraisal not found")
		}
		return response.InternalServerError(c, "failed to update appraisal")
	}

	return response.Success(c, appraisal, "Appraisal updated successfully")
}

// SubmitAppraisal submits an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.submit_appraisal
func (h *PerformanceHandler) SubmitAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	appraisal, err := h.performanceService.SubmitAppraisal(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return response.NotFound(c, "appraisal not found")
		}
		return response.InternalServerError(c, "failed to submit appraisal")
	}

	return response.Success(c, appraisal, "Appraisal submitted successfully")
}

// CompleteAppraisal completes an appraisal
// POST /api/method/hrms.hr.doctype.appraisal.appraisal.complete_appraisal
func (h *PerformanceHandler) CompleteAppraisal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	finalScoreStr := c.Query("final_score")

	if idStr == "" || finalScoreStr == "" {
		return response.BadRequest(c, "id and final_score are required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	finalScore, err := strconv.ParseFloat(finalScoreStr, 64)
	if err != nil {
		return response.BadRequest(c, "invalid final_score")
	}

	appraisal, err := h.performanceService.CompleteAppraisal(c.Context(), uint(id), finalScore)
	if err != nil {
		if err == performance.ErrAppraisalNotFound {
			return response.NotFound(c, "appraisal not found")
		}
		return response.InternalServerError(c, "failed to complete appraisal")
	}

	return response.Success(c, appraisal, "Appraisal completed successfully")
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
		return response.InternalServerError(c, "failed to list appraisals")
	}

	return response.Paginated(c, appraisals, page, pageSize, total)
}

// ========== Goal Operations ==========

// CreateGoal creates a new goal
// POST /api/method/hrms.hr.doctype.goal.goal.create_goal
func (h *PerformanceHandler) CreateGoal(c *fiber.Ctx) error {
	var req performance.CreateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	goal, err := h.performanceService.CreateGoal(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return response.BadRequest(c, "employee is required")
		}
		return response.InternalServerError(c, "failed to create goal")
	}

	return response.Created(c, goal, "Goal created successfully")
}

// GetGoal retrieves a goal by ID
// GET /api/method/hrms.hr.doctype.goal.goal.get_goal
func (h *PerformanceHandler) GetGoal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	goal, err := h.performanceService.GetGoalByID(c.Context(), uint(id))
	if err != nil {
		if err == performance.ErrGoalNotFound {
			return response.NotFound(c, "goal not found")
		}
		return response.InternalServerError(c, "failed to get goal")
	}

	return response.Success(c, goal, "success")
}

// UpdateGoal updates a goal
// POST /api/method/hrms.hr.doctype.goal.goal.update_goal
func (h *PerformanceHandler) UpdateGoal(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return response.BadRequest(c, "id is required")
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid id")
	}

	var req performance.UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	goal, err := h.performanceService.UpdateGoal(c.Context(), uint(id), &req)
	if err != nil {
		switch err {
		case performance.ErrGoalNotFound:
			return response.NotFound(c, "goal not found")
		case performance.ErrInvalidProgress:
			return response.BadRequest(c, "progress must be between 0 and 100")
		default:
			return response.InternalServerError(c, "failed to update goal")
		}
	}

	return response.Success(c, goal, "Goal updated successfully")
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
		return response.InternalServerError(c, "failed to list goals")
	}

	return response.Paginated(c, goals, page, pageSize, total)
}

// GetActiveGoalsForEmployee retrieves active goals for an employee
// GET /api/method/hrms.hr.doctype.goal.goal.get_active_goals
func (h *PerformanceHandler) GetActiveGoalsForEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	goals, err := h.performanceService.GetActiveGoalsForEmployee(c.Context(), employee)
	if err != nil {
		return response.InternalServerError(c, "failed to get active goals")
	}

	return response.Success(c, goals, "success")
}

// ========== 360-Degree Feedback Operations ==========

// CreateFeedback creates a new 360-degree feedback
// POST /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.create_feedback
func (h *PerformanceHandler) CreateFeedback(c *fiber.Ctx) error {
	var req performance.CreateFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	feedback, err := h.performanceService.CreateFeedback(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return response.BadRequest(c, "employee is required")
		}
		return response.InternalServerError(c, "failed to create feedback")
	}

	return response.Created(c, feedback, "Feedback created successfully")
}

// GetFeedbackForEmployee retrieves all feedback for an employee
// GET /api/method/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_feedback
func (h *PerformanceHandler) GetFeedbackForEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	cycle := c.Query("cycle")

	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	feedbacks, err := h.performanceService.GetFeedbackForEmployee(c.Context(), employee, cycle)
	if err != nil {
		return response.InternalServerError(c, "failed to get feedback")
	}

	return response.Success(c, feedbacks, "success")
}

// ========== Skill Map Operations ==========

// CreateOrUpdateSkillMap creates or updates an employee's skill map
// POST /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.create_or_update
func (h *PerformanceHandler) CreateOrUpdateSkillMap(c *fiber.Ctx) error {
	var req performance.CreateOrUpdateSkillMapRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	skillMap, err := h.performanceService.CreateOrUpdateSkillMap(c.Context(), &req)
	if err != nil {
		if err == performance.ErrEmployeeRequired {
			return response.BadRequest(c, "employee is required")
		}
		return response.InternalServerError(c, "failed to create or update skill map")
	}

	return response.Success(c, skillMap, "Skill map saved successfully")
}

// GetSkillMapByEmployee retrieves skill map for an employee
// GET /api/method/hrms.hr.doctype.employee_skill_map.employee_skill_map.get_skill_map
func (h *PerformanceHandler) GetSkillMapByEmployee(c *fiber.Ctx) error {
	employee := c.Query("employee")
	if employee == "" {
		return response.BadRequest(c, "employee is required")
	}

	skillMap, err := h.performanceService.GetSkillMapByEmployee(c.Context(), employee)
	if err != nil {
		if err == performance.ErrSkillMapNotFound {
			return response.NotFound(c, "skill map not found")
		}
		return response.InternalServerError(c, "failed to get skill map")
	}

	return response.Success(c, skillMap, "success")
}

// ========== Template & Cycle Operations ==========

// GetActiveAppraisalCycles retrieves active appraisal cycles
// GET /api/method/hrms.hr.doctype.appraisal_cycle.appraisal_cycle.get_active_cycles
func (h *PerformanceHandler) GetActiveAppraisalCycles(c *fiber.Ctx) error {
	company := c.Query("company")

	cycles, err := h.performanceService.GetActiveAppraisalCycles(c.Context(), company)
	if err != nil {
		return response.InternalServerError(c, "failed to get active cycles")
	}

	return response.Success(c, cycles, "success")
}

// GetActiveAppraisalTemplates retrieves active appraisal templates
// GET /api/method/hrms.hr.doctype.appraisal_template.appraisal_template.get_active_templates
func (h *PerformanceHandler) GetActiveAppraisalTemplates(c *fiber.Ctx) error {
	company := c.Query("company")

	templates, err := h.performanceService.GetActiveAppraisalTemplates(c.Context(), company)
	if err != nil {
		return response.InternalServerError(c, "failed to get active templates")
	}

	return response.Success(c, templates, "success")
}
