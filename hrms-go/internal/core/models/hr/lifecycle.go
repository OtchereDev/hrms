package hr

import (
	"time"

	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

// EmployeeOnboarding represents employee onboarding process
type EmployeeOnboarding struct {
	base.BaseModel

	// Naming
	NamingSeries     string     `gorm:"size:50" json:"naming_series"`

	// Job Details
	JobApplicant     string     `gorm:"size:255;index" json:"job_applicant"`
	JobOffer         string     `gorm:"size:255" json:"job_offer"`

	// Employee Info (before employee is created)
	FirstName        string     `gorm:"size:150" json:"first_name"`
	MiddleName       string     `gorm:"size:150" json:"middle_name"`
	LastName         string     `gorm:"size:150" json:"last_name"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Gender           string     `gorm:"size:20" json:"gender"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	DateOfJoining    time.Time  `gorm:"type:date;not null" json:"date_of_joining"`

	// Organization
	Company          string     `gorm:"size:255;not null" json:"company"`
	Department       string     `gorm:"size:255" json:"department"`
	Designation      string     `gorm:"size:255;not null" json:"designation"`
	Branch           string     `gorm:"size:255" json:"branch"`
	Grade            string     `gorm:"size:255" json:"grade"`
	EmploymentType   string     `gorm:"size:255" json:"employment_type"`

	// Created Employee Reference
	Employee         string     `gorm:"size:255;index" json:"employee"`

	// Onboarding Template
	EmployeeOnboardingTemplate string `gorm:"size:255" json:"employee_onboarding_template"`

	// Activities/Tasks
	Activities       []EmployeeBoardingActivity `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"activities"`

	// Status
	BoardingStatus   string     `gorm:"size:50;default:'Pending'" json:"boarding_status"` // Pending, In Progress, Completed

	// Dates
	ExpectedCompletionDate *time.Time `gorm:"type:date" json:"expected_completion_date"`
	ActualCompletionDate   *time.Time `gorm:"type:date" json:"actual_completion_date"`

	// Project
	Project          string     `gorm:"size:255" json:"project"`
}

// TableName specifies the table name
func (EmployeeOnboarding) TableName() string {
	return "employee_onboardings"
}

// EmployeeSeparation represents employee exit/separation process
type EmployeeSeparation struct {
	base.BaseModel

	// Naming
	NamingSeries     string     `gorm:"size:50" json:"naming_series"`

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`
	Department       string     `gorm:"size:255" json:"department"`
	Designation      string     `gorm:"size:255" json:"designation"`

	// Separation Details
	ResignationLetterDate *time.Time `gorm:"type:date" json:"resignation_letter_date"`
	ExitInterviewDate     *time.Time `gorm:"type:date" json:"exit_interview_date"`
	RelievingDate         *time.Time `gorm:"type:date;not null" json:"relieving_date"`
	ExitInterviewHeld     bool       `gorm:"default:false" json:"exit_interview_held"`
	NewWorkplace          string     `gorm:"type:text" json:"new_workplace"`

	// Reason
	ReasonForLeaving      string     `gorm:"type:text;not null" json:"reason_for_leaving"`
	FeedbackFromEmployee  string     `gorm:"type:text" json:"leave_feedback"`
	FeedbackFromHR        string     `gorm:"type:text" json:"hr_remarks"`

	// Template
	EmployeeSeparationTemplate string `gorm:"size:255" json:"employee_separation_template"`

	// Activities/Tasks
	Activities       []EmployeeBoardingActivity `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"activities"`

	// Status
	BoardingStatus   string     `gorm:"size:50;default:'Pending'" json:"boarding_status"` // Pending, In Progress, Completed

	// References
	ReferenceDocument string    `gorm:"size:255" json:"reference_document_name"`
}

// TableName specifies the table name
func (EmployeeSeparation) TableName() string {
	return "employee_separations"
}

// EmployeeBoardingActivity represents tasks/activities for onboarding/separation
type EmployeeBoardingActivity struct {
	gorm.Model
	ParentID         uint       `gorm:"index" json:"parent_id"`
	ParentType       string     `gorm:"size:50" json:"parent_type"` // "Onboarding" or "Separation"

	ActivityName     string     `gorm:"size:255;not null" json:"activity_name"`
	Role             string     `gorm:"size:255" json:"role"`
	RequiredForCompletion bool  `gorm:"default:false" json:"required_for_completion"`
	Description      string     `gorm:"type:text" json:"description"`

	// Task
	Task             string     `gorm:"size:255" json:"task"`

	// Completion
	User             string     `gorm:"size:255" json:"user"` // Assigned to
	Completed        bool       `gorm:"default:false" json:"completed"`
	CompletedBy      string     `gorm:"size:255" json:"completed_by"`
	CompletedOn      *time.Time `json:"completed_on"`
}

// TableName specifies the table name
func (EmployeeBoardingActivity) TableName() string {
	return "employee_boarding_activities"
}

// EmployeeTransfer represents internal transfer of employee
type EmployeeTransfer struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Transfer Date
	TransferDate     time.Time  `gorm:"type:date;not null" json:"transfer_date"`

	// Current Details
	CurrentCompany   string     `gorm:"size:255" json:"current_company"`
	CurrentDepartment string    `gorm:"size:255" json:"current_department"`
	CurrentDesignation string   `gorm:"size:255" json:"current_designation"`
	CurrentBranch    string     `gorm:"size:255" json:"current_branch"`
	CurrentGrade     string     `gorm:"size:255" json:"current_grade"`
	CurrentShiftType string     `gorm:"size:255" json:"current_shift"`

	// New Details
	NewCompany       string     `gorm:"size:255" json:"new_company"`
	NewDepartment    string     `gorm:"size:255" json:"new_department"`
	NewDesignation   string     `gorm:"size:255" json:"new_designation"`
	NewBranch        string     `gorm:"size:255" json:"new_branch"`
	NewGrade         string     `gorm:"size:255" json:"new_grade"`
	NewShiftType     string     `gorm:"size:255" json:"new_shift"`

	// Transfer Details
	TransferReason   string     `gorm:"type:text" json:"reason_for_transfer"`
	CreateNewEmployee bool      `gorm:"default:false" json:"create_new_employee_id"`
	NewEmployeeID    string     `gorm:"size:255" json:"new_employee_id"`

	// References
	ReferenceDocument string    `gorm:"size:255" json:"reference_document"`
}

// TableName specifies the table name
func (EmployeeTransfer) TableName() string {
	return "employee_transfers"
}

// EmployeePromotion represents employee promotion
type EmployeePromotion struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Promotion Date
	PromotionDate    time.Time  `gorm:"type:date;not null" json:"promotion_date"`

	// Current Details
	CurrentDesignation string   `gorm:"size:255" json:"current_designation"`
	CurrentGrade       string   `gorm:"size:255" json:"current_grade"`

	// New Details
	NewDesignation   string     `gorm:"size:255;not null" json:"new_designation"`
	NewGrade         string     `gorm:"size:255" json:"new_grade"`

	// Promotion Details
	PromotionDetails string     `gorm:"type:text" json:"promotion_details"`
}

// TableName specifies the table name
func (EmployeePromotion) TableName() string {
	return "employee_promotions"
}

// EmployeeOnboardingTemplate represents onboarding template
type EmployeeOnboardingTemplate struct {
	base.BaseModel
	TemplateName     string                     `gorm:"size:255;not null;uniqueIndex" json:"template_name"`
	Department       string                     `gorm:"size:255" json:"department"`
	Designation      string                     `gorm:"size:255" json:"designation"`
	Description      string                     `gorm:"type:text" json:"description"`
	Activities       []EmployeeBoardingTemplateActivity `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"activities"`
}

// TableName specifies the table name
func (EmployeeOnboardingTemplate) TableName() string {
	return "employee_onboarding_templates"
}

// EmployeeSeparationTemplate represents separation template
type EmployeeSeparationTemplate struct {
	base.BaseModel
	TemplateName     string                     `gorm:"size:255;not null;uniqueIndex" json:"template_name"`
	Department       string                     `gorm:"size:255" json:"department"`
	Designation      string                     `gorm:"size:255" json:"designation"`
	Description      string                     `gorm:"type:text" json:"description"`
	Activities       []EmployeeBoardingTemplateActivity `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"activities"`
}

// TableName specifies the table name
func (EmployeeSeparationTemplate) TableName() string {
	return "employee_separation_templates"
}

// EmployeeBoardingTemplateActivity represents template activities
type EmployeeBoardingTemplateActivity struct {
	gorm.Model
	ParentID         uint   `gorm:"index" json:"parent_id"`
	ParentType       string `gorm:"size:50" json:"parent_type"` // "Onboarding Template" or "Separation Template"

	ActivityName     string `gorm:"size:255;not null" json:"activity_name"`
	Role             string `gorm:"size:255" json:"role"`
	RequiredForCompletion bool `gorm:"default:false" json:"required_for_completion"`
	Description      string `gorm:"type:text" json:"description"`
}

// TableName specifies the table name
func (EmployeeBoardingTemplateActivity) TableName() string {
	return "employee_boarding_template_activities"
}

// ExitInterview represents exit interview details
type ExitInterview struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`
	Department       string     `gorm:"size:255" json:"department"`
	Designation      string     `gorm:"size:255" json:"designation"`

	// Interview
	InterviewDate    time.Time  `gorm:"type:date;not null" json:"date"`
	InterviewedBy    string     `gorm:"size:255" json:"interviewed_by"`

	// Employee Feedback
	HeldOn           *time.Time `json:"held_on"`
	Status           string     `gorm:"size:50" json:"status"`
	Comments         string     `gorm:"type:text" json:"comments"`

	// Separation Reference
	EmployeeSeparation string   `gorm:"size:255" json:"employee_separation"`
}

// TableName specifies the table name
func (ExitInterview) TableName() string {
	return "exit_interviews"
}

// EmployeeGrievance represents employee grievance/complaint
type EmployeeGrievance struct {
	base.BaseModel

	// Employee
	Employee         string     `gorm:"size:255;not null;index" json:"employee"`
	EmployeeName     string     `gorm:"size:255" json:"employee_name"`
	Company          string     `gorm:"size:255;not null" json:"company"`

	// Grievance Details
	GrievanceAgainst string     `gorm:"size:255" json:"grievance_against"`
	GrievanceAgainstParty string `gorm:"size:255" json:"grievance_against_party"` // Employee, Department, etc.
	Subject          string     `gorm:"size:255;not null" json:"subject"`
	GrievanceDate    time.Time  `gorm:"type:date;not null" json:"raised_on"`
	Description      string     `gorm:"type:text;not null" json:"description"`

	// Actions
	CauseOfGrievance string     `gorm:"type:text" json:"cause_of_grievance"`
	ActionsTaken     string     `gorm:"type:text" json:"actions_taken"`

	// Resolution
	ResolutionDate   *time.Time `gorm:"type:date" json:"resolution_date"`
	ResolutionDetails string    `gorm:"type:text" json:"resolution_details"`

	// Status
	Status           string     `gorm:"size:50;default:'Open'" json:"status"` // Open, In Progress, Resolved, Closed

	// Assigned To
	AssignedTo       string     `gorm:"size:255" json:"assigned_to"`
}

// TableName specifies the table name
func (EmployeeGrievance) TableName() string {
	return "employee_grievances"
}
