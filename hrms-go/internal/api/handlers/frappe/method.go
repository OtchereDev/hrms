package frappe

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	frappeCore "github.com/OtchereDev/hrms-go/internal/core/frappe"
	"github.com/OtchereDev/hrms-go/internal/api/handlers"
)

// MethodHandler routes Frappe method calls to appropriate handlers
// This handles calls like /api/method/hrms.hr.doctype.leave_application.leave_application.apply_leave
type MethodHandler struct {
	db                 *gorm.DB
	attendanceHandler  *handlers.AttendanceHandler
	leaveHandler       *handlers.LeaveHandler
	payrollHandler     *handlers.PayrollHandler
	performanceHandler *handlers.PerformanceHandler
}

// NewMethodHandler creates a new method handler
func NewMethodHandler(
	db *gorm.DB,
	attendanceHandler *handlers.AttendanceHandler,
	leaveHandler *handlers.LeaveHandler,
	payrollHandler *handlers.PayrollHandler,
	performanceHandler *handlers.PerformanceHandler,
) *MethodHandler {
	return &MethodHandler{
		db:                 db,
		attendanceHandler:  attendanceHandler,
		leaveHandler:       leaveHandler,
		payrollHandler:     payrollHandler,
		performanceHandler: performanceHandler,
	}
}

