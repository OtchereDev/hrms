# Employee Lifecycle Module - Backend Documentation

## Overview

This module manages the complete employee journey from recruitment through exit, including onboarding, transfers, promotions, separations, and exit interviews. The backend implements complex boarding workflows, task automation, and employee property tracking.

---

## Core Doctypes

### 1. Employee Onboarding

**Purpose**: Automates new hire onboarding with task management and tracking.

**Auto-naming**: `HR-EMP-ONB-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| job_applicant | Link | Yes | Source job applicant |
| job_offer | Link | Yes | Accepted job offer |
| employee | Link | No | Created employee (auto-set) |
| employee_name | Data | Yes | Applicant name |
| date_of_joining | Date | Yes | Employee start date |
| employee_onboarding_template | Link | No | Template with pre-defined activities |
| company | Link | Yes | Hiring company |
| department | Link | No | Target department |
| designation | Link | No | Target designation |
| employee_grade | Link | No | Target grade |
| boarding_begins_on | Date | Yes | Onboarding process start date |
| boarding_status | Select | - | Pending/In Process/Completed (read-only) |
| activities | Table | - | Onboarding tasks |
| project | Link | - | Auto-created project for task tracking |

#### Business Logic

**Validation Rules**:
```python
def validate_duplicate_employee_onboarding():
    # Prevents multiple onboardings for same job applicant
    # Query: job_applicant = current AND name != self AND docstatus != 2

def set_employee():
    # Auto-links employee if already created from job applicant
    # Fetches existing employee where job_applicant matches

def validate_employee_creation():
    # Ensures all required_for_employee_creation activities are completed
    # Only applies when creating employee from onboarding
```

**On Submit**:
```python
def on_submit():
    # 1. Create Project
    project_name = f"Employee Onboarding: {job_applicant}"
    project.company = company
    project.expected_end_date = date_of_joining

    # 2. Create Tasks for each activity
    for activity in activities:
        task = frappe.new_doc("Task")
        task.project = project.name
        task.subject = activity.activity_name
        task.description = activity.description

        # Calculate dates
        start_date = add_days(boarding_begins_on, activity.begin_on or 0)
        end_date = add_days(start_date, activity.duration or 0)

        # Adjust for holidays
        if holiday_list:
            start_date = add_working_days(boarding_begins_on, activity.begin_on)
            end_date = add_working_days(start_date, activity.duration)

        task.exp_start_date = start_date
        task.exp_end_date = end_date

        # Assign users
        if activity.user:
            task.add_assigned_users([activity.user])
        if activity.role:
            users = get_users_with_role(activity.role)
            task.add_assigned_users(users)

        task.save()
        activity.task = task.name

    # 3. Update status
    boarding_status = "Pending"

    # 4. Send notifications
    if notify_users_by_email:
        send_activity_notifications()
```

**Creating Employee**:
```python
@frappe.whitelist()
def make_employee(source_name):
    # Maps onboarding to employee
    employee = frappe.new_doc("Employee")
    employee.first_name = job_applicant.applicant_name
    employee.personal_email = job_applicant.email_id
    employee.date_of_joining = date_of_joining
    employee.company = company
    employee.department = department
    employee.designation = designation
    employee.grade = employee_grade  # Note: maps to 'grade' field
    employee.status = "Active"
    employee.job_applicant = job_applicant
    employee.holiday_list = holiday_list

    return employee
```

**Status Updates**:
```python
def update_employee_boarding_status(project, user):
    # Called via hook when project/tasks are updated
    # Calculates boarding_status based on project completion

    if project.percent_complete == 0:
        boarding_status = "Pending"
    elif 0 < project.percent_complete < 100:
        boarding_status = "In Process"
    elif project.percent_complete == 100:
        boarding_status = "Completed"
```

**Mark as Completed**:
```python
def mark_onboarding_as_completed():
    # Marks all tasks as completed
    # Updates project status to "Completed"
    # Updates boarding_status to "Completed"

    for activity in activities:
        if activity.task:
            task = frappe.get_doc("Task", activity.task)
            if task.status != "Completed":
                task.status = "Completed"
                task.save()

    project.status = "Completed"
    project.percent_complete = 100
    project.save()
