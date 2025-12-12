package hr

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// Appraisal represents employee performance appraisal
type Appraisal struct {
	base.BaseModel

	// Naming
	NamingSeries        string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`

	// Appraisal Period
	AppraisalCycle      string     `gorm:"size:255" json:"appraisal_cycle"`
	StartDate           time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time  `gorm:"type:date;not null" json:"end_date"`

	// Template
	AppraisalTemplate   string     `gorm:"size:255" json:"appraisal_template"`

	// Goals & KRAs
	Goals               []AppraisalGoal `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"goals"`

	// Ratings
	TotalScore          float64    `gorm:"type:decimal(5,2)" json:"total_score"`

	// Feedback
	SelfAppraisal       string     `gorm:"type:text" json:"self_appraisal"`
	ManagerFeedback     string     `gorm:"type:text" json:"remarks"`

	// Status
	Status              string     `gorm:"size:50;default:'Draft'" json:"status"` // Draft, Submitted, Completed, Cancelled
	FinalScore          float64    `gorm:"type:decimal(5,2)" json:"final_score"`
}

// TableName specifies the table name
func (Appraisal) TableName() string {
	return "appraisals"
}

// AppraisalGoal represents goals/KRAs in an appraisal
type AppraisalGoal struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	KRA                 string  `gorm:"size:255;not null" json:"kra"` // Key Result Area
	KRATitle            string  `gorm:"size:255" json:"kra_title"`
	Weightage           float64 `gorm:"type:decimal(5,2)" json:"per_weightage"`

	// Scoring
	Score               float64 `gorm:"type:decimal(5,2)" json:"score"`
	ScoreEarnedPoints   float64 `gorm:"type:decimal(5,2)" json:"score_earned"`
}

// TableName specifies the table name
func (AppraisalGoal) TableName() string {
	return "appraisal_goals"
}

// AppraisalTemplate represents appraisal template
type AppraisalTemplate struct {
	base.BaseModel
	TemplateName        string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description         string `gorm:"type:text" json:"description"`
	Company             string `gorm:"size:255" json:"company"`

	// KRA Templates
	Goals               []AppraisalTemplateGoal `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"goals"`

	// Rating Scale
	RatingScale         []AppraisalTemplateRating `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"rating_scale"`
}

// TableName specifies the table name
func (AppraisalTemplate) TableName() string {
	return "appraisal_templates"
}

// AppraisalTemplateGoal represents template goals
type AppraisalTemplateGoal struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	KRA                 string  `gorm:"size:255;not null" json:"kra"`
	KRATitle            string  `gorm:"size:255" json:"kra_title"`
	Weightage           float64 `gorm:"type:decimal(5,2)" json:"per_weightage"`
}

// TableName specifies the table name
func (AppraisalTemplateGoal) TableName() string {
	return "appraisal_template_goals"
}

// AppraisalTemplateRating represents rating scale
type AppraisalTemplateRating struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	RatingCriteria      string  `gorm:"size:255;not null" json:"rating_criteria"`
	ScoreFrom           float64 `gorm:"type:decimal(5,2)" json:"score_from"`
	ScoreTo             float64 `gorm:"type:decimal(5,2)" json:"score_to"`
}

// TableName specifies the table name
func (AppraisalTemplateRating) TableName() string {
	return "appraisal_template_ratings"
}

// AppraisalCycle represents appraisal period
type AppraisalCycle struct {
	base.BaseModel
	CycleName           string    `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Company             string    `gorm:"size:255;not null" json:"company"`
	StartDate           time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate             time.Time `gorm:"type:date;not null" json:"end_date"`
	AppraisalTemplate   string    `gorm:"size:255" json:"appraisal_template"`
	Status              string    `gorm:"size:50;default:'Active'" json:"status"` // Active, Completed
}

// TableName specifies the table name
func (AppraisalCycle) TableName() string {
	return "appraisal_cycles"
}