// Call routes a method call to the appropriate handler
// This is the main router for /api/method/* calls
func (h *MethodHandler) Call(c *fiber.Ctx) error {
	method := c.Params("*")

	// Route to appropriate handler based on method path
	switch method {
	// Attendance methods
	case "hrms.hr.doctype.attendance.attendance.mark_attendance":
		return h.wrapHandler(h.attendanceHandler.MarkAttendance)(c)
	case "hrms.hr.doctype.attendance.attendance.get_attendance":
		return h.wrapHandler(h.attendanceHandler.GetAttendance)(c)
	case "hrms.hr.doctype.attendance.attendance.list_attendance":
		return h.wrapHandler(h.attendanceHandler.ListAttendance)(c)
	case "hrms.hr.doctype.attendance.attendance.get_monthly_attendance":
		return h.wrapHandler(h.attendanceHandler.GetMonthlyAttendance)(c)

	// Employee Checkin methods
	case "hrms.hr.doctype.employee_checkin.employee_checkin.checkin":
		return h.wrapHandler(h.attendanceHandler.Checkin)(c)
	case "hrms.hr.doctype.employee_checkin.employee_checkin.get_today_checkins":
		return h.wrapHandler(h.attendanceHandler.GetTodayCheckins)(c)

	// Attendance Request methods
	case "hrms.hr.doctype.attendance_request.attendance_request.create":
		return h.wrapHandler(h.attendanceHandler.CreateAttendanceRequest)(c)
	case "hrms.hr.doctype.attendance_request.attendance_request.approve":
		return h.wrapHandler(h.attendanceHandler.ApproveAttendanceRequest)(c)
	case "hrms.hr.doctype.attendance_request.attendance_request.reject":
		return h.wrapHandler(h.attendanceHandler.RejectAttendanceRequest)(c)

	// Attendance utility methods
	case "hrms.hr.doctype.attendance.attendance.get_unmarked_days":
		return h.wrapHandler(h.attendanceHandler.GetUnmarkedDays)(c)
	case "hrms.hr.doctype.attendance.attendance.mark_bulk_attendance":
		return h.wrapHandler(h.attendanceHandler.MarkBulkAttendance)(c)
	case "hrms.hr.doctype.attendance.attendance.get_events":
		return h.wrapHandler(h.attendanceHandler.GetEvents)(c)

	// Leave Application methods
	case "hrms.hr.doctype.leave_application.leave_application.apply_leave":
		return h.wrapHandler(h.leaveHandler.ApplyLeave)(c)
	case "hrms.hr.doctype.leave_application.leave_application.get":
		return h.wrapHandler(h.leaveHandler.GetLeaveApplication)(c)
	case "hrms.hr.doctype.leave_application.leave_application.list":
		return h.wrapHandler(h.leaveHandler.ListLeaveApplications)(c)
	case "hrms.hr.doctype.leave_application.leave_application.approve":
		return h.wrapHandler(h.leaveHandler.ApproveLeaveApplication)(c)
	case "hrms.hr.doctype.leave_application.leave_application.reject":
		return h.wrapHandler(h.leaveHandler.RejectLeaveApplication)(c)
	case "hrms.hr.doctype.leave_application.leave_application.cancel":
		return h.wrapHandler(h.leaveHandler.CancelLeaveApplication)(c)

	// Leave Allocation methods
	case "hrms.hr.doctype.leave_allocation.leave_allocation.allocate":
		return h.wrapHandler(h.leaveHandler.AllocateLeave)(c)
	case "hrms.hr.doctype.leave_allocation.leave_allocation.get_leave_balance":
		return h.wrapHandler(h.leaveHandler.GetLeaveBalance)(c)

	// Leave Type methods
	case "hrms.hr.doctype.leave_type.leave_type.get_active":
		return h.wrapHandler(h.leaveHandler.GetActiveLeaveTypes)(c)

	// Leave Encashment methods
	case "hrms.hr.doctype.leave_encashment.leave_encashment.create":
		return h.wrapHandler(h.leaveHandler.CreateLeaveEncashment)(c)

	// Leave utility methods
	case "hrms.hr.doctype.leave_application.leave_application.get_leave_details":
		return h.wrapHandler(h.leaveHandler.GetLeaveDetails)(c)
	case "hrms.hr.doctype.leave_application.leave_application.get_number_of_leave_days":
		return h.wrapHandler(h.leaveHandler.GetNumberOfLeaveDays)(c)
	case "hrms.hr.doctype.leave_application.leave_application.get_leave_balance_on":
		return h.wrapHandler(h.leaveHandler.GetLeaveBalanceOn)(c)
	case "hrms.hr.doctype.leave_application.leave_application.get_leaves_for_period":
		return h.wrapHandler(h.leaveHandler.GetLeavesForPeriod)(c)

	// Salary Slip methods
	case "hrms.payroll.doctype.salary_slip.salary_slip.generate":
		return h.wrapHandler(h.payrollHandler.GenerateSalarySlip)(c)
	case "hrms.payroll.doctype.salary_slip.salary_slip.get":
		return h.wrapHandler(h.payrollHandler.GetSalarySlip)(c)
	case "hrms.payroll.doctype.salary_slip.salary_slip.submit":
		return h.wrapHandler(h.payrollHandler.SubmitSalarySlip)(c)
	case "hrms.payroll.doctype.salary_slip.salary_slip.list":
		return h.wrapHandler(h.payrollHandler.ListSalarySlips)(c)

	// Salary Structure methods
	case "hrms.payroll.doctype.salary_structure.salary_structure.assign":
		return h.wrapHandler(h.payrollHandler.AssignSalaryStructure)(c)
	case "hrms.payroll.doctype.salary_structure.salary_structure.get_active_assignment":
		return h.wrapHandler(h.payrollHandler.GetActiveSalaryAssignment)(c)

	// Loan methods
	case "hrms.payroll.doctype.loan.loan.create":
		return h.wrapHandler(h.payrollHandler.CreateLoan)(c)
	case "hrms.payroll.doctype.loan.loan.approve":
		return h.wrapHandler(h.payrollHandler.ApproveLoan)(c)
	case "hrms.payroll.doctype.loan.loan.list":
		return h.wrapHandler(h.payrollHandler.ListLoans)(c)

	// Employee Advance methods
	case "hrms.payroll.doctype.employee_advance.employee_advance.create":
		return h.wrapHandler(h.payrollHandler.CreateEmployeeAdvance)(c)
	case "hrms.payroll.doctype.employee_advance.employee_advance.approve":
		return h.wrapHandler(h.payrollHandler.ApproveEmployeeAdvance)(c)

	// Expense Claim methods
	case "hrms.payroll.doctype.expense_claim.expense_claim.create":
		return h.wrapHandler(h.payrollHandler.CreateExpenseClaim)(c)
	case "hrms.payroll.doctype.expense_claim.expense_claim.approve":
		return h.wrapHandler(h.payrollHandler.ApproveExpenseClaim)(c)
	case "hrms.payroll.doctype.expense_claim.expense_claim.list":
		return h.wrapHandler(h.payrollHandler.ListExpenseClaims)(c)

	// Payroll utility methods
	case "hrms.payroll.doctype.salary_slip.salary_slip.get_salary_slip_details":
		return h.wrapHandler(h.payrollHandler.GetSalarySlipDetails)(c)
	case "hrms.payroll.doctype.loan.loan.calculate_amounts":
		return h.wrapHandler(h.payrollHandler.CalculateLoanAmounts)(c)
	case "hrms.payroll.doctype.salary_slip.salary_slip.calculate_net_pay":
		return h.wrapHandler(h.payrollHandler.CalculateNetPay)(c)
	case "hrms.payroll.doctype.payroll_entry.payroll_entry.get_payroll_summary":
		return h.wrapHandler(h.payrollHandler.GetPayrollSummary)(c)

	// Appraisal methods
	case "hrms.hr.doctype.appraisal.appraisal.create":
		return h.wrapHandler(h.performanceHandler.CreateAppraisal)(c)
	case "hrms.hr.doctype.appraisal.appraisal.get":
		return h.wrapHandler(h.performanceHandler.GetAppraisal)(c)
	case "hrms.hr.doctype.appraisal.appraisal.update":
		return h.wrapHandler(h.performanceHandler.UpdateAppraisal)(c)
	case "hrms.hr.doctype.appraisal.appraisal.submit":
		return h.wrapHandler(h.performanceHandler.SubmitAppraisal)(c)
	case "hrms.hr.doctype.appraisal.appraisal.complete":
		return h.wrapHandler(h.performanceHandler.CompleteAppraisal)(c)
	case "hrms.hr.doctype.appraisal.appraisal.list":
		return h.wrapHandler(h.performanceHandler.ListAppraisals)(c)

	// Goal methods
	case "hrms.hr.doctype.goal.goal.create":
		return h.wrapHandler(h.performanceHandler.CreateGoal)(c)
	case "hrms.hr.doctype.goal.goal.get":
		return h.wrapHandler(h.performanceHandler.GetGoal)(c)
	case "hrms.hr.doctype.goal.goal.update":
		return h.wrapHandler(h.performanceHandler.UpdateGoal)(c)
	case "hrms.hr.doctype.goal.goal.list":
		return h.wrapHandler(h.performanceHandler.ListGoals)(c)
	case "hrms.hr.doctype.goal.goal.get_active_for_employee":
		return h.wrapHandler(h.performanceHandler.GetActiveGoalsForEmployee)(c)

	// Employee Performance Feedback methods
	case "hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.create":
		return h.wrapHandler(h.performanceHandler.CreateFeedback)(c)
	case "hrms.hr.doctype.employee_performance_feedback.employee_performance_feedback.get_for_employee":
		return h.wrapHandler(h.performanceHandler.GetFeedbackForEmployee)(c)

	// Employee Skill Map methods
	case "hrms.hr.doctype.employee_skill_map.employee_skill_map.create_or_update":
		return h.wrapHandler(h.performanceHandler.CreateOrUpdateSkillMap)(c)
	case "hrms.hr.doctype.employee_skill_map.employee_skill_map.get_by_employee":
		return h.wrapHandler(h.performanceHandler.GetSkillMapByEmployee)(c)

	// Appraisal Cycle methods
	case "hrms.hr.doctype.appraisal_cycle.appraisal_cycle.get_active":
		return h.wrapHandler(h.performanceHandler.GetActiveAppraisalCycles)(c)

	// Appraisal Template methods
	case "hrms.hr.doctype.appraisal_template.appraisal_template.get_active":
		return h.wrapHandler(h.performanceHandler.GetActiveAppraisalTemplates)(c)

	default:
		return frappeCore.SendError(c, fiber.StatusNotFound, "Method not found: "+method, frappeCore.ErrTypeNotFound)
	}
}

// wrapHandler wraps our existing handlers to convert their responses to Frappe format
func (h *MethodHandler) wrapHandler(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Call the original handler
		// Our handlers already return JSON responses, but we may need to intercept
		// and convert them to Frappe format if they're not already using it

		// For now, just call the handler directly
		// In the future, we might want to intercept the response and convert it
		return handler(c)
	}
}