```

---

### 2. Employee Separation

**Purpose**: Manages offboarding process with task tracking.

**Auto-naming**: `HR-EMP-SEP-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee leaving |
| employee_name | Data | - | Employee name (read-only) |
| resignation_letter_date | Date | - | From employee record |
| boarding_begins_on | Date | Yes | Separation process start |
| employee_separation_template | Link | No | Template with exit activities |
| company | Link | Yes | Company |
| department | Link | - | Department (read-only) |
| boarding_status | Select | - | Status (read-only) |
| activities | Table | - | Exit tasks |
| project | Link | - | Auto-created project |

#### Business Logic

Inherits from `EmployeeBoardingController` - same logic as onboarding but for exits.

**Key Differences**:
- Uses resignation_letter_date as project start reference
- Links to departing employee instead of job applicant
- No employee creation functionality
- May link to Exit Interview

---

### 3. Employee Transfer

**Purpose**: Manages employee transfers with property updates and optional new employee ID creation.

**Auto-naming**: `HR-EMP-TRN-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee being transferred |
| transfer_date | Date | Yes | Effective transfer date |
| company | Link | - | Current company |
| new_company | Link | No | Target company (for inter-company transfers) |
| transfer_details | Table | Yes | Properties being changed |
| create_new_employee_id | Check | No | Create separate employee record |
| new_employee_id | Link | - | Created employee (if applicable) |

#### Business Logic

**Validation**:
```python
def validate_active_employee():
    employee = frappe.get_doc("Employee", self.employee)
    if employee.status != "Active":
        frappe.throw(_("Only Active Employees can be transferred"))

def validate_transfer_date():
    # On before_submit
    if transfer_date > today():
        frappe.throw(_("Transfer date cannot be in the future"))
```

**Property Transfer Logic**:
```python
def on_submit():
    employee = frappe.get_doc("Employee", self.employee)

    if create_new_employee_id:
        # Option 1: Create new employee ID
        new_employee = copy_employee(employee)

        # Update properties on new employee
        for detail in transfer_details:
            setattr(new_employee, detail.fieldname, detail.new)

        # Update work history
        new_employee.append("internal_work_history", {
            "from_date": transfer_date,
            "department": detail.new (if property is department),
            "designation": detail.new (if property is designation),
            # etc.
        })

        new_employee.insert()
        new_employee_id = new_employee.name

        # Close old employee
        employee.relieving_date = transfer_date
        employee.status = "Left"

        # Transfer user_id if not changing
        if user_id not in transfer_details:
            new_employee.user_id = employee.user_id
            employee.user_id = None

    else:
        # Option 2: Update in place
        for detail in transfer_details:
            setattr(employee, detail.fieldname, detail.new)

        # Update work history
        update_employee_work_history(employee, transfer_details, transfer_date)

    employee.save()
```

**Work History Tracking**:
```python
def update_employee_work_history(employee, changes, effective_date):
    # Close previous work history entry
    for history in employee.internal_work_history:
        if not history.to_date:
            history.to_date = add_days(effective_date, -1)

    # Create new work history entry
    employee.append("internal_work_history", {
        "from_date": effective_date,
        "to_date": None,
        "department": get_value_from_changes("department"),
        "designation": get_value_from_changes("designation"),
        "branch": get_value_from_changes("branch"),
    })
```

---

### 4. Employee Promotion

**Purpose**: Records promotions with property changes and CTC updates.

**Auto-naming**: `HR-EMP-PRO-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Employee being promoted |
| promotion_date | Date | Yes | Effective date |
| promotion_details | Table | - | Properties changing |
| current_ctc | Currency | - | Current CTC (if revising) |
| revised_ctc | Currency | - | New CTC |

#### Business Logic

**On Submit**:
```python
def on_submit():
    employee = frappe.get_doc("Employee", self.employee)

    # Update properties
    for detail in promotion_details:
        setattr(employee, detail.fieldname, detail.new)

    # Update work history
    update_employee_work_history(employee, promotion_details, promotion_date)

    # Update CTC
    if revised_ctc:
        employee.ctc = revised_ctc

    employee.save()
```

**On Cancel**:
```python
def on_cancel():
    employee = frappe.get_doc("Employee", self.employee)

    # Revert properties
    for detail in promotion_details:
        setattr(employee, detail.fieldname, detail.current)

    # Remove work history entry
    employee.internal_work_history = [
        h for h in employee.internal_work_history
        if h.from_date != promotion_date
    ]

    # Revert CTC
    if current_ctc:
        employee.ctc = current_ctc

    employee.save()
```

