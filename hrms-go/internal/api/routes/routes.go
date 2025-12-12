package routes

import (
	"github.com/OtchereDev/hrms-go/internal/api/handlers"
	"github.com/OtchereDev/hrms-go/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	employeeHandler *handlers.EmployeeHandler,
	attendanceHandler *handlers.AttendanceHandler,
	leaveHandler *handlers.LeaveHandler,
	payrollHandler *handlers.PayrollHandler,
	performanceHandler *handlers.PerformanceHandler,
	authMiddleware *middleware.AuthMiddleware,
	permMiddleware *middleware.PermissionMiddleware,
) {
	// API version prefix
	api := app.Group("/api")

	// Method routes (Frappe-compatible API pattern)
	method := api.Group("/method")

	// Public authentication routes
	method.Post("/login", authHandler.Login)
	method.Post("/logout", authMiddleware.Authenticate, authHandler.Logout)
	method.Post("/refresh_token", authHandler.RefreshToken)

	// Protected routes - require authentication
	authenticated := method.Use(authMiddleware.Authenticate)

	// User info routes
	authenticated.Get("/hrms.api.get_current_user_info", authHandler.GetCurrentUserInfo)
	authenticated.Get("/hrms.api.get_current_employee_info", employeeHandler.GetCurrentEmployeeInfo)

	// Employee routes
	employeeRoutes := authenticated.Group("")

	// List employees (accessible to HR User and above)
	employeeRoutes.Get("/hrms.api.get_employees",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.ListEmployees,
	)

	employeeRoutes.Get("/hrms.api.get_active_employees",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.GetActiveEmployees,
	)

	employeeRoutes.Get("/hrms.api.get_employees_by_department",
		permMiddleware.RequireRole("HR User", "HR Manager", "System Manager"),
		employeeHandler.GetEmployeesByDepartment,
	)

	// Get employee (accessible to all authenticated users)
	employeeRoutes.Get("/hrms.hr.doctype.employee.employee.get_employee",
		employeeHandler.GetEmployee,
	)

	employeeRoutes.Get("/hrms.hr.doctype.employee.employee.get_employee_details",
		employeeHandler.GetEmployeeDetails,
	)

	// Create/Update/Delete employee (HR Manager and above only)
	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.create_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.CreateEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.update_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.UpdateEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.delete_employee",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.DeleteEmployee,
	)

	employeeRoutes.Post("/hrms.hr.doctype.employee.employee.update_employee_status",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		employeeHandler.UpdateEmployeeStatus,
	)

	// ========== Attendance Routes ==========
	attendanceRoutes := authenticated.Group("")

	// Mark attendance (HR Manager and above)
	attendanceRoutes.Post("/hrms.hr.doctype.attendance.attendance.mark_attendance",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		attendanceHandler.MarkAttendance,
	)

	// Get attendance (accessible to all authenticated users)
	attendanceRoutes.Get("/hrms.hr.doctype.attendance.attendance.get_attendance",
		attendanceHandler.GetAttendance,
	)

	// List attendance
	attendanceRoutes.Get("/hrms.hr.doctype.attendance.attendance.list_attendance",
		attendanceHandler.ListAttendance,
	)

	// Get monthly attendance
	attendanceRoutes.Get("/hrms.hr.doctype.attendance.attendance.get_monthly_attendance",
		attendanceHandler.GetMonthlyAttendance,
	)

	// Check-in/Check-out (accessible to all employees)
	attendanceRoutes.Post("/hrms.hr.doctype.employee_checkin.employee_checkin.checkin",
		attendanceHandler.Checkin,
	)

	// Get today's check-ins
	attendanceRoutes.Get("/hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins",
		attendanceHandler.GetTodayCheckins,
	)

	// Attendance requests
	attendanceRoutes.Post("/hrms.hr.doctype.attendance_request.attendance_request.create_request",
		attendanceHandler.CreateAttendanceRequest,
	)

	attendanceRoutes.Post("/hrms.hr.doctype.attendance_request.attendance_request.approve_request",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		attendanceHandler.ApproveAttendanceRequest,
	)

	attendanceRoutes.Post("/hrms.hr.doctype.attendance_request.attendance_request.reject_request",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		attendanceHandler.RejectAttendanceRequest,
	)

	// ========== Leave Routes ==========
	leaveRoutes := authenticated.Group("")

	// Apply for leave (accessible to all employees)
	leaveRoutes.Post("/hrms.hr.doctype.leave_application.leave_application.apply_leave",
		leaveHandler.ApplyLeave,
	)

	// Get leave application
	leaveRoutes.Get("/hrms.hr.doctype.leave_application.leave_application.get_leave_application",
		leaveHandler.GetLeaveApplication,
	)

	// List leave applications
	leaveRoutes.Get("/hrms.hr.doctype.leave_application.leave_application.list_leave_applications",
		leaveHandler.ListLeaveApplications,
	)

	// Approve/Reject leave (HR Manager and above)
	leaveRoutes.Post("/hrms.hr.doctype.leave_application.leave_application.approve_leave",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		leaveHandler.ApproveLeaveApplication,
	)

	leaveRoutes.Post("/hrms.hr.doctype.leave_application.leave_application.reject_leave",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		leaveHandler.RejectLeaveApplication,
	)

	// Cancel leave (employee can cancel their own)
	leaveRoutes.Post("/hrms.hr.doctype.leave_application.leave_application.cancel_leave",
		leaveHandler.CancelLeaveApplication,
	)

	// Leave allocation (HR Manager and above)
	leaveRoutes.Post("/hrms.hr.doctype.leave_allocation.leave_allocation.allocate_leave",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		leaveHandler.AllocateLeave,
	)

	// Get leave balance
	leaveRoutes.Get("/hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance",
		leaveHandler.GetLeaveBalance,
	)

	// Get active leave types
	leaveRoutes.Get("/hrms.hr.doctype.leave_type.leave_type.get_active_leave_types",
		leaveHandler.GetActiveLeaveTypes,
	)

	// Leave encashment
	leaveRoutes.Post("/hrms.hr.doctype.leave_encashment.leave_encashment.create_encashment",
		leaveHandler.CreateLeaveEncashment,
	)

	// ========== Payroll Routes ==========
	payrollRoutes := authenticated.Group("")

	// Salary Slip operations (HR Manager and above for generation)
	payrollRoutes.Post("/hrms.payroll.doctype.salary_slip.salary_slip.generate_salary_slip",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.GenerateSalarySlip,
	)

	payrollRoutes.Get("/hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip",
		payrollHandler.GetSalarySlip,
	)

	payrollRoutes.Post("/hrms.payroll.doctype.salary_slip.salary_slip.submit_salary_slip",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.SubmitSalarySlip,
	)

	payrollRoutes.Get("/hrms.payroll.doctype.salary_slip.salary_slip.list_salary_slips",
		payrollHandler.ListSalarySlips,
	)

	// Salary structure assignment (HR Manager and above)
	payrollRoutes.Post("/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.assign_structure",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.AssignSalaryStructure,
	)

	payrollRoutes.Get("/hrms.payroll.doctype.salary_structure_assignment.salary_structure_assignment.get_active_assignment",
		payrollHandler.GetActiveSalaryAssignment,
	)

	// Loan operations
	payrollRoutes.Post("/hrms.payroll.doctype.loan.loan.create_loan",
		payrollHandler.CreateLoan,
	)

	payrollRoutes.Post("/hrms.payroll.doctype.loan.loan.approve_loan",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.ApproveLoan,
	)

	payrollRoutes.Get("/hrms.payroll.doctype.loan.loan.list_loans",
		payrollHandler.ListLoans,
	)

	// Employee advance operations
	payrollRoutes.Post("/hrms.payroll.doctype.employee_advance.employee_advance.create_advance",
		payrollHandler.CreateEmployeeAdvance,
	)

	payrollRoutes.Post("/hrms.payroll.doctype.employee_advance.employee_advance.approve_advance",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.ApproveAdvance,
	)

	// Expense claim operations
	payrollRoutes.Post("/hrms.payroll.doctype.expense_claim.expense_claim.create_expense_claim",
		payrollHandler.CreateExpenseClaim,
	)

	payrollRoutes.Post("/hrms.payroll.doctype.expense_claim.expense_claim.approve_claim",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		payrollHandler.ApproveExpenseClaim,
	)

	payrollRoutes.Get("/hrms.payroll.doctype.expense_claim.expense_claim.list_expense_claims",
		payrollHandler.ListExpenseClaims,
	)

	// ========== Performance Routes ==========
	performanceRoutes := authenticated.Group("")

	// Appraisal operations (HR Manager and above for creation)
	performanceRoutes.Post("/hrms.hr.doctype.appraisal.appraisal.create_appraisal",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		performanceHandler.CreateAppraisal,
	)

	performanceRoutes.Get("/hrms.hr.doctype.appraisal.appraisal.get_appraisal",
		performanceHandler.GetAppraisal,
	)

	performanceRoutes.Post("/hrms.hr.doctype.appraisal.appraisal.update_appraisal",
		performanceHandler.UpdateAppraisal,
	)

	performanceRoutes.Post("/hrms.hr.doctype.appraisal.appraisal.submit_appraisal",
		performanceHandler.SubmitAppraisal,
	)

	performanceRoutes.Post("/hrms.hr.doctype.appraisal.appraisal.complete_appraisal",
		permMiddleware.RequireRole("HR Manager", "System Manager"),
		performanceHandler.CompleteAppraisal,
	)

	performanceRoutes.Get("/hrms.hr.doctype.appraisal.appraisal.list_appraisals",
		performanceHandler.ListAppraisals,
	)

	// Goal operations
	performanceRoutes.Post("/hrms.hr.doctype.goal.goal.create_goal",
		performanceHandler.CreateGoal,
	)

	performanceRoutes.Get("/hrms.hr.doctype.goal.goal.get_goal",
		performanceHandler.GetGoal,
	)

	performanceRoutes.Post("/hrms.hr.doctype.goal.goal.update_goal",
		performanceHandler.UpdateGoal,
	)

	performanceRoutes.Get("/hrms.hr.doctype.goal.goal.list_goals",
		performanceHandler.ListGoals,
	)

	performanceRoutes.Get("/hrms.hr.doctype.goal.goal.get_active_goals",
		performanceHandler.GetActiveGoalsForEmployee,
	)

	// 360-degree feedback operations
	performanceRoutes.Post("/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.create_feedback",
		performanceHandler.CreateFeedback,
	)

	performanceRoutes.Get("/hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_feedback",
		performanceHandler.GetFeedbackForEmployee,
	)

	// Skill map operations
	performanceRoutes.Post("/hrms.hr.doctype.employee_skill_map.employee_skill_map.create_or_update",
		performanceHandler.CreateOrUpdateSkillMap,
	)

	performanceRoutes.Get("/hrms.hr.doctype.employee_skill_map.employee_skill_map.get_skill_map",
		performanceHandler.GetSkillMapByEmployee,
	)

	// Appraisal templates and cycles
	performanceRoutes.Get("/hrms.hr.doctype.appraisal_cycle.appraisal_cycle.get_active_cycles",
		performanceHandler.GetActiveAppraisalCycles,
	)

	performanceRoutes.Get("/hrms.hr.doctype.appraisal_template.appraisal_template.get_active_templates",
		performanceHandler.GetActiveAppraisalTemplates,
	)

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "hrms-api",
		})
	})

	// API info route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "HRMS API",
			"version": "1.0.0",
			"status":  "running",
		})
	})
}