// Goal represents employee goals/objectives
type Goal struct {
	base.BaseModel

	// Goal Details
	GoalName            string     `gorm:"size:255;not null" json:"subject"`
	Description         string     `gorm:"type:text" json:"description"`

	// Assignment
	Employee            string     `gorm:"size:255;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	AssignedBy          string     `gorm:"size:255" json:"assigned_by"`
	Company             string     `gorm:"size:255;not null" json:"organization"`

	// Timeline
	StartDate           time.Time  `gorm:"type:date" json:"from_time"`
	EndDate             time.Time  `gorm:"type:date" json:"to_time"`

	// Progress
	Progress            float64    `gorm:"type:decimal(5,2);default:0" json:"progress"`
	Status              string     `gorm:"size:50;default:'Pending'" json:"status"` // Pending, Achieved, Partially Achieved, Failed, Archived

	// KPI/Alignment
	KPIGoal             bool       `gorm:"default:false" json:"is_kpi_goal"`
	AlignWithGoal       string     `gorm:"size:255" json:"align_with_goal"`
}

// TableName specifies the table name
func (Goal) TableName() string {
	return "goals"
}

// EmployeePerformanceFeedback represents 360-degree feedback
type EmployeePerformanceFeedback struct {
	base.BaseModel

	// Employee Being Reviewed
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`

	// Reviewer
	Reviewer            string     `gorm:"size:255;not null" json:"reviewer"`
	ReviewerName        string     `gorm:"size:255" json:"reviewer_name"`

	// Feedback Period
	FeedbackDate        time.Time  `gorm:"type:date;not null" json:"feedback_date"`
	AppraisalCycle      string     `gorm:"size:255" json:"appraisal_cycle"`

	// Template
	FeedbackTemplate    string     `gorm:"size:255" json:"feedback_template"`

	// Feedback Criteria
	FeedbackCriteria    []EmployeePerformanceFeedbackCriteria `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"feedback_criteria"`

	// Overall
	OverallFeedback     string     `gorm:"type:text" json:"overall_feedback"`
	OverallRating       float64    `gorm:"type:decimal(5,2)" json:"total_score"`
}

// TableName specifies the table name
func (EmployeePerformanceFeedback) TableName() string {
	return "employee_performance_feedbacks"
}

// EmployeePerformanceFeedbackCriteria represents feedback criteria
type EmployeePerformanceFeedbackCriteria struct {
	gorm.Model
	ParentID            uint    `gorm:"index" json:"parent_id"`
	Criteria            string  `gorm:"size:255;not null" json:"criteria"`
	Weightage           float64 `gorm:"type:decimal(5,2)" json:"per_weightage"`
	Rating              float64 `gorm:"type:decimal(5,2)" json:"rating"`
	Feedback            string  `gorm:"type:text" json:"feedback"`
}

// TableName specifies the table name
func (EmployeePerformanceFeedbackCriteria) TableName() string {
	return "employee_performance_feedback_criteria"
}

// EmployeeSkillMap represents employee skill assessment
type EmployeeSkillMap struct {
	base.BaseModel

	// Employee
	Employee            string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName        string     `gorm:"size:255" json:"employee_name"`
	Company             string     `gorm:"size:255;not null" json:"company"`
	Department          string     `gorm:"size:255" json:"department"`
	Designation         string     `gorm:"size:255" json:"designation"`

	// Skills
	Skills              []EmployeeSkillMapSkill `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"employee_skills"`

	// Trainings
	Trainings           []EmployeeSkillMapTraining `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"trainings"`
}

// TableName specifies the table name
func (EmployeeSkillMap) TableName() string {
	return "employee_skill_maps"
}

// EmployeeSkillMapSkill represents skills in skill map
type EmployeeSkillMapSkill struct {
	gorm.Model
	ParentID            uint   `gorm:"index" json:"parent_id"`
	Skill               string `gorm:"size:255;not null" json:"skill"`
	Proficiency         string `gorm:"size:50" json:"proficiency"` // Beginner, Intermediate, Advanced, Expert
	EvaluationDate      *time.Time `gorm:"type:date" json:"evaluation_date"`
}

// TableName specifies the table name
func (EmployeeSkillMapSkill) TableName() string {
	return "employee_skill_map_skills"
}

// EmployeeSkillMapTraining represents trainings in skill map
type EmployeeSkillMapTraining struct {
	gorm.Model
	ParentID            uint   `gorm:"index" json:"parent_id"`
	TrainingName        string `gorm:"size:255;not null" json:"training"`
	Status              string `gorm:"size:50" json:"status"` // Completed, In Progress, Planned
}

// TableName specifies the table name
func (EmployeeSkillMapTraining) TableName() string {
	return "employee_skill_map_trainings"
}