---

### 5. Exit Interview

**Purpose**: Manages exit interviews with questionnaire distribution and decision tracking.

**Auto-naming**: `HR-EXIT-INT-`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| employee | Link | Yes | Departing employee |
| company | Link | Yes | Company |
| status | Select | Yes | Pending/Scheduled/Completed/Cancelled |
| date | Date | Cond. | Interview date (required if Scheduled) |
| interviewers | Table | Cond. | Interviewers (required if Scheduled) |
| ref_doctype | Link | No | Questionnaire doctype |
| reference_document_name | Dynamic Link | No | Questionnaire document |
| questionnaire_email_sent | Check | - | Email sent flag |
| interview_summary | Text Editor | - | Interview notes |
| employee_status | Select | Cond. | Employee Retained/Exit Confirmed (required if Completed) |

#### Business Logic

**Validation**:
```python
def validate_relieving_date():
    employee = frappe.get_doc("Employee", self.employee)
    if not employee.relieving_date:
        frappe.throw(_("Please set relieving date for employee: {0}").format(employee.name))

def validate_duplicate_interview():
    # Only one interview per employee
    existing = frappe.db.exists("Exit Interview", {
        "employee": self.employee,
        "name": ("!=", self.name),
        "docstatus": ("<", 2)
    })
    if existing:
        frappe.throw(_("Exit Interview already exists for employee {0}").format(self.employee))
```

**Questionnaire Email**:
```python
@frappe.whitelist()
def send_exit_questionnaire(interviews):
    # Get settings
    hr_settings = frappe.get_single("HR Settings")
    if not hr_settings.exit_questionnaire_web_form:
        frappe.throw(_("Please set Exit Questionnaire Web Form in HR Settings"))

    if not hr_settings.exit_questionnaire_notification_template:
        frappe.throw(_("Please set Exit Questionnaire Notification Template in HR Settings"))

    template = frappe.get_doc("Email Template", hr_settings.exit_questionnaire_notification_template)

    for interview_name in interviews:
        interview = frappe.get_doc("Exit Interview", interview_name)

        # Generate questionnaire link
        web_form_url = get_url(hr_settings.exit_questionnaire_web_form)

        # Send email
        frappe.sendmail(
            recipients=[interview.email],
            subject=template.subject,
            message=frappe.render_template(template.response, {
                "doc": interview,
                "questionnaire_link": web_form_url
            })
        )

        interview.questionnaire_email_sent = 1
        interview.save()
```

**On Submit**:
```python
def on_submit():
    # Only allow if status is Completed
    if status != "Completed":
        frappe.throw(_("Only Completed interviews can be submitted"))

    # Update employee held_on date
    employee = frappe.get_doc("Employee", self.employee)
    employee.held_on = date
    employee.save()
```

---

### 6. Appointment Letter

**Purpose**: Generates formal appointment letters from templates.

**Auto-naming**: `HR-APP-LETTER-.#####`

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| job_applicant | Link | Yes | Selected candidate |
| applicant_name | Data | Yes | Applicant name (read-only) |
| company | Link | Yes | Offering company |
| appointment_date | Date | Yes | Appointment date |
| appointment_letter_template | Link | Yes | Letter template |
| introduction | Long Text | Yes | Letter intro (from template) |
| terms | Table | Yes | Terms and conditions |
| closing_notes | Text | No | Closing remarks |

#### Business Logic

**Template Loading**:
```python
@frappe.whitelist()
def get_appointment_letter_details(template):
    template_doc = frappe.get_doc("Appointment Letter Template", template)

    return {
        "introduction": template_doc.introduction,
        "closing_notes": template_doc.closing_notes,
        "terms": [
            {"offer_term": term.offer_term, "value": term.value}
            for term in template_doc.terms
        ]
    }
```

---

## Supporting Doctypes

### Employee Boarding Activity (Child Table)

Used in both Onboarding and Separation.

**Fields**:
- `activity_name`: Task description
- `user`: Specific user to assign
- `role`: Role whose members should be assigned
- `description`: Detailed instructions
- `begin_on`: Days from boarding_begins_on to start
- `duration`: Days to complete
- `task`: Created task reference
- `task_weight`: Task weight for project completion
- `required_for_employee_creation`: Must complete before creating employee (onboarding only)

**Task Assignment Logic**:
```python
def get_assigned_users(activity):
    users = []

    if activity.user:
        users.append(activity.user)

    if activity.role:
        role_users = frappe.get_all("Has Role",
            filters={"role": activity.role, "parenttype": "User"},
            fields=["parent"]
        )
        users.extend([u.parent for u in role_users if u.parent != "Administrator"])

    return list(set(users))  # Remove duplicates
```

---

### Employee Property History (Child Table)

Used in Transfer and Promotion to track property changes.

**Fields**:
- `property`: Display name of property
- `current`: Current value
- `new`: New value
- `fieldname`: Actual field name in Employee doctype

**Workflow**:
1. User clicks "Add Employee Property" button
2. Dialog shows: Select property → Shows current value → Enter new value
3. Row added to table with property, current, new, fieldname
4. On submit: Updates employee.<fieldname> = new
5. On cancel: Reverts employee.<fieldname> = current

---

## Automation & Scheduled Jobs

### Boarding Status Updates

**Hook**: `on_update` for Project and Task

```python
def update_employee_boarding_status(project, user=None):
    # Triggered when project/tasks are updated
    # Finds associated onboarding/separation

    onboarding = frappe.get_all("Employee Onboarding",
        filters={"project": project.name, "docstatus": 1},
        limit=1
    )

    if onboarding:
        doc = frappe.get_doc("Employee Onboarding", onboarding[0].name)

        # Calculate status from project completion
        if project.percent_complete == 0:
            status = "Pending"
        elif 0 < project.percent_complete < 100:
            status = "In Process"
        elif project.percent_complete == 100:
            status = "Completed"

        if doc.boarding_status != status:
            frappe.db.set_value("Employee Onboarding", doc.name, "boarding_status", status)
```

---

## Integration Points

### With Employee Master

**employee_master.py** in `/home/user/hrms/hrms/overrides/`:

```python
def validate_onboarding_process(self, method=None):
    # Links employee to onboarding when created
    if self.job_applicant:
        onboarding = frappe.db.get_value("Employee Onboarding",
            {"job_applicant": self.job_applicant, "docstatus": 1},
            "name"
        )

        if onboarding:
            frappe.db.set_value("Employee Onboarding", onboarding, "employee", self.name)

def update_job_applicant_and_offer(self, method=None):
    # Updates job applicant and offer status on employee creation
    if self.job_applicant:
        frappe.db.set_value("Job Applicant", self.job_applicant, "status", "Accepted")

        job_offer = frappe.db.get_value("Job Offer",
            {"job_applicant": self.job_applicant, "docstatus": 1},
            "name"
        )

        if job_offer:
            frappe.db.set_value("Job Offer", job_offer, "status", "Accepted")
```

### With Exit Interview

**Employee Separation** can link to **Exit Interview** for comprehensive offboarding tracking.

---

## Permissions

| Doctype | System Manager | HR Manager | HR User | Employee |
|---------|---------------|------------|---------|----------|
| Employee Onboarding | Full | Full (no delete) | - | - |
| Employee Separation | Full | - | - | - |
| Employee Transfer | - | Full (no delete) | Create/Submit | Read |
| Employee Promotion | - | Full (no delete) | Create/Submit | Read |
| Exit Interview | Full | - | - | - |
| Appointment Letter | Full | Full | - | - |

---

## Formula Reference

### Date Calculations

```python
# Task start date with holiday adjustment
start_date = add_working_days(boarding_begins_on, activity.begin_on, holiday_list)

# Task end date
end_date = add_working_days(start_date, activity.duration, holiday_list)
```

### Project Completion

```python
# Percent complete
total_weight = sum(task.task_weight for task in project.tasks)
completed_weight = sum(task.task_weight for task in project.tasks if task.status == "Completed")
percent_complete = (completed_weight / total_weight) * 100 if total_weight else 0
```

---

## Error Handling

**Common Validations**:
- Active employee check
- Future date prevention
- Duplicate prevention
- Required field validation
- Status-based action restrictions

**Error Messages**:
```python
# Examples from codebase
_("Only Active Employees can be transferred")
_("Transfer date cannot be in the future")
_("Employee Onboarding already exists for Job Applicant: {0}")
_("Please set relieving date for employee")
_("All required activities must be completed before creating employee")
```

---

This backend documentation provides complete technical details for implementing the Employee Lifecycle module functionality.
